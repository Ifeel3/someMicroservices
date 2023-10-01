package src

import (
	"io"
	"net/http"
)

func RegistrationHandler(w http.ResponseWriter, r *http.Request, addr string) {
	request, _ := http.NewRequest(http.MethodPost, addr+"/api/login", r.Body)
	response, _ := http.DefaultClient.Do(request)
	if response.StatusCode == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login created"))
		return
	}
	w.WriteHeader(response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	w.Write(body)
}
