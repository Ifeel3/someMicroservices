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
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Bad request"}})
		return
	}
	var check LoginStruct
	err = conn.QueryRow(context.TODO(), "select login, pass from auth where login=$1", login.Login).Scan(&check.Login, &check.Pass)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Login already created"}})
		return
	}
	err = conn.QueryRow(context.TODO(), "insert into auth (login, pass) values ($1, $2) returning login, pass", login.Login, login.Pass).Scan(&login.Login, &login.Pass)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Login not created"}})
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "OK"}})
		return
	}
}
