package vrm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type apierr struct {
	Success   bool        `json:"success,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	ErrorCode string      `json:"error_code,omitempty"`
}

type logon struct {
	Token            string `json:"token"`
	UserID           int64  `json:"idUser,omitempty"`
	VerificationMode string `json:"verification_mode,omitempty"`
	VerificationSent bool   `json:"verification_sent,omitempty"`
}

type access struct {
	Success bool   `json:"success,omitempty"`
	Token   string `json:"token"`
	TokenID string `json:"idAccessToken,omitempty"`
}

type user struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type client struct {
	user   user
	logon  logon
	access access
	mux    sync.Mutex
}

func New() *client {
	return &client{}
}

//SetUser - set authorization parameters
func (c *client) SetUser(name, pass string) error {
	c.mux.Lock()
	c.user.Name, c.user.Password = name, pass
	c.mux.Unlock()
	return nil
}

//SetLogon - set current logon, if connection is not closed
func (c *client) SetLogon(token string, userID int64) error {
	c.mux.Lock()
	c.logon.Token, c.logon.UserID = token, userID
	c.mux.Unlock()
	return nil
}

//SetAccess - set access, if access is not closed
func (c *client) SetAccess(token, tokenID string) error {
	c.mux.Lock()
	c.access.Token, c.access.TokenID = token, tokenID
	c.mux.Unlock()
	return nil
}

func (c *client) GetToken() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.logon.Token
}

func (c *client) GetUserID() int64 {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.logon.UserID
}

func (c *client) GetAccessToken() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.access.Token
}

func (c *client) GetAccessTokenID() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.access.TokenID
}

//SetUserJson - set authorization parameters
func (c *client) SetUserJson(config []byte) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	return json.Unmarshal(config, &c.user)
}

func (c *client) GetUserJson() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	res, _ := json.MarshalIndent(c.user, "", "\t")
	return string(res)
}

//SetLogonJson - set current logon, if connection is not closed
func (c *client) SetLogonJson(config []byte) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	return json.Unmarshal(config, &c.logon)
}

func (c *client) GetLogonJson() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	res, _ := json.MarshalIndent(c.logon, "", "\t")
	return string(res)
}

//SetAccessJson - set access, if access is not closed
func (c *client) SetAccessJson(config []byte) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	return json.Unmarshal(config, &c.logon)
}

func (c *client) GetAccessJson() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	res, _ := json.MarshalIndent(c.access, "", "\t")
	return string(res)
}

type User interface {
	GetName() string
	GetPass() string
}

//SetUserInterface - set authorization parameters
func (c *client) SetUserInterface(u User) error {
	c.mux.Lock()
	c.user.Name, c.user.Password = u.GetName(), u.GetPass()
	c.mux.Unlock()
	return nil
}

type Access interface {
	GetToken() string
	GetTokenID() string
}

//SetAccessInterface - set access, if access is not closed
func (c *client) SetAccessInterface(a Access) error {
	c.mux.Lock()
	c.access.Token, c.access.TokenID = a.GetToken(), a.GetTokenID()
	c.mux.Unlock()
	return nil
}

type Logon interface {
	GetToken() string
	GetUserID() int64
}

//SetLogonInterface - set current logon, if connection is not closed
func (c *client) SetLogonInterface(l Logon) error {
	c.mux.Lock()
	c.logon.Token, c.logon.UserID = l.GetToken(), l.GetUserID()
	c.mux.Unlock()
	return nil
}

// Connect - Log in using an e-mail and password
//  Used to authenticate as a user to access authenticated routes.
//  2FA token must be included if 2FA is enabled on the account.
//  Returns a bearer token (JWT).
// https://vrm-api-docs.victronenergy.com/#/operations/auth/login
func (c *client) Connect() error {
	c.mux.Lock()
	defer c.mux.Unlock()
	buff, err := json.Marshal(c.user)
	if err != nil {
		return err
	}
	response, err := http.Post(apiurl("/auth/login").String(), "application/json", bytes.NewBuffer(buff))
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
		if err := json.Unmarshal(res, &c.logon); err != nil {
			return err
		}
		return nil
	}
	aerr := apierr{}
	if err := json.Unmarshal(res, &aerr); err == nil {
		return fmt.Errorf("connect failed, %v", aerr.Errors)
	}
	return fmt.Errorf("connect failed, code %v", response.StatusCode)
}

// Close - Log out using a token
//  Used to log out a user. The token provided in the authorization header will be blacklisted
//  from the server and can no longer be used for authentication purposes.
// https://vrm-api-docs.victronenergy.com/#/operations/auth/logout
func (c *client) Close() error {
	res, err := c.get(apiurl("/auth/logout").String())
	if err != nil {
		return err
	}
	if string(res) != "{\"token\":\"\"}" {
		return errors.New("logout fail")
	}
	return nil
}

// CreateAccessTokens - Create an access token for a user.
//  Users can create personal access tokens for usage with external services.
//  These tokens can be used as an alternative way of authentication against the VRM API.
//  The token is returned, after which it is not possible to ever retrieve it again.
// https://vrm-api-docs.victronenergy.com/#/operations/users/idUser/accesstokens/create
func (c *client) CreateAccessTokens(name string) error {
	type Data struct {
		Name string `json:"name"`
	}
	data := Data{Name: name}
	buff, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(buff)
	u := apiurl("/users/%v/accesstokens/create", c.GetUserID())
	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization", "Bearer "+c.GetToken())
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
		c.mux.Lock()
		defer c.mux.Unlock()
		if err := json.Unmarshal(res, &c.access); err != nil {
			return err
		}
		return nil
	}
	aerr := apierr{}
	if err := json.Unmarshal(res, &aerr); err == nil {
		return fmt.Errorf("create access tokens failed, %v", aerr.Errors)
	}
	return fmt.Errorf("create access tokens failed, code %v", response.StatusCode)
}

// RevokeAccessTokens - Revoke an access token for a user.
//  Revokes one or more personal access token for a user.
// https://vrm-api-docs.victronenergy.com/#/operations/users/idUser/accesstokens/revoke
func (c *client) RevokeAccessTokens(name string) error {
	type Data struct {
		Success bool `json:"success"`
		Data    struct {
			Removed int64 `json:"removed"`
		} `json:"data"`
	}
	u := apiurl("/users/%v/accesstokens/%v/revoke", c.GetUserID(), c.GetAccessTokenID())
	res, err := c.get(u.String())
	if err != nil {
		return err
	}
	data := Data{}
	if err := json.Unmarshal(res, &data); err != nil {
		return fmt.Errorf("revoke access tokens failed, %v", err)
	}
	if !data.Success {
		return fmt.Errorf("revoke access tokens failed, %+v", data)
	}
	return nil
}

func (c *client) GetInstallations(i *Installations) error {
	u := apiurl("/users/%v/installations", c.GetUserID())
	u.RawQuery = "extended=1"
	res, err := c.get(u.String())
	if err != nil {
		return err
	}
	return json.Unmarshal(res, i)
}

func (c *client) GetAccessTokensList(atl *AccessTokensList) error {
	u := apiurl("/users/%v/accesstokens/list", c.GetUserID())
	u.RawQuery = "extended=1"
	res, err := c.get(u.String())
	if err != nil {
		return err
	}
	return json.Unmarshal(res, atl)
}

// Get - get api /installations/...
//  recommend call RequestsList
func (c *client) Get(siteID int, request, query string) ([]byte, error) {
	u := apiurl("/installations/%v/%v", siteID, request)
	u.RawQuery = query
	return c.get(u.String())
}

// GetObject - get api /installations/...
//  recommend call RequestsList
func (c *client) GetObject(object any, siteID int, request, query string) error {
	res, err := c.Get(siteID, request, query)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, object)
}

func (c *client) get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.access.Token != "" {
		req.Header.Set("X-Authorization", "Token "+c.GetAccessToken())
	} else {
		req.Header.Set("X-Authorization", "Bearer "+c.GetToken())
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	response.Close = true
	res, err := ioutil.ReadAll(response.Body)
	if response.StatusCode == 200 {
		return res, err
	}
	apierr := apierr{}
	if err := json.Unmarshal(res, &apierr); err == nil {
		return nil, fmt.Errorf("get url failed, %v", apierr.Errors)
	}
	return nil, fmt.Errorf("get url failed, code %v", response.StatusCode)
}

// RequestList - return requests list
//
// https://vrm-api-docs.victronenergy.com
func (c *client) RequestsList() []string {
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

// ***********************************************
// ***********************************************
// ***********************************************

func apiurl(path string, v ...any) *url.URL {
	return &url.URL{Scheme: "https", Host: "vrmapi.victronenergy.com", Path: fmt.Sprintf("/v2"+path, v...)}
}
