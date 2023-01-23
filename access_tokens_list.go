package vrm

type AccessTokensList struct {
	Success bool `json:"success,omitempty"`
	Tokens  []struct {
		Name          string      `json:"name,omitempty"`
		IDAccessToken string      `json:"idAccessToken,omitempty"`
		CreatedOn     string      `json:"createdOn,omitempty"`
		Scope         string      `json:"scope,omitempty"`
		Expires       interface{} `json:"expires,omitempty"`
	} `json:"tokens,omitempty"`
}
