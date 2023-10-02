package src

import (
	. "api-gate/src/structs"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func CheckToken(w http.ResponseWriter, r *http.Request, user string, token TokenStruct, authAddr string) bool {
	jsonBody, _ := json.Marshal(token)
	request, _ := http.NewRequest(http.MethodGet, authAddr+"/api/token/"+user, bytes.NewReader(jsonBody))
	response, _ := http.DefaultClient.Do(request)
	if response.StatusCode != http.StatusOK {
		w.WriteHeader(response.StatusCode)
		body, _ := io.ReadAll(response.Body)
		w.Write(body)
		return false
	}
	return true
}
