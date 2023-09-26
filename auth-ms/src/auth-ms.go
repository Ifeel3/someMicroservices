package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Login struct {
	Id    int
	Login string
	Hash  string
}

type Status struct {
	Status string
	Str    string
}

func addLogin(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), "postgres://postgres:pass1234@db:5432/postgres")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Status{"Error", "Database not connected"})
		return
	}
	defer conn.Close(context.TODO())
	str := strings.Split(r.RequestURI, "/")
	if len(str) < 4 || len(str) > 4 || str[3] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Status{"Error", "Bad request"})
		return
	} else {
		var ans Login
		err := conn.QueryRow(context.TODO(), "select id, login, hash from auth where login=$1", str[3]).Scan(&ans.Id, &ans.Login, &ans.Hash)
		if err == nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Status{"Error", "Login already created"})
			return
		}
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(str[3]+time.Now().GoString())))
		_, err = conn.Exec(context.TODO(), "insert into auth (login, hash) values ($1,$2)", str[3], hash)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Status{"Error", "Wrong query"})
			return
		} else {
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(Status{"OK", hash})
		}
	}
}

func getLogin(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), "postgres://postgres:pass1234@db:5432/postgres")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Status{"Error", "Database not connected"})
		return
	}
	defer conn.Close(context.TODO())
	str := strings.Split(r.RequestURI, "/")
	if len(str) < 4 || len(str) > 4 || str[3] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Status{"Error", "Bad request"})
		return
	} else {
		var ans Login
		err := conn.QueryRow(context.TODO(), "select id, login, hash from auth where login=$1", str[3]).Scan(&ans.Id, &ans.Login, &ans.Hash)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Status{"Error", "Login not found"})
			return
		} else {
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(Status{"OK", ans.Hash})
		}
	}
}

func delLogin(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), "postgres://postgres:pass1234@db:5432/postgres")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Status{"Error", "Database not connected"})
		return
	}
	defer conn.Close(context.TODO())
	str := strings.Split(r.RequestURI, "/")
	if len(str) < 4 || len(str) > 4 || str[3] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Status{"Error", "Bad request"})
		return
	} else {
		_, err := conn.Exec(context.TODO(), "delete from auth where login=$1", str[3])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Status{"Error", "Not executed"})
			return
		} else {
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(Status{"OK", "Login deleted"})
		}
	}
}

func updLogin(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), "postgres://postgres:pass1234@db:5432/postgres")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Status{"Error", "Database not connected"})
		return
	}
	defer conn.Close(context.TODO())
	str := strings.Split(r.RequestURI, "/")
	if len(str) < 4 || len(str) > 4 || str[3] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Status{"Error", "Bad request"})
		return
	} else {
		var ans Login
		err := conn.QueryRow(context.TODO(), "select id, login, hash from auth where login=$1", str[3]).Scan(&ans.Id, &ans.Login, &ans.Hash)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Status{"Error", "Login not found"})
			return
		}
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(str[3]+time.Now().GoString())))
		_, err = conn.Exec(context.TODO(), "update auth set hash=$1 where login=$2", hash, str[3])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Status{"Error", "Not executed"})
			return
		} else {
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(Status{"OK", hash})
		}
	}
}

func main() {
	http.HandleFunc("/api/add/", addLogin)
	http.HandleFunc("/api/get/", getLogin)
	http.HandleFunc("/api/del/", delLogin)
	http.HandleFunc("/api/upd/", updLogin)
	http.ListenAndServe(":8081", nil)
}
