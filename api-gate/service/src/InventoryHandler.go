package src

import (
	. "api-gate/src/structs"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func InventoryHandler(w http.ResponseWriter, r *http.Request, authAddr string, someAddr string) {
	splitted := strings.Split(r.RequestURI, "/")
	if len(splitted) < 3 || splitted[3] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var itemAndToken ItemWithTokenStruct
	json.NewDecoder(r.Body).Decode(&itemAndToken)
	if !CheckToken(w, r, splitted[3], TokenStruct{Token: itemAndToken.Token}, authAddr) {
		return
	}
	item := ItemStruct{Item: itemAndToken.Item, Info: itemAndToken.Info, Price: itemAndToken.Price, Owner: itemAndToken.Owner}
	itemJson, _ := json.Marshal(item)
	request, _ := http.NewRequest(http.MethodPost, someAddr+"/api/items/"+splitted[3], bytes.NewReader(itemJson))
	response, _ := http.DefaultClient.Do(request)
	w.WriteHeader(response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	w.Write(body)
}
