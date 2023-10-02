package src

import (
	"bytes"
	"io"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, r *http.Request, addr string) {
	body, _ := io.ReadAll(r.Body)
	request, _ := http.NewRequest(http.MethodPatch, addr+"/api/token", bytes.NewReader(body))
	response, _ := http.DefaultClient.Do(request)
	w.WriteHeader(response.StatusCode)
	body, _ = io.ReadAll(response.Body)
	w.Write(body)
}
