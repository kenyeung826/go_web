package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"app/data_access/db"
	appMiddleware "app/middleware"
	"app/util"
)

const (
	port         = ":8080"	
)

func main() {
	util.LoadConfig()

	defer func() {
		fmt.Println("defer Close")
		db.CloseAll()
		util.CloseGlobalLog()
	}()
	r := chi.NewRouter()
	//r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(appMiddleware.ServerLogger)

	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(appMiddleware.ErrorHandler)

	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.ThrottleBacklog(1000, 2000, 20*time.Second))
	r.Use(appMiddleware.BufferedResponseHandler)
	Route(r)
	http.ListenAndServe(port, r)

	// Gracefully shut down the server

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		db.CloseAll()
		util.CloseGlobalLog()
		fmt.Println("close Database")
	}()
	fmt.Println("End of Service?")

}
