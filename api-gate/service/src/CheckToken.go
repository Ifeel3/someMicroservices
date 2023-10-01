package src

import (
	. "api-gate/src/structs"
	"io"
	"net/http"
	"strings"
)

func CheckToken(w http.ResponseWriter, r *http.Request, token TokenStruct, authAddr string) bool {
	splitted := strings.Split(r.RequestURI, "/")
	request, _ := http.NewRequest(http.MethodGet, authAddr+"/api/token/"+splitted[3], r.Body)
	response, _ := http.DefaultClient.Do(request)
	if response.StatusCode != http.StatusOK {
		w.WriteHeader(response.StatusCode)
		body, _ := io.ReadAll(response.Body)
		w.Write(body)
		return false
	}
	return true
}
