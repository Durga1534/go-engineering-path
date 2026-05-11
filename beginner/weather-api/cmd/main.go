package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"weather.api/internal/cache"
	"weather.api/internal/weather"
)

var myCache *cache.Cache

func main() {
	myCache = cache.NewCache(os.Getenv("REDIS_ADDR"))

	http.HandleFunc("/weather", handleWeather)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Weather API started in %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "Query parameter 'city' is required", http.StatusBadRequest)
		return
	}

	cachedData, err := myCache.Get(city)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "HIT")
		fmt.Fprint(w, cachedData)
		return
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Println("ERROR: WEATHER_API_KEY is not set in environment variables")
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}
	log.Printf("Cache Miss for '%s'. Fetching from API...", city)
	freshData, err := weather.FetchWeather(city, apiKey)
	if err != nil {
		log.Printf("API Error: %v", err)
		http.Error(w, "Could not fetch weather data", http.StatusServiceUnavailable)
		return
	}

	err = myCache.Set(city, freshData, 12*time.Hour)
	if err != nil {
		log.Printf("Cache set Error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	fmt.Fprint(w, freshData)
}
