package handler

import (
	"database/sql"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func ReadyHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if err := dbConn.Ping(); err != nil {
			http.Error(w, "db not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	}
}
