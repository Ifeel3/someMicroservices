package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

type Login struct {
	Status string
	Login  string `json:"login"`
	Pass   string `json:"pass"`
	Token  string `json:"token"`
	Date   int64  `json:"date"`
}

var addr string = "postgres://postgres:pass1234@db:5432/postgres"

func addLogin(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), addr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Database not connected", "", "", 0})
		return
	}
	defer conn.Close(context.TODO())
	var login Login
	err = json.NewDecoder(r.Body).Decode(&login)
	if err != nil || login.Login == "" || login.Pass == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Bad request", "", "", 0})
		return
	}
	var tmp Login
	err = conn.QueryRow(context.TODO(), "select login, pass from auth where login=$1", login.Login).Scan(&tmp.Login, &tmp.Pass)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Login already created", "", "", 0})
		return
	}
	login.Date = time.Now().Unix()
	login.Token = fmt.Sprintf("%x", sha256.Sum256([]byte(login.Pass+time.Now().GoString())))
	err = conn.QueryRow(context.TODO(), "insert into auth (login, pass, token, date) values ($1, $2, $3, $4) returning login, pass, token, date", login.Login, login.Pass, login.Token, login.Date).Scan(&login.Login, &login.Pass, &login.Token, &login.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Login login not created", "", "", 0})
		return
	} else {
		w.WriteHeader(http.StatusOK)
		login.Status = "OK"
		json.NewEncoder(w).Encode(login)
		return
	}
}

func checkToken(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), addr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Database not connected", "", "", 0})
		return
	}
	defer conn.Close(context.TODO())
	var login Login
	err = json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Bad request", "", "", 0})
		return
	}
	var tmp Login
	err = conn.QueryRow(context.TODO(), "select login, token, date from auth where login=$1", login.Login).Scan(&tmp.Login, &tmp.Token, &tmp.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Bad request", "", "", 0})
		return
	} else if login.Login != tmp.Login {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Wrong token", "", "", 0})
		return
	}
	t := time.Unix(tmp.Date, 0).Add(time.Hour).Unix()
	current := time.Now().Unix()
	if current > t {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Expired token", "", "", 0})
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"OK", "", "", "", 0})
	}
}

func delLogin(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), addr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Database not connected", "", "", 0})
		return
	}
	defer conn.Close(context.TODO())
	var login Login
	err = json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Bad request", "", "", 0})
		return
	}
	var tmp Login
	err = conn.QueryRow(context.TODO(), "select login, pass from auth where login=$1", login.Login).Scan(&tmp.Login, &tmp.Pass)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Bad request", "", "", 0})
		return
	} else if login.Login != tmp.Login || login.Pass != tmp.Pass {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Wrong login or pass", "", "", 0})
		return
	} else {
		err = conn.QueryRow(context.TODO(), "delete from auth where login=$1 returning login, pass, token, date", login.Login).Scan(&login.Login, &login.Pass, &login.Token, &login.Date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Login{"Error", "Bad request", "", "", 0})
			return
		}
		w.WriteHeader(http.StatusOK)
		login.Status = "OK"
		json.NewEncoder(w).Encode(login)
	}
}

func updToken(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), addr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Database not connected", "", "", 0})
		return
	}
	defer conn.Close(context.TODO())
	var login Login
	err = json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Bad request", "", "", 0})
		return
	}
	var tmp Login
	err = conn.QueryRow(context.TODO(), "select login, pass from auth where login=$1", login.Login).Scan(&tmp.Login, &tmp.Pass)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Bad request", "", "", 0})
		return
	} else if login.Login != tmp.Login || login.Pass != tmp.Pass {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Wrong login or pass", "", "", 0})
		return
	} else {
		login.Date = time.Now().Unix()
		login.Token = fmt.Sprintf("%x", sha256.Sum256([]byte(login.Pass+time.Now().GoString())))
		err = conn.QueryRow(context.TODO(), "update auth set token=$1, date=$2 where login=$3 returning login, token, date", login.Token, login.Date, login.Login).Scan(&login.Login, &login.Token, &login.Date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Login{"Error", "Bad request", "", "", 0})
			return
		} else {
			w.WriteHeader(http.StatusAccepted)
			login.Status = "OK"
			json.NewEncoder(w).Encode(login)
		}
	}
}

func updPass(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), addr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Database not connected", "", "", 0})
		return
	}
	defer conn.Close(context.TODO())
	var login struct{ Login, Pass, NewPass string }
	err = json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Bad request", "", "", 0})
		return
	}
	var tmp Login
	err = conn.QueryRow(context.TODO(), "select login, pass from auth where login=$1", login.Login).Scan(&tmp.Login, &tmp.Pass)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Bad request", "", "", 0})
		return
	} else if login.Login != tmp.Login || login.Pass != tmp.Pass {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Login{"Error", "Wrong login or pass", "", "", 0})
		return
	} else {
		date := time.Now().Unix()
		token := fmt.Sprintf("%x", sha256.Sum256([]byte(login.Pass+time.Now().GoString())))
		err = conn.QueryRow(context.TODO(), "update auth set pass=$1, token=$2, date=$3 where login=$4 returning login, pass, token, date", login.NewPass, token, date, login.Login).Scan(&login.Login, &login.Pass, &token, &date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Login{"Error", "Bad request", "", "", 0})
			return
		} else {
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(Login{"OK", login.Login, login.Pass, token, date})
		}
	}
}

func main() {
	http.HandleFunc("/api/add", addLogin)
	http.HandleFunc("/api/check", checkToken)
	http.HandleFunc("/api/del", delLogin)
	http.HandleFunc("/api/upd/token", updToken)
	http.HandleFunc("/api/upd/pass", updPass)
	http.ListenAndServe(":8081", nil)
}
