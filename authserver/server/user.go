package server

import (
	"net/http"

	"github.com/EnsurityTechnologies/ensweb"
)

func (s *Server) getUser(username string) (*User, error) {
	var uid UserID
	err := s.ts.Read(UserIDTable, username, &uid)
	if err != nil {
		return nil, err
	}
	var u User
	err = s.ts.Read(UserTable, uid.UserID, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Server) getUserByID(uid string) (*User, error) {
	if uid == "0000000000000000000000000000000000000000000000000000000000000000" || uid == "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" {
		return &User{UserName: "default", IsActive: true}, nil
	}
	us, ok := s.users.Get(uid)
	if ok {
		us := us.(User)
		return &us, nil
	} else {
		var u User
		err := s.ts.Read(UserTable, uid, &u)
		if err != nil {
			return nil, err
		}
		s.users.Set(uid, u, 0)
		return &u, nil
	}

}

func (s *Server) putUser(u *User) error {
	uid := UserID{
		UserID: u.UserID,
	}
	err := s.ts.Write(UserIDTable, u.UserName, uid)
	if err != nil {
		return err
	}
	err = s.ts.Write(UserTable, u.UserID, u)
	if err != nil {
		return err
	}
	s.users.Set(uid.UserID, *u, 0)
	return nil
}
// fingerEnrolled godoc
// @Summary			 User FingerEnrolled Status
// @Description 	 Update the finger enrolled status of a user.
// @Produce			 json
// @Accept 			 json
// @Param 			 data body UserRequest true "JSON object containing user details"
// @Success 		 200 {object} UserResponse
// @Failure 		 400 {object} ensweb.BaseResponse
// @Failure			 500 {object} ensweb.BaseResponse
// @Router 			 /api/fingerenrolled [post]
func (s *Server) fingerEnrolled(req *ensweb.Request) *ensweb.Result {
	var ur UserRequest
	urp := &UserResponse{
		BaseResponse: ensweb.BaseResponse{
			Status: false,
		},
	}
	err := s.ParseJSON(req, &ur)
	if err != nil {
		urp.Message = "failed to add finger enroll, failed to unmarshal json"
		s.log.Error(urp.Message)
		return s.RenderJSON(req, urp, http.StatusBadRequest)
	}
	u, err := s.getUser(ur.UserName)
	if err != nil {
		s.log.Error("failed to get user", "err", err, "username", u.UserName)
		urp.Message = "User does not exist or invalid user"
		return s.RenderJSON(req, urp, http.StatusUnauthorized)
	}
	u.EnrolledFingers++
	err = s.putUser(u)
	if err != nil {
		s.log.Error("failed to update user", "err", err, "username", u.UserName)
		urp.Message = "Failed to update user"
		return s.RenderJSON(req, urp, http.StatusInternalServerError)
	}
	urp.Status = true
	urp.Message = "User finger enrolled status updated"
	us := UserStatus{
		IsPinRequired:   u.IsPinRequired,
		AllowedFingers:  u.AllowedFingers,
		EnrolledFingers: u.EnrolledFingers,
	}
	urp.UserStatus = us
	return s.RenderJSON(req, urp, http.StatusOK)
}
// keyAdded godoc
// @Summary 		Update User Registration Status
// @Description 	Update the user registration status to indicate whether the key is added.
// @Produce 		json
// @Accept 			json
// @Param 			data body UserRequest true "JSON object containing user details"
// @Success			200 {object} UserResponse
// @Failure			400 {object} ensweb.BaseResponse
// @Failure			500 {object} ensweb.BaseResponse
// @Router 			/api/keyadded [post]
func (s *Server) keyAdded(req *ensweb.Request) *ensweb.Result {
	var ur UserRequest
	urp := &UserResponse{
		BaseResponse: ensweb.BaseResponse{
			Status: false,
		},
	}
	err := s.ParseJSON(req, &ur)
	if err != nil {
		urp.Message = "failed to add key, failed to unmarshal json"
		s.log.Error(urp.Message)
		return s.RenderJSON(req, urp, http.StatusBadRequest)
	}
	u, err := s.getUser(ur.UserName)
	if err != nil {
		s.log.Error("failed to get user", "err", err, "username", u.UserName)
		urp.Message = "User does not exist or invalid user"
		return s.RenderJSON(req, urp, http.StatusUnauthorized)
	}
	u.KeyAdded = true
	err = s.putUser(u)
	if err != nil {
		s.log.Error("failed to update user", "err", err, "username", u.UserName)
		urp.Message = "Failed to update user"
		return s.RenderJSON(req, urp, http.StatusInternalServerError)
	}
	urp.Status = true
	urp.Message = "User registration status updated"
	us := UserStatus{
		IsPinRequired:   u.IsPinRequired,
		AllowedFingers:  u.AllowedFingers,
		EnrolledFingers: u.EnrolledFingers,
	}
	urp.UserStatus = us
	return s.RenderJSON(req, urp, http.StatusOK)
}
// reset User godoc
// @Summary 		Reset User
// @Description		Reset a user by removing their data and updating enrollment status.
// @Produce 		json
// @Accept			json
// @Param 			data body UserRequest true "JSON object containing user details"
// @Success			200 {object} UserResponse
// @Failure         400 {object} ensweb.BaseResponse
// @Failure			500 {object} ensweb.BaseResponse
// @Router 			/api/resetuser [post]
func (s *Server) resetUser(req *ensweb.Request) *ensweb.Result {
	var ur UserRequest
	urp := &UserResponse{
		BaseResponse: ensweb.BaseResponse{
			Status: false,
		},
	}
	err := s.ParseJSON(req, &ur)
	if err != nil {
		urp.Message = "failed to reset user, failed to unmarshal json"
		s.log.Error(urp.Message)
		return s.RenderJSON(req, urp, http.StatusBadRequest)
	}
	u, err := s.getUser(ur.UserName)
	if err != nil {
		s.log.Error("failed to get user", "err", err, "username", u.UserName)
		urp.Message = "User does not exist or invalid user"
		return s.RenderJSON(req, urp, http.StatusInternalServerError)
	}
	err = s.RemoveUser(u.UserID)
	if err != nil {
		s.log.Error("failed to reset user", "err", err, "username", u.UserName)
		urp.Message = "Failed to reset user"
		return s.RenderJSON(req, urp, http.StatusInternalServerError)
	}
	u.EnrolledFingers = 0
	u.KeyAdded = false
	err = s.putUser(u)
	if err != nil {
		s.log.Error("failed to update user", "err", err, "username", u.UserName)
		urp.Message = "Failed to update user"
		return s.RenderJSON(req, urp, http.StatusInternalServerError)
	}
	urp.Status = true
	urp.Message = "User reset successfully"
	us := UserStatus{
		IsPinRequired:   u.IsPinRequired,
		AllowedFingers:  u.AllowedFingers,
		EnrolledFingers: u.EnrolledFingers,
	}
	urp.UserStatus = us
	return s.RenderJSON(req, urp, http.StatusOK)
}
// setNumFingers godoc
// @Summary 			SetNumFingers
// @Description 		Set the maximum number of fingers allowed for a user.
// @Produce 			json
// @Accept 				json
// @Param 				data body UserRequest true "JSON object containing user details"
// @Success 			200 {object} UserResponse
// @Failure 			400 {object} ensweb.BaseResponse
// @Failure 			500 {object} ensweb.BaseResponse
// @Router 				/api/setnumfingers [post]
func (s *Server) setNumFingers(req *ensweb.Request) *ensweb.Result {
	var ur UserRequest
	urp := &UserResponse{
		BaseResponse: ensweb.BaseResponse{
			Status: false,
		},
	}
	err := s.ParseJSON(req, &ur)
	if err != nil {
		urp.Message = "failed to set number of fingers, failed to unmarshal json"
		s.log.Error(urp.Message)
		return s.RenderJSON(req, urp, http.StatusBadRequest)
	}
	u, err := s.getUser(ur.UserName)
	if err != nil {
		s.log.Error("failed to get user", "err", err, "username", u.UserName)
		urp.Message = "User does not exist or invalid user"
		return s.RenderJSON(req, urp, http.StatusInternalServerError)
	}
	if u.EnrolledFingers > ur.MaxNumFingers {
		s.log.Error("Enrolled fingers are more", "err", err, "username", u.UserName)
		urp.Message = "Enrolled fingers are more"
		return s.RenderJSON(req, urp, http.StatusUnauthorized)
	}
	u.AllowedFingers = ur.MaxNumFingers
	err = s.putUser(u)
	if err != nil {
		s.log.Error("failed to update user", "err", err, "username", u.UserName)
		urp.Message = "Failed to update user"
		return s.RenderJSON(req, urp, http.StatusInternalServerError)
	}
	urp.Status = true
	urp.Message = "Maximum fingers updated successfully"
	us := UserStatus{
		IsPinRequired:   u.IsPinRequired,
		AllowedFingers:  u.AllowedFingers,
		EnrolledFingers: u.EnrolledFingers,
	}
	urp.UserStatus = us
	return s.RenderJSON(req, urp, http.StatusOK)
}
// setPinRequired godoc
// @Summary			SetPinRequired
// @Description 	Set whether PIN is required for a user.
// @Produce 		json
// @Accept 			json
// @Param			data body UserRequest true "JSON object containing user details"
// @Success 		200 {object} UserResponse
// @Failure 		400 {object} ensweb.BaseResponse
// @Failure 		500 {object} ensweb.BaseResponse
// @Router 			/api/setpinrequired [post]
func (s *Server) setPinRequired(req *ensweb.Request) *ensweb.Result {
	var ur UserRequest
	urp := &UserResponse{
		BaseResponse: ensweb.BaseResponse{
			Status: false,
		},
	}
	err := s.ParseJSON(req, &ur)
	if err != nil {
		urp.Message = "failed to set pin required, failed to unmarshal json"
		s.log.Error(urp.Message)
		return s.RenderJSON(req, urp, http.StatusInternalServerError)
	}
	u, err := s.getUser(ur.UserName)
	if err != nil {
		s.log.Error("failed to get user", "err", err, "username", u.UserName)
		urp.Message = "User does not exist or invalid user"
		return s.RenderJSON(req, urp, http.StatusInternalServerError)
	}
	if u.IsPinRequired == ur.IsPinRequired {
		urp.Message = "User PIN settings already set"
		return s.RenderJSON(req, urp, http.StatusOK)
	}
	s.log.Info("PIN settings requested, resetting the user")
	err = s.RemoveUser(u.UserID)
	if err != nil {
		s.log.Error("failed to reset user", "err", err, "username", u.UserName)
		urp.Message = "Failed to reset user"
		return s.RenderJSON(req, urp, http.StatusInternalServerError)
	}
	u.IsPinRequired = ur.IsPinRequired
	u.EnrolledFingers = 0
	u.KeyAdded = false
	err = s.putUser(u)
	if err != nil {
		s.log.Error("failed to update user", "err", err, "username", u.UserName)
		urp.Message = "Failed to update user"
		return s.RenderJSON(req, urp, http.StatusInternalServerError)
	}
	urp.Status = true
	urp.Message = "User PIN settings done successfully"
	us := UserStatus{
		IsPinRequired:   u.IsPinRequired,
		AllowedFingers:  u.AllowedFingers,
		EnrolledFingers: u.EnrolledFingers,
	}
	urp.UserStatus = us
	return s.RenderJSON(req, urp, http.StatusOK)
}
