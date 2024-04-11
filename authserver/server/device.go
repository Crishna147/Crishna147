package server

import (
	"fmt"
	"net/http"

	"github.com/EnsurityTechnologies/ensweb"
)

func (s *Server) getDevice(serialNumber string) (*Device, error) {
	var d Device
	err := s.ts.Read(DeviceTable, serialNumber, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *Server) putDevice(d *Device) error {
	err := s.ts.Write(DeviceTable, d.SerialNumber, d)
	if err != nil {
		return err
	}
	return nil
}
// addDevice godoc
// @Summary 			Add Device
// @Description 		Add multiple devices.
// @Tags  				Devices
// @Produce 			json
// @Accept 				json
// @Param 				data body AddDeviceRequest true "Array of AddDeviceRequest objects"
// @Success 			200 {object} DeviceRequest
// @Failure             401 {object} ensweb.BaseResponse
// @Failure             400 {object} ensweb.BaseResponse
// @Failure             500 {object} ensweb.BaseResponse
// @Router 				/api/adddevice [post]
func (s *Server) addDevice(req *ensweb.Request) *ensweb.Result {
	var adr []AddDeviceRequest
	br := &ensweb.BaseResponse{
		Status: false,
	}
	err := s.ParseJSON(req, &adr)
	if err != nil {
		br.Message = "failed add device devices, failed to unmarshal json"
		return s.RenderJSON(req, br, http.StatusBadRequest)
	}
	l, err := s.getLicense()
	if err != nil {
		br.Message = "failed add device devices, failed to read license"
		return s.RenderJSON(req, br, http.StatusInternalServerError)
	}
	defer s.releaseLicense()
	numDevices := l.NumberDevices
	count := 0
	for i := range adr {
		if l.MaxNumberDevices == numDevices {
			br.Message = "Only " + fmt.Sprintf("%d", count) + " devices licenses are remaining"
			return s.RenderJSON(req, br, http.StatusOK)
		}
		_, err = s.getDevice(adr[i].SerialNumber)
		if err == nil {
			br.Message = fmt.Sprintf("Device (%s) already added", adr[i].SerialNumber)
			return s.RenderJSON(req, br, http.StatusOK)
		}
		numDevices++
		count++
	}
	for i := range adr {
		d := &Device{
			SerialNumber: adr[i].SerialNumber,
			PublicKey:    adr[i].PublicKey,
			IsBlocked:    false,
		}
		err = s.putDevice(d)
		if err != nil {
			s.log.Error("failed to add device", "serialnumber", d.SerialNumber)
			br.Message = "Failed to add device " + d.SerialNumber
			err = s.updateLicense(l)
			if err != nil {
				s.log.Error("failed to update license", "err", err, "devices", adr)
			}
			return s.RenderJSON(req, br, http.StatusOK)
		}
		l.NumberDevices++
	}
	err = s.updateLicense(l)
	if err != nil {
		s.log.Error("failed to update license", "err", err, "devices", adr)
	}
	
	br.Status = true
	br.Message = "Devices added successfully"
	return s.RenderJSON(req, br, http.StatusOK)
}
// blockDevice godoc
// @Summary 			Block Device
// @Description 		Block a device by serial number.
// @Tags  				Devices
// @Produce 			json
// @Accept				json
// @Param 				data body DeviceRequest true "JSON object containing serial number of the device to block"
// @Success 			200 {object} DeviceRequest
// @Failure 			400 {object} ensweb.BaseResponse
// @Failure             500 {object} ensweb.BaseResponse
// @Router 				/api/blockdevice [post]
func (s *Server) blockDevice(req *ensweb.Request) *ensweb.Result {
	var devr DeviceRequest
	br := &ensweb.BaseResponse{
		Status: false,
	}
	err := s.ParseJSON(req, &devr)
	if err != nil {
		br.Message = "failed to block device, failed to unmarshal json"
		s.log.Error(br.Message)
		return s.RenderJSON(req, br, http.StatusInternalServerError)
	}
	d, err := s.getDevice(devr.SerialNumber)
	if err != nil {
		br.Message = fmt.Sprintf("Device (%s) not found", devr.SerialNumber)
		s.log.Error(br.Message)
		return s.RenderJSON(req, br, http.StatusOK)
	}
	d.IsBlocked = true
	err = s.putDevice(d)
	if err != nil {
		br.Message = "Failed to block device " + d.SerialNumber
		s.log.Error(br.Message)
		return s.RenderJSON(req, br, http.StatusOK)
	}
	br.Status = true
	br.Message = "Devices blocked successfully"
	return s.RenderJSON(req, br, http.StatusOK)
}
// unblockDevice godoc
// @Summary			 Unblock Device
// @Description		 Unblock a device by serial number.
// @Tags			 Devices
// @Produce			 json
// @Accept           json
// @Param 			 data body DeviceRequest true "JSON object containing serial number of the device to unblock"
// @Success			 200 {object} DeviceRequest
// @Failure          400 {object} ensweb.BaseResponse
// @Failure 		 500 {object} ensweb.BaseResponse
// @Router 			/api/unblockdevice [post]
func (s *Server) unblockDevice(req *ensweb.Request) *ensweb.Result {
	var devr DeviceRequest
	br := &ensweb.BaseResponse{
		Status: false,
	}
	err := s.ParseJSON(req, &devr)
	if err != nil {
		br.Message = "failed to unblock device, failed to unmarshal json"
		s.log.Error(br.Message)
		return s.RenderJSON(req, br, http.StatusInternalServerError)
	}
	d, err := s.getDevice(devr.SerialNumber)
	if err != nil {
		br.Message = fmt.Sprintf("Device (%s) not found", devr.SerialNumber)
		s.log.Error(br.Message)
		return s.RenderJSON(req, br, http.StatusOK)
	}
	d.IsBlocked = false
	err = s.putDevice(d)
	if err != nil {
		
		s.log.Error("failed to unblock device", "serialnumber", d.SerialNumber)
		br.Message = "Failed to unblock device " + d.SerialNumber
		return s.RenderJSON(req, br, http.StatusOK)
	}
	br.Status = true
	br.Message = "Devices unblocked successfully"
	return s.RenderJSON(req, br, http.StatusOK)
}
// deviceStatus godoc
// @Summary 			Get Device Status
// @Description 		Get the status of a device by serial number.
// @Produce 			json
// @Accept 				json
// @Param 				data body DeviceRequest true "JSON object containing serial number of the device to check status"
// @Success 			200 {object} DeviceRequest
// @Failure  			400 {object} DeviceStatus
// @Failure				500 {object} DeviceStatus
// @Router 				/api/devicestatus [post]
func (s *Server) deviceStatus(req *ensweb.Request) *ensweb.Result {
	var devr DeviceRequest
	ds := &DeviceStatus{
		SerialNumber: devr.SerialNumber,
	}
	err := s.ParseJSON(req, &devr)
	if err != nil {
		ds.Message = "failed to get device status, failed to unmarshal json"
		return s.RenderJSON(req, ds, http.StatusInternalServerError)
	}
	d, err := s.getDevice(devr.SerialNumber)
	if err != nil {
		ds.Message = fmt.Sprintf("Device (%s) not found", devr.SerialNumber)
		return s.RenderJSON(req, ds, http.StatusOK)
	}
	ds.IsBlocked = d.IsBlocked
	ds.Status = true
	ds.Message = "Got device status successfully"
	return s.RenderJSON(req, ds, http.StatusOK)
}
