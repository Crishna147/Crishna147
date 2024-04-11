package server

import (
	"time"

	"github.com/EnsurityTechnologies/ensweb"
)

const (
	UserIDTable  string = "usersid"
	UserTable    string = "users"
	LicenseTable string = "license"
	DeviceTable  string = "devices"
)

const (
	AdminRole string = "admin"
	UserRole  string = "user"
)

type License struct {
	MaxNumberUsers   int       `json:"max_number_users"`
	NumberUsers      int       `json:"number_users"`
	MaxNumberDevices int       `json:"max_number_devices"`
	NumberDevices    int       `json:"number_devices"`
	ExpiryTime       time.Time `json:"expiry_time"`
}

type UserID struct {
	UserID string `json:"uid"`
}

type User struct {
	UserName        string `json:"username"`
	UserID          string `json:"uid"`
	Role            string `json:"role"`
	Password        string `json:"password"`
	IsPinRequired   bool   `json:"is_pin_required"`
	Pin             string `json:"pin"`
	IsActive        bool   `json:"is_active"`
	AllowedFingers  int    `json:"allowed_fingers"`
	EnrolledFingers int    `json:"enrolled_fingers"`
	KeyAdded        bool   `json:"key_added"`
}

type Device struct {
	SerialNumber string `json:"serial_number"`
	PublicKey    string `json:"pub_key"`
	IsBlocked    bool   `json:"is_blocked"`
}

type RequestType struct {
	ID        string `json:"uuid"`
	JourneyID string `json:"journeyId"`
	TS        int64  `json:"ts"`
	AppID     string `json:"appid"`
}

type DataRequest struct {
	Data   string
	Secret string
}

type CommonData struct {
	Data string `json:"data"`
}

type PublicKeyResponse struct {
	PublicKey string `json:"publicKey"`
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserRequest struct {
	UserName      string `json:"username"`
	UID           string `json:"uid"`
	SerialNumber  string `json:"serialNumber"`
	MaxNumFingers int    `json:"maxNumFingers"`
	IsPinRequired bool   `json:"isPinRequired"`
}

type LoginResponse struct {
	ensweb.BaseResponse
	Token string `json:"token"`
}

type UserStatus struct {
	IsPinRequired   bool   `json:"isPinRequired"`
	Pin             string `json:"pin"`
	AllowedFingers  int    `json:"allowedFingers"`
	EnrolledFingers int    `json:"enrolledFingers"`
}

type UserResponse struct {
	ensweb.BaseResponse
	Token      string     `json:"token"`
	UserStatus UserStatus `json:"userStatus"`
}

type AddDeviceRequest struct {
	SerialNumber string `json:"serialNumber"`
	PublicKey    string `json:"publicKey"`
}

type DeviceRequest struct {
	SerialNumber string `json:"serialNumber"`
}

type DeviceStatus struct {
	ensweb.BaseResponse
	SerialNumber string `json:"serialNumber"`
	IsBlocked    bool   `json:"isBlocked"`
}
