package main

import (
	. "auth-ms/src"
	. "auth-ms/src/structs"
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

var addr string = "postgres://postgres:pass1234@localhost:5432/postgres"

var conn *pgxpool.Pool
var connErr error

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if connErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Database not connected"}})
		return
	}
	switch r.Method {
	case "POST":
		AddLogin(w, r, conn)
	case "PATCH":
		UpdPass(w, r, conn)
	case "DELETE":
		DelLogin(w, r, conn)
	case "GET":
		CheckAuth(w, r, conn)
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Bad request"}})
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	if connErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Database not connected"}})
		return
	}
	switch r.Method {
	case "PATCH":
		UpdToken(w, r, conn)
	case "GET":
		CheckToken(w, r, conn)
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Bad request"}})
	}
}

func main() {
	conn, connErr = pgxpool.New(context.TODO(), addr)
	if connErr == nil {
		defer conn.Close()
	}
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/login/", loginHandler)
	http.HandleFunc("/api/token", tokenHandler)
	http.HandleFunc("/api/token/", tokenHandler)
	http.ListenAndServe(":8081", nil)
}
