package api

import (
	"app/middleware"
	"database/sql"
	"encoding/json"
	"net/http"

	"app/data_access/db"
	"app/util"
)

type AppInjector struct {
	db *sql.DB
}

type App struct {
	AppId   int
	AppName string
}

func Index(w http.ResponseWriter, r *http.Request) {
	entry := middleware.GetAppLogger(r)
	entry.Println("test")
	panic("test panic error")

}

func Data(w http.ResponseWriter, r *http.Request) {

	

}
