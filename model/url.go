package model

import (
	"errors"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
)

type URL struct {
	UrlId       uint   `gorm:"primary_key"`
	UserId      uint   `gorm:"unique_index:index_addr_user"` // for preventing url duplication for a single user
	Address     string `gorm:"unique_index:index_addr_user"`
	Threshold   int
	FailedTimes int
	Requests    []Request `gorm:"foreignkey:url_id"`
}

type Request struct {
	ReqId  uint `gorm:"primary_key"`
	UrlId  uint
	Result int
}

func NewURL(userID uint, address string, threshold int) (*URL, error) {
	url := new(URL)
	url.UserId = userID
	url.Threshold = threshold
	url.FailedTimes = 0

	isValid := govalidator.IsURL(address)
	if !strings.HasPrefix("http://", address) {
		address = "http://" + address
	}
	if isValid {
		url.Address = address
		return url, nil
	}
	return nil, errors.New("not a valid URL address")
}

func (url *URL) ShouldTriggerAlarm() bool {
	return url.FailedTimes >= url.Threshold
}

// SendRequest sends a HTTP GET request to the url
// returns a *Request with result status code
func (url *URL) SendRequest() (*Request, error) {
	resp, err := http.Get(url.Address)
	req := new(Request)
	req.UrlId = url.UrlId
	if err != nil {
		return req, err
	}
	req.Result = resp.StatusCode
	return req, nil
}
