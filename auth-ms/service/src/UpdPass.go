package src

import (
	. "auth-ms/src/structs"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdPass(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	splitted := strings.Split(r.RequestURI, "/")
	if len(splitted) < 4 || splitted[3] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Bad request")
		return
	}
	login := splitted[3]
	var newpass NewPassStruct
	err := json.NewDecoder(r.Body).Decode(&newpass)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Bad request")
		return
	}
	var currentpass LoginStruct
	err = conn.QueryRow(context.TODO(), "select login, pass from auth where login=$1", login).Scan(&currentpass.Login, &currentpass.Pass)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if newpass.Pass != currentpass.Pass {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Wrong login or pass")
		return
	} else {
		err = conn.QueryRow(context.TODO(), "update auth set pass=$1 where login=$2 returning login, pass", newpass.NewPass, login).Scan(&login, &newpass.Pass)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
