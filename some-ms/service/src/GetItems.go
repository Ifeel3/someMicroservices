package src

import (
	"context"
	"encoding/json"
	"net/http"
	. "some-ms/src/structs"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetItems(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	str := strings.Split(r.RequestURI, "/")
	if len(str) == 4 && str[3] != "" {
		data, err := conn.Query(context.TODO(), "select id, info, price, owner from items where owner=$1", str[3])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Not found"}})
			return
		}
		defer data.Close()
		var result []ItemStruct
		for data.Next() {
			var tmp ItemStruct
			err := data.Scan(&tmp.Id, &tmp.Info, &tmp.Price, &tmp.Owner)
			if err != nil {
				break
			}
			result = append(result, tmp)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "OK"}, result})
		return
	} else if len(str) == 5 && str[4] != "" {
		var result ItemStruct
		err := conn.QueryRow(context.TODO(), "select id, info, price, owner from items where owner=$1 and id=$2", str[3], str[4]).Scan(&result.Id, &result.Info, &result.Price, &result.Owner)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Not found"}})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "OK"}, result})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Bad request"}})
		return
	}
}
