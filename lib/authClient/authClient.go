package authClient

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	pbauth "lib/proto"
	"os"
)

type AuthClient struct {
	client pbauth.AuthRpcClient
	conn   *grpc.ClientConn
}

func InitAuthClient() (authClient AuthClient, err error) {
	options := []grpc.DialOption{grpc.WithInsecure()}
	address := os.Getenv("AUTH_ADDRESS")
	if len(address) == 0 {
		return authClient, errors.New("Auth server address is not provided")
	}
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(address, options...); err != nil {
		return
	}
	authClient.conn = conn
	authClient.client = pbauth.NewAuthRpcClient(conn)
	return
}

type ErrorRespStatus struct {
	StatusCode int
	ErrorResp  error
}

func (e ErrorRespStatus) Error() string {
	return e.ErrorResp.Error()
}

type UserPermission struct {
	Username string
	Admin    bool
}

func (c *AuthClient) Validate(token string) (permission UserPermission, err error) {
	req := &pbauth.ValidateRequest{
		AccessToken: token,
	}
	var resp *pbauth.ValidateResponse
	if resp, err = c.client.Validate(context.Background(), req); err != nil {
		return
	}
	permission.Username = resp.Username
	permission.Admin = resp.Admin
	return
}

var errForbidden = errors.New("Not enough permissions")

func (c *AuthClient) EnsurePermission(token string, adminRequired bool) (err error) {
	var userPermission UserPermission
	if userPermission, err = c.Validate(token); err != nil {
		return
	}
	if !userPermission.Admin {
		if adminRequired {
			return errForbidden
		}
	}
	return
}
