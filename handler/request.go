package handler

import (
	"net/http"

	"github.com/arvnd79/ie-proj/common"
	"github.com/arvnd79/ie-proj/model"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

type userAuthRequest struct {
	Username string `valid:"stringlength(3|12), alphanum" json:"username"`
	Password string `valid:"stringlength(3|12)" json:"password"`
}

// binding user auth request with model.User instance
func (r *userAuthRequest) bind(c echo.Context, user *model.User) error {
	if err := c.Bind(r); err != nil {
		return common.NewRequestError("error binding user request", err, http.StatusBadRequest)
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		e := common.NewValidationError(err, "Error validating sign-up request")
		return e
	}
	user.Username = r.Username
	user.Password = r.Password
	return nil
}

type urlCreateRequest struct {
	Address   string `json:"address" valid:"url"`
	Threshold int    `json:"threshold" valid:"int"`
}

func (r *urlCreateRequest) bind(c echo.Context, url *model.URL) error {
	if err := c.Bind(r); err != nil {
		return common.NewRequestError("error binding url create request, check json structure and try again", err, http.StatusBadRequest)
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		e := common.NewValidationError(err, "Error validating create url request")
		return e
	}
	url.Address = r.Address
	url.Threshold = r.Threshold
	url.FailedTimes = 0
	return nil
}

type alertDismissRequest struct {
	URLID uint `json:"url_id"`
}
