package vrm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const VRM_API string = "https://vrmapi.victronenergy.com/v2"

type resError struct {
	Success   bool   `json:"success"`
	Errors    string `json:"errors,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
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

func NewClient(user, pass string) *Client {
	return &Client{User: user, Password: pass}
}

type Client struct {
	User     string `json:"username"`
	Password string `json:"password"`
	logon    resLogin
	access   resAccessTokens
}

func (client *Client) SetAccount(user, pass string) {
	client.User, client.Password = user, pass
}

func (client *Client) SetLogon(token string, userID int64) {
	client.logon.Token, client.logon.UserID = token, userID
}

func (client *Client) SetAccess(token, tokenID string) {
	client.access.Token, client.access.TokenID = token, tokenID
}

func (client *Client) GetToken() string {
	return client.logon.Token
}

func (client *Client) GetUserID() int64 {
	return client.logon.UserID
}

func (client *Client) GetVerificationMode() string {
	return client.logon.VerificationMode
}

func (client *Client) GetVerificationSent() bool {
	return client.logon.VerificationSent
}

func (client *Client) GetAccessToken() string {
	return client.access.Token
}

func (client *Client) GetAccessTokenID() string {
	return client.access.TokenID
}

func (client *Client) GetLogonJson() string {
	res, _ := json.MarshalIndent(client.logon, "", "\t")
	return string(res)
}

func (client *Client) GetAccessJson() string {
	res, _ := json.MarshalIndent(client.access, "", "\t")
	return string(res)
}

// Connect - log in using an e-mail and password
//  Used to authenticate as a user to access authenticated routes.
//  2FA token must be included if 2FA is enabled on the account.
//  Returns a bearer token (JWT).
// https://vrm-api-docs.victronenergy.com/#/operations/auth/login
func (client *Client) Connect() error {
	buff, err := json.Marshal(client)
	if err != nil {
		return err
	}
	response, err := http.Post(VRM_API+"/auth/login", "application/json", bytes.NewBuffer(buff))
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
		if err := json.Unmarshal(res, &client.logon); err != nil {
			return err
		}
		return nil
	}
	resErr := resError{}
	if err := json.Unmarshal(res, &resErr); err == nil {
		if resErr.Success {
			return fmt.Errorf("%s", resErr.Errors)
		}
	}
	return fmt.Errorf("connect failed, code %v", response.StatusCode)
}

// Close - log out using a token
//  Used to log out a user. The token provided in the authorization header will be blacklisted
//  from the server and can no longer be used for authentication purposes.
// https://vrm-api-docs.victronenergy.com/#/operations/auth/logout
func (client *Client) Close() error {
	code, res, err := client.get(VRM_API + "/auth/logout")
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
		if resErr.Success {
			return fmt.Errorf("%s", resErr.Errors)
		}
	}
	return fmt.Errorf("connect failed, code %v", code)
}

// CreateAccessTokens - Create an access token for a user.
//  Users can create personal access tokens for usage with external services.
//  These tokens can be used as an alternative way of authentication against the VRM API.
//  The token is returned, after which it is not possible to ever retrieve it again.
// https://vrm-api-docs.victronenergy.com/#/operations/users/idUser/accesstokens/create
func (client *Client) CreateAccessTokens(name string) error {
	type Data struct {
		Name string `json:"name"`
	}
	data := Data{Name: name}
	buff, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(buff)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/%d/accesstokens/create", VRM_API, client.logon.UserID), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization", "Bearer "+client.logon.Token)
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
		if err := json.Unmarshal(res, &client.access); err != nil {
			return err
		}
		return nil
	}
	resErr := resError{}
	if err := json.Unmarshal(res, &resErr); err == nil {
		if resErr.Success {
			return fmt.Errorf("%s", resErr.Errors)
		}
	}
	return fmt.Errorf("connect failed, code %v", response.StatusCode)
}

// RequestList - return requests list
//
// https://vrm-api-docs.victronenergy.com
func (client *Client) RequestsList() []string {
	return []string{
		"/system-overview",
		"/diagnostics",
		"/gps-download",
		"/tags",
		"/data-download",
		"/stats",
		"/overallstats",
		"/widgets/Graph",
		"/widgets/GPS",
		"/widgets/HoursOfAc",
		"/widgets/GeneratorState",
		"/widgets/InputState",
		"/widgets/InverterState",
		"/widgets/MPPTState",
		"/widgets/ChargerState",
		"/widgets/EssBatteryLifeState",
		"/widgets/FuelCellState",
		"/widgets/BatteryExternalRelayState",
		"/widgets/BatteryRelayState",
		"/widgets/BatteryMonitorWarningsAndAlarms",
		"/widgets/GatewayRelayState",
		"/widgets/GatewayRelayTwoState",
		"/widgets/ChargerRelayState",
		"/widgets/SolarChargerRelayState",
		"/widgets/VeBusState",
		"/widgets/VeBusWarningsAndAlarms",
		"/widgets/InverterChargerState",
		"/widgets/InverterChargerWarningsAndAlarms",
		"/widgets/BatterySummary",
		"/widgets/BMSDiagnostics",
		"/widgets/HistoricData",
		"/widgets/IOExtenderInOut",
		"/widgets/LithiumBMS",
		"/widgets/DCMeter",
		"/widgets/EvChargerSummary",
		"/widgets/MeteorologicalSensor",
		"/widgets/GlobalLinkSummary",
		"/widgets/MotorSummary",
		"/widgets/PVInverterStatus",
		"/widgets/SolarChargerSummary",
		"/widgets/Status",
		"/widgets/TankSummary",
		"/widgets/TempSummaryAndGraph",
	}
}

func (client *Client) Get(siteID, req string) ([]byte, error) {
	code, res, err := client.get(fmt.Sprintf("%s/installations/%s%s", VRM_API, siteID, req))
	if err != nil {
		return res, err
	}
	if code == 200 {
		return res, err
	}
	resErr := resError{}
	if err := json.Unmarshal(res, &resErr); err == nil {
		if resErr.Success {
			return res, fmt.Errorf("%s", resErr.Errors)
		}
	}
	return res, fmt.Errorf("connect failed, code %v", code)
}

func (client *Client) get(url string) (code int, response []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if client.access.Token != "" {
		req.Header.Set("X-Authorization", "Token "+client.access.Token)
	} else {
		req.Header.Set("X-Authorization", "Bearer "+client.logon.Token)
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
