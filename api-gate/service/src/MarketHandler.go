package src

import (
	. "api-gate/src/structs"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func MarketHandler(w http.ResponseWriter, r *http.Request, authAddr string, someAddr string) {
	splitted := strings.Split(r.RequestURI, "/")
	if len(splitted) < 4 || splitted[3] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var token TokenStruct
	json.NewDecoder(r.Body).Decode(&token)
	if !CheckToken(w, r, splitted[3], token, authAddr) {
		return
	}
	request, _ := http.NewRequest(http.MethodGet, someAddr+"/api/items/"+splitted[3], nil)
	response, _ := http.DefaultClient.Do(request)
	w.WriteHeader(response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	w.Write(body)

}
