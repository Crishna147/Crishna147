package client

import (
	"fmt"

	"github.com/EnsurityTechnologies/ensweb"
	"github.com/EnsurityTechnologies/thincvirtual/authserver/server"
)

func (c *Client) AddDevices(adr []server.AddDeviceRequest) error {
	var br ensweb.BaseResponse
	err := c.SendJSON("POST", server.APIAddDevice, true, nil, adr, &br, nil)
	if err != nil {
		c.log.Error("failed to add devices", "err", err)
		return err
	}
	if !br.Status {
		c.log.Error("failed to add devices", "msg", br.Message)
		return fmt.Errorf(br.Message)
	}
	c.log.Info(br.Message)
	return nil
}

func (c *Client) BlockDevice(serialNumber string) error {
	dr := server.DeviceRequest{
		SerialNumber: serialNumber,
	}
	var br ensweb.BaseResponse
	err := c.SendJSON("POST", server.APIBlockDevice, true, nil, dr, &br, nil)
	if err != nil {
		c.log.Error("failed to block device", "err", err)
		return err
	}
	if !br.Status {
		c.log.Error("failed to block device", "msg", br.Message)
		return fmt.Errorf(br.Message)
	}
	c.log.Info(br.Message)
	return nil
}

func (c *Client) UnBlockDevice(serialNumber string) error {
	dr := server.DeviceRequest{
		SerialNumber: serialNumber,
	}
	var br ensweb.BaseResponse
	err := c.SendJSON("POST", server.APIUnBlockDevice, true, nil, dr, &br, nil)
	if err != nil {
		c.log.Error("failed to unblock device", "err", err)
		return err
	}
	if !br.Status {
		c.log.Error("failed to unblock device", "msg", br.Message)
		return fmt.Errorf(br.Message)
	}
	c.log.Info(br.Message)
	return nil
}

func (c *Client) GetDeviceStatus(serialNumber string) (bool, error) {
	dr := server.DeviceRequest{
		SerialNumber: serialNumber,
	}
	var br server.DeviceStatus
	err := c.SendJSON("POST", server.APIDeviceStatus, false, nil, dr, &br, nil)
	if err != nil {
		c.log.Error("failed to get device status", "err", err)
		return false, err
	}
	if !br.Status {
		c.log.Error("failed to get device status", "msg", br.Message)
		return false, fmt.Errorf(br.Message)
	}
	c.log.Info(br.Message)
	return br.IsBlocked, nil
}
