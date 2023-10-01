package src

import (
	"io"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, r *http.Request, addr string) {
	request, _ := http.NewRequest(http.MethodPatch, addr+"/api/token", r.Body)
	response, _ := http.DefaultClient.Do(request)
	w.WriteHeader(response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	w.Write(body)
}
