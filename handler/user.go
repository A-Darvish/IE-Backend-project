package handler

import (
	"net/http"

	"github.com/arvnd79/ie-proj/common"
	"github.com/arvnd79/ie-proj/model"
	"github.com/labstack/echo/v4"
)

// SignUp user handler to handle sign up request, if successful it returns JWT with username followed by 201 status code
// json request format:
//
//	{
//		"username": "ali" [alpha numeric,12 > len > 2],
//		"password": "ah0" [12 > len > 2]
//	}
func (h *Handler) SignUp(c echo.Context) error {
	req := &userAuthRequest{}
	user := &model.User{}
	if err := req.bind(c, user); err != nil {
		return err
	}
	// saving user
	if err := h.st.AddUser(user); err != nil {
		return common.NewRequestError("could not save user in database", err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, NewResponseData(NewUserResponse(user)))
}

// Login handler function, if login is successful a response with JWT and username is returned followed by a 200 status code
// json request format:
//
//	{
//		"username": "foo" [alpha numeric, 12 > len > 2],
//		"password": "bar1" [12 > len > 2]
//	}
func (h *Handler) Login(c echo.Context) error {
	req := &userAuthRequest{}
	user := &model.User{}
	if err := req.bind(c, user); err != nil {
		return err
	}
	// retrieving user from database
	u, err := h.st.GetUserByUserName(user.Username)
	if err != nil || !u.ValidatePassword(user.Password) {
		return common.NewRequestError("Invalid username or password", err, http.StatusUnauthorized)
	}
	return c.JSON(http.StatusOK, NewResponseData(NewUserResponse(u)))
}

// FetchAlerts retrieves all alerts for the user, returns a list of urls with alert
func (h *Handler) FetchAlerts(c echo.Context) error {
	userID := extractID(c)
	alerts, err := h.st.FetchAlerts(userID)
	if err != nil {
		return common.NewRequestError("coult not get alerts from database", err, http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, NewResponseData(alerts))
}
