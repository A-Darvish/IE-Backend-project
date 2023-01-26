package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arvnd79/ie-proj/common"
	"github.com/arvnd79/ie-proj/model"
	"github.com/labstack/echo/v4"
)

// TODO: add pagination support
// FetchURLs is used to retrieve a user's urls
// accessible with GET /api/urls
func (h *Handler) FetchURLs(c echo.Context) error {
	userID := extractID(c)
	urls, err := h.st.GetURLsByUser(userID)
	if err != nil {
		return common.NewRequestError("Error retrieving urls from database, maybe check your token again", err, http.StatusBadRequest)
	}
	resp := newURLListResponse(urls)
	return c.JSON(http.StatusOK, NewResponseData(resp))
}

// CreateURL is used to add a url to monitor service
// urls are validated and if there isn't any error a response code 201 is returned
// json request format:
//
//	{
//		"address": "http://google.com",
//		"threshold": 10
//	}
func (h *Handler) CreateURL(c echo.Context) error {
	userID := extractID(c)
	req := &urlCreateRequest{}
	url := &model.URL{}

	if err := req.bind(c, url); err != nil {
		return err
	}
	url.UserId = userID
	// adding url to database
	if err := h.st.AddURL(url); err != nil {
		// internal error
		return common.NewRequestError("error adding url to database", err, http.StatusInternalServerError)
	}
	// adding url to monitor scheduler
	h.sch.W.AddURL([]model.URL{*url})
	return c.JSON(http.StatusCreated, NewResponseData("URL created successfully"))
}

// DismissAlert updates a url inside database, resetting it's "failed_times" to 0
// returns an error in case of bad format request or invalid url_id
// json request format:
//
//{
//	"url_id": id
//}

// GetURLStats reports stats of a url
// returns error in case of invalid url_id or unauthenticated request
// param request format :
//
// /api/urls/:urlID
// you can also specify time intervals to get stats in
// just use unix timestamp with the syntax below (to_time is optional):
// /api/urls/:urlID?from_time=1579184689[&to_time]
func (h *Handler) GetURLStats(c echo.Context) error {
	userID := extractID(c)
	urlID, err := strconv.Atoi(c.Param("urlID"))
	if err != nil {
		return common.NewRequestError("Invalid path parameter", err, http.StatusBadRequest)
	}

	url := new(model.URL)

	url, err = h.st.GetURLById(uint(urlID))
	if err != nil {
		return common.NewRequestError("error retrieving url stats, invalid url id", err, http.StatusBadRequest)
	}
	if url.UserId != userID {
		return common.NewRequestError("operation not permitted", errors.New("user is not the owner of url"), http.StatusUnauthorized)
	}
	return c.JSON(http.StatusOK, NewResponseData(newRequestListResponse(url.Requests, url.Address)))
}

// DeleteURL deletes a url with given id
// returns error if url_id is invalid or user can't modify this url
// request format :
//
// DELETE /api/urls/:urlID
