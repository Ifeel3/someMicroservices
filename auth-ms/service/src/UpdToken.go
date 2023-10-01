package src

import (
	. "auth-ms/src/structs"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdToken(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
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
		date := time.Now().Unix()
		token := fmt.Sprintf("%x", sha256.Sum256([]byte(login.Pass+time.Now().GoString())))
		err = conn.QueryRow(context.TODO(), "update auth set token=$1, date=$2 where login=$3 returning token, date", token, date, login.Login).Scan(&token, &date)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(TokenStruct{Token: token, Date: date})
		}
	}
}
