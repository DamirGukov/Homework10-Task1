package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	exchangeRes := &GetMeteostatWeather{}

	r.HandleFunc("/weather", exchangeRes.GetWeather).
		Methods(http.MethodGet)

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

type GetMeteostatWeather struct {
}

const apiURL = "https://goweather.herokuapp.com/weather/New-York"

type ApiResponse struct {
	Temperature string `json:"temperature"`
	Wind        string `json:"wind"`
	Description string `json:"description"`
}

func (gmw *GetMeteostatWeather) GetWeather(w http.ResponseWriter, r *http.Request) {
	api, err := http.Get(apiURL)
	if err != nil {
		log.Println("Failed to get response data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer api.Body.Close()

	if api.StatusCode != http.StatusOK {
		log.Println("Got non-OK status from responce data:", api.Status)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response ApiResponse

	err = json.NewDecoder(api.Body).Decode(&response)
	if err != nil {
		log.Println("Failed to parse exchange API response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, response)
}

func respondWithJSON(w http.ResponseWriter, body interface{}) {
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Printf("Something went wrong while writing response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
