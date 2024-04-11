package client

import (
	"fmt"

	"github.com/EnsurityTechnologies/thincvirtual/authserver/server"
)

func (c *Client) ResetUser(username string) error {
	ureq := server.UserRequest{
		UserName: username,
	}
	var uresp server.UserResponse
	err := c.SendJSON("POST", server.APIResetUser, true, nil, ureq, &uresp, nil)
	if err != nil {
		c.log.Error("failed to reset user", "err", err)
		return err
	}
	if !uresp.Status {
		c.log.Error("failed to reset user", "msg", uresp.Message)
		return fmt.Errorf(uresp.Message)
	}
	c.log.Info(uresp.Message)
	return nil
}

func (c *Client) SetMaxFingers(username string, numFingers int) error {
	ureq := server.UserRequest{
		UserName:      username,
		MaxNumFingers: numFingers,
	}
	var uresp server.UserResponse
	err := c.SendJSON("POST", server.APISetNumFingers, true, nil, ureq, &uresp, nil)
	if err != nil {
		c.log.Error("failed to set maximum fingers", "err", err)
		return err
	}
	if !uresp.Status {
		c.log.Error("failed to set maximum fingers", "msg", uresp.Message)
		return fmt.Errorf(uresp.Message)
	}
	c.log.Info(uresp.Message)
	return nil
}

func (c *Client) SetPinRequired(username string, isPinRequired bool) error {
	ureq := server.UserRequest{
		UserName:      username,
		IsPinRequired: isPinRequired,
	}
	var uresp server.UserResponse
	err := c.SendJSON("POST", server.APISetPinRequired, true, nil, ureq, &uresp, nil)
	if err != nil {
		c.log.Error("failed to set pin settings", "err", err)
		return err
	}
	if !uresp.Status {
		c.log.Error("failed to set pin settings", "msg", uresp.Message)
		return fmt.Errorf(uresp.Message)
	}
	c.log.Info(uresp.Message)
	return nil
}

func (c *Client) FingerEnrolled(username string, serialNumber string) error {
	ureq := server.UserRequest{
		UserName:     username,
		SerialNumber: serialNumber,
	}
	var uresp server.UserResponse
	err := c.SendJSON("POST", server.APIFingerEnrolled, true, nil, ureq, &uresp, nil)
	if err != nil {
		c.log.Error("failed to update finger enrolled state", "err", err)
		return err
	}
	if !uresp.Status {
		c.log.Error("failed to update finger enrolled state", "msg", uresp.Message)
		return fmt.Errorf(uresp.Message)
	}
	c.log.Info(uresp.Message)
	return nil
}

func (c *Client) KeyAdded(username string, serialNumber string) error {
	ureq := server.UserRequest{
		UserName:     username,
		SerialNumber: serialNumber,
	}
	var uresp server.UserResponse
	err := c.SendJSON("POST", server.APIKeyAdded, true, nil, ureq, &uresp, nil)
	if err != nil {
		c.log.Error("failed to update key status", "err", err)
		return err
	}
	if !uresp.Status {
		c.log.Error("failed to update key status", "msg", uresp.Message)
		return fmt.Errorf(uresp.Message)
	}
	c.log.Info(uresp.Message)
	return nil
}
