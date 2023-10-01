package src

import (
	. "auth-ms/src/structs"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckToken(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	splitted := strings.Split(r.RequestURI, "/")
	if len(splitted) < 4 || splitted[3] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	login := splitted[3]
	var token TokenStruct
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var tmpToken TokenStruct
	err = conn.QueryRow(context.TODO(), "select login, token, date from auth where login=$1", login).Scan(&login, &tmpToken.Token, &tmpToken.Date)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	t := time.Unix(tmpToken.Date, 0).Add(time.Hour).Unix()
	current := time.Now().Unix()
	if token.Token != tmpToken.Token {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Token not valid")
		return
	} else if current > t {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Token expired")
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
