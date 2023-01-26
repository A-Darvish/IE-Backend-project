package main

import (
	"log"
	"time"

	"github.com/arvnd79/ie-proj/common"
	"github.com/arvnd79/ie-proj/db"
	"github.com/arvnd79/ie-proj/endpointWatch"
	"github.com/arvnd79/ie-proj/handler"
	"github.com/arvnd79/ie-proj/store"
	"github.com/labstack/echo/v4"
)

func main() {
	our_db := db.Setup("http-endpoint-watch.db")
	st := store.NewStore(our_db)
	w := endpointWatch.NewWatch(st, nil, 10)

	sch, _ := endpointWatch.NewScheduler(w)
	sch.DoWithIntervals(time.Minute * 1)

	err := w.LoadFromDatabase()
	if err != nil {
		log.Println(err)
	}
	e := echo.New()

	// Group creates a new router group with prefix and optional group-level middleware.
	g := e.Group("/api")
	h := handler.NewHandler(st, sch)
	h.RegisterRoutes(g)

	e.HTTPErrorHandler = common.CustomHTTPErrorHandler

	e.Logger.Fatal(e.Start(":8080"))

}
