package server

import (
	"net/http"
	"time"

	"github.com/EnsurityTechnologies/enscrypt"
	"github.com/EnsurityTechnologies/ensweb"
	"github.com/golang-jwt/jwt/v5"
)

// admin Login godoc
// @Summary				Admin Login
// @Description			Authenticate an administrator.
// @Produce 			json
// @Accept   			json
// @Param 				data body LoginRequest true "JSON object containing login credentials"
// @Success				200 {object} LoginResponse
// @Failure  			400 {object} ensweb.BaseResponse
// @Failure 			401 {object} ensweb.BaseResponse
// @Failure 			500 {object} ensweb.BaseResponse
// @Router 				/api/adminlogin [post]
func (s *Server) adminLogin(req *ensweb.Request) *ensweb.Result {
	var lr LoginRequest
	lrp := &LoginResponse{
		BaseResponse: ensweb.BaseResponse{
			Status: false,
		},
	}
	err := s.ParseJSON(req, &lr)
	if err != nil {
		lrp.Message = "invalid login request, failed to unmarshal json"
		return s.RenderJSON(req, lrp, http.StatusBadRequest)
	}
	u, err := s.getUser(lr.UserName)
	if err != nil {
		s.log.Error("user does not exist", "err", err, "username", u.UserName)
		lrp.Message = "user does not exist"
		return s.RenderJSON(req, lrp, http.StatusUnauthorized)
	}
	if u.Role != AdminRole {
		s.log.Error("unauthorized user", "role", u.Role, "username", u.UserName)
		lrp.Message = "unauthorized user"
		return s.RenderJSON(req, lrp, http.StatusUnauthorized)
	}
	if !enscrypt.VerifyPassword(lr.Password, u.Password) {
		s.log.Error("password does not match", "username", u.UserName)
		lrp.Message = "password does not match"
		return s.RenderJSON(req, lrp, http.StatusUnauthorized)
	}
	expiresAt := time.Now().Add(time.Minute * 10)
	bt := &BearerToken{
		UserName: u.UserName,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "authserver",
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	tk := s.GenerateJWTToken(bt)
	lrp.Status = true
	lrp.Message = "user logged in successfully"
	lrp.Token = tk
	return s.RenderJSON(req, lrp, http.StatusOK)
}

// userLogin godoc
// @Summary 			User Login
// @Description			 Authenticate and log in a user.
// @Produce 			json
// @Accept 				json
// @Param 				data body UserRequest true "JSON object containing user login details"
// @Success 			200 {object} UserResponse
// @Failure				400 {object} ensweb.BaseResponse
// @Failure 			401 {object} ensweb.UserResponse
// @Failure 			500 {object} ensweb.UserResponse
// @Router 				/api/userlogin [post]
func (s *Server) userLogin(req *ensweb.Request) *ensweb.Result {
	var lr UserRequest
	lrp := &UserResponse{
		BaseResponse: ensweb.BaseResponse{
			Status: false,
		},
	}
	err := s.ParseJSON(req, &lr)
	if err != nil {
		s.log.Error("invalid json input", "err", err)
		lrp.Message = "invalid login request, failed to unmarshal json"
		return s.RenderJSON(req, lrp, http.StatusBadRequest)
	}
	d, err := s.getDevice(lr.SerialNumber)
	if err != nil || d.SerialNumber != lr.SerialNumber {
		s.log.Error("invalid device", "sn", lr.SerialNumber)
		lrp.Message = "Invalid device"
		return s.RenderJSON(req, lrp, http.StatusUnauthorized)
	}
	if d.IsBlocked {
		s.log.Error("device is blocked", "sn", lr.SerialNumber)
		lrp.Message = "Device is blocked"
		return s.RenderJSON(req, lrp, http.StatusUnauthorized)
	}
	u, err := s.getUser(lr.UserName)
	if err != nil {
		l, err := s.getLicense()
		if err != nil || l.MaxNumberUsers == 0 {
			s.log.Error("failed to get license", "err", err, "username", u.UserName)
			lrp.Message = "Failed to get license"
			return s.RenderJSON(req, lrp, http.StatusUnauthorized)
		}
		defer s.releaseLicense()
		if l.MaxNumberUsers == l.NumberUsers {
			s.log.Error("user license already reached maximum")
			lrp.Message = "User license already reached maximum"
			return s.RenderJSON(req, lrp, http.StatusUnauthorized)
		}
		u = &User{
			UserName:        lr.UserName,
			UserID:          lr.UID,
			Role:            UserRole,
			IsPinRequired:   false,
			Pin:             RandString(16),
			IsActive:        true,
			AllowedFingers:  MaxAllowedFingers,
			EnrolledFingers: 0,
			KeyAdded:        false,
		}
		err = s.putUser(u)
		if err != nil {
			s.log.Error("failed to create user", "username", lr.UserName)
			lrp.Message = "Failed to create user"
			return s.RenderJSON(req, lrp, http.StatusInternalServerError)
		}
		l.NumberUsers++
		err = s.updateLicense(l)
		if err != nil {
			s.log.Error("failed to update license", "username", lr.UserName)
			lrp.Message = "Failed to update license"
			return s.RenderJSON(req, lrp, http.StatusInternalServerError)
		}
	}
	if u.UserID != lr.UID {
		s.log.Error("user id does not match", "err", err, "username", u.UserName)
		lrp.Message = "user id does not match"
		return s.RenderJSON(req, lrp, http.StatusInternalServerError)
	}
	if !u.IsActive {
		s.log.Error("user is not active", "username", u.UserName)
		lrp.Message = "user is not active"
		return s.RenderJSON(req, lrp, http.StatusUnauthorized)
	}

	expiresAt := time.Now().Add(time.Minute * 10)
	bt := &BearerToken{
		UserName: u.UserName,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "authserver",
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	tk := s.GenerateJWTToken(bt)
	lrp.Status = true
	lrp.Message = "user logged in successfully"
	lrp.Token = tk
	us := UserStatus{
		IsPinRequired:   u.IsPinRequired,
		Pin:             u.Pin,
		AllowedFingers:  u.AllowedFingers,
		EnrolledFingers: u.EnrolledFingers,
	}
	lrp.UserStatus = us
	return s.RenderJSON(req, lrp, http.StatusOK)
}
