package main

import (
	. "api-gate/src"
	"net/http"
)

var someAddr string = "http://some-ms:8082"
var authAddr string = "http://auth-ms:8081"

func AuthHandlerWrap(w http.ResponseWriter, r *http.Request) {
	AuthHandler(w, r, authAddr)
}

func RegistrationHandlerWrap(w http.ResponseWriter, r *http.Request) {
	RegistrationHandler(w, r, authAddr)
}

func MarketHandlerWrap(w http.ResponseWriter, r *http.Request) {
	MarketHandler(w, r, authAddr, someAddr)
}

func InventoryHandlerWrap(w http.ResponseWriter, r *http.Request) {
	InventoryHandler(w, r, authAddr, someAddr)
}

func main() {
	http.HandleFunc("/api/auth/login", AuthHandlerWrap)
	http.HandleFunc("/api/auth/reg", RegistrationHandlerWrap)
	http.HandleFunc("/api/market/", MarketHandlerWrap)
	http.HandleFunc("/api/market", MarketHandlerWrap)
	http.HandleFunc("/api/inv", InventoryHandlerWrap)
	http.HandleFunc("/api/inv/", InventoryHandlerWrap)
	http.ListenAndServe(":8080", nil)
}
