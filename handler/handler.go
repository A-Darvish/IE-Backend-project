package handler

import (
	"github.com/arvnd79/ie-proj/endpointWatch"
	"github.com/arvnd79/ie-proj/store"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	st  *store.Store
	sch *endpointWatch.Scheduler
}

// NewHandler creates a new handler with given store instance
func NewHandler(st *store.Store, sch *endpointWatch.Scheduler) *Handler {
	return &Handler{st: st, sch: sch}
}

func extractID(c echo.Context) uint {
	e := c.Get("user").(*jwt.Token)    // Get retrieves data from the context.
	claims := e.Claims.(jwt.MapClaims) //payload
	id := uint(claims["id"].(float64))
	return id
}
