package api

import (
	"log"
	"net/http"
	"weather/internal/client"
)

func GetCity(w http.ResponseWriter, r *http.Request) {
	// addr := client.GetClientIP(r) doesn't work on localhost
	addr := "188.126.76.169"
	log.Print(addr)
	redirecturl := client.GetCityFromIp(addr)

	http.Redirect(w, r, redirecturl, http.StatusFound)
}
