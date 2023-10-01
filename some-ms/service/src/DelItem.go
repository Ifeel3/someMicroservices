package src

import (
	"context"
	"encoding/json"
	"net/http"
	. "some-ms/src/structs"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DelItem(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	str := strings.Split(r.RequestURI, "/")
	if str[3] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Wrong id")
		return
	} else {
		var item ItemStruct
		err := conn.QueryRow(context.TODO(), "delete from items where id=$1 returning id, info, price, owner", str[3]).Scan(&item.Id, &item.Info, &item.Price, &item.Owner)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("Item not found")
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(item)
		}
	}
}
