package src

import (
	. "api-gate/src/structs"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func MarketHandler(w http.ResponseWriter, r *http.Request, authAddr string, someAddr string) {
	var token TokenStruct
	json.NewDecoder(r.Body).Decode(&token)
	if !CheckToken(w, r, token, authAddr) {
		return
	}
	splitted := strings.Split(r.RequestURI, "/")
	request, _ := http.NewRequest(http.MethodGet, someAddr+"/api/items/"+splitted[3], r.Body)
	response, _ := http.DefaultClient.Do(request)
	w.WriteHeader(response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	w.Write(body)

}
