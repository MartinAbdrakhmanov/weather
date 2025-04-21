package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"weather/internal/model"
)

func GetCityFromIp(addr string) string {
	// addr := getClientIP(r) doesn't work on localhost
	url := fmt.Sprintf("http://ip-api.com/json/%s", addr)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	city := model.CityResponse{}
	json.Unmarshal(body, &city)
	if city.Status != "success" {
		log.Fatal("error status")
	}
	redirecturl := fmt.Sprintf("/%v", city.City)
	return redirecturl
}

func GetClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
