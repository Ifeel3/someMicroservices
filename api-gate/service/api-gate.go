package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LoginStruct struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

type TokenStruct struct {
	Token string `json:"token"`
	Data  int64  `json:"data"`
}

type StatusStruct struct {
	Status     string `json:"status"`
	StatusInfo string `json:"statusinfo"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	var login LoginStruct
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Bad request"}})
		return
	}
	request, _ := http.NewRequest(http.MethodGet, authAddr+"/api/login", r.Body)
	response, _ := http.DefaultClient.Do(request)
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "Bad request"}})
		return
	}
	status := res[0].(StatusStruct)
	if status.Status == "Error" && status.StatusInfo == "Not found" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]any{StatusStruct{Status: "Error", StatusInfo: "User not found"}})
		return
	}
}

// var someAddr string = "http://some-ms:8082"
var authAddr string = "http://localhost:8081"

func main() {
	http.HandleFunc("/api/auth/login", AuthHandler)
	err := http.ListenAndServe(":8079", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
