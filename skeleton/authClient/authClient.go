package authClient

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type AuthClient struct {
	client  http.Client
	address string
}

func InitAuthClient() (authClient AuthClient, err error) {
	authClient = AuthClient{}
	authClient.address = os.Getenv("AUTH_ADDRESS")
	if len(authClient.address) == 0 {
		return authClient, errors.New("Auth server address is not provided")
	}
	return
}

type errorResp struct {
	Error string `json:"error"`
}

type ErrorRespStatus struct {
	StatusCode int
	ErrorResp  error
}

func (e ErrorRespStatus) Error() string {
	return e.ErrorResp.Error()
}

func (c *AuthClient) Validate(token string) (err error) {
	var req *http.Request
	if req, err = http.NewRequest("GET", c.address, nil); err != nil {
		return
	}
	req.Header.Set("auth", token)
	var resp *http.Response
	if resp, err = c.client.Do(req); err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		var errResp errorResp
		var err error
		if err = json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			err = errors.New(errResp.Error)
		}
		return &ErrorRespStatus{
			StatusCode: resp.StatusCode,
			ErrorResp:  err,
		}
	}
	return
}
