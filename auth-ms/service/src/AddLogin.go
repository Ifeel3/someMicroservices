package src

import (
	. "auth-ms/src/structs"
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddLogin(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	var login LoginStruct
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil || login.Login == "" || login.Pass == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var check LoginStruct
	err = conn.QueryRow(context.TODO(), "select login, pass from auth where login=$1", login.Login).Scan(&check.Login, &check.Pass)
	if err == nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Login already created")
		return
	}
	err = conn.QueryRow(context.TODO(), "insert into auth (login, pass) values ($1, $2) returning login, pass", login.Login, login.Pass).Scan(&login.Login, &login.Pass)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Login not created")
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}
