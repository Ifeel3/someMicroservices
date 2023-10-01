package src

import (
	"context"
	"encoding/json"
	"net/http"
	. "some-ms/src/structs"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdItem(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	var item, tmp ItemStruct
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = conn.QueryRow(context.TODO(), "select id, info, price, owner from items where id=$1", item.Id).Scan(&tmp.Id, &tmp.Info, &tmp.Price, &tmp.Owner)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Item not found")
		return
	}
	err = conn.QueryRow(context.TODO(), "update items set info=$1, price=$2, owner=$3 where id=$4 returning id, info, price, owner", item.Info, item.Price, item.Owner, item.Id).Scan(&item.Id, &item.Info, &item.Price, &item.Owner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(item)
	}
}
