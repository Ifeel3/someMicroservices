package main

import (
	"context"
	"encoding/json"
	"net/http"
	. "some-ms/src"

	"github.com/jackc/pgx/v5/pgxpool"
)

var addr string = "postgres://postgres:pass1234@db:5432/postgres"

var conn *pgxpool.Pool
var connErr error

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	if connErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Database not connected")
		return
	}
	switch r.Method {
	case "GET":
		GetItems(w, r, conn)
	case "POST":
		AddItem(w, r, conn)
	case "PATCH":
		UpdItem(w, r, conn)
	case "DELETE":
		DelItem(w, r, conn)
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Bad request")
	}
}

func main() {
	conn, connErr = pgxpool.New(context.TODO(), addr)
	if connErr == nil {
		defer conn.Close()
	}
	http.HandleFunc("/api/items/", RequestHandler)
	http.HandleFunc("/api/items", RequestHandler)
	http.ListenAndServe(":8082", nil)
}
