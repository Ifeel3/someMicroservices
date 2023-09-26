package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type Item struct {
	Status string
	Id     int32  `json:"id"`
	Info   string `json:"info"`
	Price  int64  `json:"price"`
	Owner  string `json:"owner"`
}

func addItem(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), "postgres://postgres:pass1234@localhost:5432/postgres")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Item{"Error", 0, "Database not connected", 0, ""})
		return
	}
	defer conn.Close(context.TODO())
	var item Item
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Item{"Error", 0, "Bad request", 0, ""})
		return
	}
	err = conn.QueryRow(context.TODO(), "insert into items (info, price, owner) values ($1, $2, $3) returning id", item.Info, item.Price, item.Owner).Scan(&item.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Item{"Error", 0, "Bad request", 0, ""})
		return
	} else {
		item.Status = "OK"
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(item)
	}
}

/*
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
		err := conn.QueryRow(context.TODO(), "select login, hash from auth where login=$1", str[3]).Scan(&ans.Login, &ans.Hash)
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
		err := conn.QueryRow(context.TODO(), "select login, hash from auth where login=$1", str[3]).Scan(&ans.Login, &ans.Hash)
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
}*/

func main() {
	http.HandleFunc("/api/add", addItem)
	/*	http.HandleFunc("/api/get", getLogin)
		http.HandleFunc("/api/del", delLogin)
		http.HandleFunc("/api/upd", updLogin)*/
	http.ListenAndServe(":8082", nil)
}
