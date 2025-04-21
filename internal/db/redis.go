package db

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"weather/internal/client"
	"weather/internal/service"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_ADDR"),
	Password: os.Getenv("REDIS_PASSWORD"),
	DB:       0,
	Protocol: 2,
})

func GetSuggestion(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		log.Fatal("city is required", http.StatusBadRequest)
		return
	}
	log.Printf("Preparing suggestion for %v \n", city)
	suggestion, err := rdb.Get(ctx, city).Result()
	if err == redis.Nil {
		forecast := client.GetForecast(city, 2)
		service.FilterForecastByTime(&forecast)
		suggestion = client.GetClothingSuggestion(forecast)
		err = rdb.Set(ctx, city, suggestion, 1*time.Hour).Err() // change me
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"suggestion": suggestion,
	})
}
