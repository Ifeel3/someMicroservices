package src

import (
	"context"
	"encoding/json"
	"net/http"
	. "some-ms/src/structs"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddItem(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	var item ItemStruct
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Bad request"}})
		return
	}
	err = conn.QueryRow(context.TODO(), "insert into items (info, price, owner) values ($1, $2, $3) returning id", item.Info, item.Price, item.Owner).Scan(&item.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Bad request"}})
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "OK"}, item})
	}
}
