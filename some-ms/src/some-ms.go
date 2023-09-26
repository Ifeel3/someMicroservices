package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Item struct {
	Status string
	Id     int32  `json:"id"`
	Info   string `json:"info"`
	Price  int64  `json:"price"`
	Owner  string `json:"owner"`
}

var addr string = "postgres://postgres:pass1234@db:5432/postgres"

func addItem(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), addr)
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
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(item)
	}
}

func getItems(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), addr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Item{"Error", 0, "Database not connected", 0, ""})
		return
	}
	defer conn.Close(context.TODO())
	str := strings.Split(r.RequestURI, "/")
	if len(str) == 4 && str[3] != "" {
		data, err := conn.Query(context.TODO(), "select id, info, price, owner from items where owner=$1", str[3])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode([]Item{{"Error", 0, "Not found", 0, ""}})
			return
		}
		defer data.Close()
		var result []Item
		for data.Next() {
			var tmp Item
			tmp.Status = "OK"
			err := data.Scan(&tmp.Id, &tmp.Info, &tmp.Price, &tmp.Owner)
			if err != nil {
				break
			}
			result = append(result, tmp)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return
	} else if len(str) == 5 && str[4] != "" {
		var result []Item = make([]Item, 1)
		err := conn.QueryRow(context.TODO(), "select id, info, price, owner from items where owner=$1 and id=$2", str[3], str[4]).Scan(&result[0].Id, &result[0].Info, &result[0].Price, &result[0].Owner)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode([]Item{{"Error", 0, "Not found", 0, ""}})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

func delItem(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), addr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Item{"Error", 0, "Database not connected", 0, ""})
		return
	}
	defer conn.Close(context.TODO())
	str := strings.Split(r.RequestURI, "/")
	if str[3] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Item{"Error", 0, "Wrong id", 0, ""})
		return
	} else {
		var item Item
		err := conn.QueryRow(context.TODO(), "delete from items where id=$1 returning id, info, price, owner", str[3]).Scan(&item.Id, &item.Info, &item.Price, &item.Owner)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Item{"Error", 0, "Item not found", 0, ""})
			return
		} else {
			item.Status = "OK"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(item)
		}
	}
}

func updItem(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.TODO(), addr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Item{"Error", 0, "Database not connected", 0, ""})
		return
	}
	defer conn.Close(context.TODO())
	var item, tmp Item
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Item{"Error", 0, "Bad request", 0, ""})
		return
	}
	err = conn.QueryRow(context.TODO(), "select id, info, price, owner from items where id=$1", item.Id).Scan(&tmp.Id, &tmp.Info, &tmp.Price, &tmp.Owner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Item{"Error", 0, "Item not found", 0, ""})
		return
	}
	err = conn.QueryRow(context.TODO(), "update items set info=$1, price=$2, owner=$3 where id=$4 returning id, info, price, owner", item.Info, item.Price, item.Owner, item.Id).Scan(&item.Id, &item.Info, &item.Price, &item.Owner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Item{"Error", 0, "Bad request", 0, ""})
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(item)
	}
}

func main() {
	http.HandleFunc("/api/add", addItem)
	http.HandleFunc("/api/get/", getItems)
	http.HandleFunc("/api/del/", delItem)
	http.HandleFunc("/api/upd", updItem)
	http.ListenAndServe(":8082", nil)
}
