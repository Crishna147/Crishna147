package client

import (
	"fmt"

	"github.com/EnsurityTechnologies/thincvirtual/authserver/server"
)

func (c *Client) Login(username string, pwd string) (*server.LoginResponse, error) {
	lr := server.LoginRequest{
		UserName: username,
		Password: pwd,
	}
	var lrs server.LoginResponse
	err := c.SendJSON("POST", server.APIAdminLogin, false, nil, lr, &lrs, nil)
	if err != nil {
		c.log.Error("user login failed", "err", err)
		return nil, err
	}
	if !lrs.Status {
		c.log.Error("user login failed", "msg", lrs.Message)
		return nil, fmt.Errorf(lrs.Message)
	}
	c.SetToken(lrs.Token)
	return &lrs, nil
}

func (c *Client) UserLogin(username string, uid string, serialNumbr string) (*server.UserResponse, error) {
	lr := server.UserRequest{
		UserName:     username,
		UID:          uid,
		SerialNumber: serialNumbr,
	}
	var ur server.UserResponse
	err := c.SendJSON("POST", server.APIUserLogin, false, nil, lr, &ur, nil)
	if err != nil {
		c.log.Error("user login failed", "err", err)
		return nil, err
	}
	if !ur.Status {
		c.log.Error("user login failed", "msg", ur.Message)
		return nil, fmt.Errorf(ur.Message)
	}
	c.SetToken(ur.Token)
	return &ur, nil
}
