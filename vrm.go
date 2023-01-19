package vrm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const VRM_API_HOST string = "https://vrmapi.victronenergy.com"

type authorization struct {
	Name string `json:"username"`
	Pass string `json:"password"`
}

type resError struct {
	Success   bool        `json:"success"`
	Errors    interface{} `json:"errors,omitempty"`
	ErrorCode string      `json:"error_code,omitempty"`
}

type resLogin struct {
	Token            string `json:"token"`
	UserID           int64  `json:"idUser,omitempty"`
	VerificationMode string `json:"verification_mode,omitempty"`
	VerificationSent bool   `json:"verification_sent,omitempty"`
}

type resAccessTokens struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	TokenID string `json:"idAccessToken"`
}

var (
	user   authorization
	logon  resLogin
	access resAccessTokens
)

func urlVRM(path string, v ...any) *url.URL {
	return &url.URL{Scheme: "https", Host: "vrmapi.victronenergy.com", Path: fmt.Sprintf("/v2"+path, v...)}
}

func SetAccount(name, pass string) {
	user.Name, user.Pass = name, pass
}

func SetLogon(token string, userID int64) {
	logon.Token, logon.UserID = token, userID
}

func SetAccess(token, tokenID string) {
	access.Token, access.TokenID = token, tokenID
}

func GetToken() string {
	return logon.Token
}

func GetUserID() int64 {
	return logon.UserID
}

func GetVerificationMode() string {
	return logon.VerificationMode
}

func GetVerificationSent() bool {
	return logon.VerificationSent
}

func GetAccessToken() string {
	return access.Token
}

func GetAccessTokenID() string {
	return access.TokenID
}

func GetLogonJson() string {
	res, _ := json.MarshalIndent(logon, "", "\t")
	return string(res)
}

func GetAccessJson() string {
	res, _ := json.MarshalIndent(access, "", "\t")
	return string(res)
}

// NewConnect - log in using an e-mail and password
//  Used to authenticate as a user to access authenticated routes.
//  2FA token must be included if 2FA is enabled on the account.
//  Returns a bearer token (JWT).
// https://vrm-api-docs.victronenergy.com/#/operations/auth/login
func NewConnect(name, pass string) error {
	user.Name, user.Pass = name, pass
	return Connect()
}

// Connect - log in using an e-mail and password
//  Used to authenticate as a user to access authenticated routes.
//  2FA token must be included if 2FA is enabled on the account.
//  Returns a bearer token (JWT).
// https://vrm-api-docs.victronenergy.com/#/operations/auth/login
func Connect() error {
	buff, err := json.Marshal(user)
	if err != nil {
		return err
	}
	response, err := http.Post(urlVRM("/auth/login").String(), "application/json", bytes.NewBuffer(buff))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	response.Close = true
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		if err := json.Unmarshal(res, &logon); err != nil {
			return err
		}
		return nil
	}
	resErr := resError{}
	if err := json.Unmarshal(res, &resErr); err == nil {
		return fmt.Errorf("connect failed, %v", resErr.Errors)
	}
	return fmt.Errorf("connect failed, code %v", response.StatusCode)
}

// Close - log out using a token
//  Used to log out a user. The token provided in the authorization header will be blacklisted
//  from the server and can no longer be used for authentication purposes.
// https://vrm-api-docs.victronenergy.com/#/operations/auth/logout
func Close() error {
	code, res, err := get(urlVRM("/auth/logout").String())
	if err != nil {
		return err
	}
	if code == 200 {
		if string(res) != "{\"token\":\"\"}" {
			return errors.New("logout fail")
		}
		return nil
	}
	resErr := resError{}
	if err := json.Unmarshal(res, &resErr); err == nil {
		return fmt.Errorf("close failed, %v", resErr.Errors)
	}
	return fmt.Errorf("close failed, code %v", code)
}

// CreateAccessTokens - Create an access token for a user.
//  Users can create personal access tokens for usage with external services.
//  These tokens can be used as an alternative way of authentication against the VRM API.
//  The token is returned, after which it is not possible to ever retrieve it again.
// https://vrm-api-docs.victronenergy.com/#/operations/users/idUser/accesstokens/create
func CreateAccessTokens(name string) error {
	type Data struct {
		Name string `json:"name"`
	}
	data := Data{Name: name}
	buff, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(buff)
	req, err := http.NewRequest("POST", urlVRM("/users/%d/accesstokens/create", logon.UserID).String(), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization", "Bearer "+logon.Token)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	response.Close = true
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		if err := json.Unmarshal(res, &access); err != nil {
			return err
		}
		return nil
	}
	resErr := resError{}
	if err := json.Unmarshal(res, &resErr); err == nil {
		return fmt.Errorf("create access tokens failed, %v", resErr.Errors)
	}
	return fmt.Errorf("create access tokens failed, code %v", response.StatusCode)
}

func Get(siteID, request, query string) ([]byte, error) {
	u := urlVRM("/installations/%s/%s", siteID, request)
	u.RawQuery = query
	code, res, err := get(u.String())
	if err != nil {
		return res, err
	}
	if code == 200 {
		return res, err
	}
	resErr := resError{}
	if err := json.Unmarshal(res, &resErr); err == nil {
		return res, fmt.Errorf("get url failed, %v", resErr.Errors)
	}
	return res, fmt.Errorf("get url failed, code %v", code)
}

func GetObject(object any, siteID, request, query string) error {
	res, err := Get(siteID, "stats", query)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, object)
}

func get(url string) (code int, response []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if access.Token != "" {
		req.Header.Set("X-Authorization", "Token "+access.Token)
	} else {
		req.Header.Set("X-Authorization", "Bearer "+logon.Token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	res.Close = true
	code = res.StatusCode
	response, err = ioutil.ReadAll(res.Body)
	return
}

// RequestList - return requests list
//
// https://vrm-api-docs.victronenergy.com
func RequestsList() []string {
	return []string{
		"system-overview",
		"diagnostics",
		"gps-download",
		"tags",
		"data-download",
		"stats",
		"overallstats",
		"widgets/Graph",
		"widgets/GPS",
		"widgets/HoursOfAc",
		"widgets/GeneratorState",
		"widgets/InputState",
		"widgets/InverterState",
		"widgets/MPPTState",
		"widgets/ChargerState",
		"widgets/EssBatteryLifeState",
		"widgets/FuelCellState",
		"widgets/BatteryExternalRelayState",
		"widgets/BatteryRelayState",
		"widgets/BatteryMonitorWarningsAndAlarms",
		"widgets/GatewayRelayState",
		"widgets/GatewayRelayTwoState",
		"widgets/ChargerRelayState",
		"widgets/SolarChargerRelayState",
		"widgets/VeBusState",
		"widgets/VeBusWarningsAndAlarms",
		"widgets/InverterChargerState",
		"widgets/InverterChargerWarningsAndAlarms",
		"widgets/BatterySummary",
		"widgets/BMSDiagnostics",
		"widgets/HistoricData",
		"widgets/IOExtenderInOut",
		"widgets/LithiumBMS",
		"widgets/DCMeter",
		"widgets/EvChargerSummary",
		"widgets/MeteorologicalSensor",
		"widgets/GlobalLinkSummary",
		"widgets/MotorSummary",
		"widgets/PVInverterStatus",
		"widgets/SolarChargerSummary",
		"widgets/Status",
		"widgets/TankSummary",
		"widgets/TempSummaryAndGraph",
	}
}
