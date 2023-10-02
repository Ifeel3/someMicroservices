package src

import (
	"bytes"
	"io"
	"net/http"
)

func RegistrationHandler(w http.ResponseWriter, r *http.Request, addr string) {
	body, _ := io.ReadAll(r.Body)
	request, _ := http.NewRequest(http.MethodPost, addr+"/api/login", bytes.NewReader(body))
	response, _ := http.DefaultClient.Do(request)
	if response.StatusCode == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login created"))
		return
	}
	w.WriteHeader(response.StatusCode)
	body, _ = io.ReadAll(response.Body)
	w.Write(body)
}
