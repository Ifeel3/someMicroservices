package src

import (
	. "auth-ms/src/structs"
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DelLogin(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	var login LoginStruct
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var tmp LoginStruct
	err = conn.QueryRow(context.TODO(), "select login, pass from auth where login=$1", login.Login).Scan(&tmp.Login, &tmp.Pass)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if login.Login != tmp.Login || login.Pass != tmp.Pass {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Wrong login or pass")
		return
	} else {
		err = conn.QueryRow(context.TODO(), "delete from auth where login=$1 returning login, pass", login.Login).Scan(&login.Login, &login.Pass)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
