package main

import (
	"encoding/json"            //Provides functions to encode and decode JSON.
	"io/ioutil"                 //Contains functions for reading and writing files
	"net/http"					//Provides HTTP client and server implementations
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`				// API key for OpenWeatherMap
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
		Celsius float64 `json:"temp_celsius"` 
	} `json:"main"`
}

func LoadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData
	err = json.Unmarshal(bytes, &c)         //Decodes the JSON data into the c variable of type apiConfigData.
	if err != nil {
		return apiConfigData{}, err
	}

	return c, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func query(city string) (weatherData, error) {
	apiConfig, err := LoadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?appid=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()                     //Ensures the response body is closed after reading.

	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil { 					//Decodes the JSON response into the d variable of type weatherData.
		return weatherData{}, err 
	}

	d.Main.Celsius = d.Main.Kelvin - 273.15
	
	return d, nil
}

func main() {
	http.HandleFunc("/hello", hello)

	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}
