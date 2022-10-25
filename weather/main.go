package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiconfigdata struct {
	OpenWeatherMapApi string `json:"openWeatherMapApi"`
}

type weatherdata struct {
	name string `json:"name"`
	main struct {
		kelvin float64 `json:"temp"`
	} `json:"main"`
}

func LoadApiConfig(filename string) (apiconfigdata, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiconfigdata{}, err
	}
	var c apiconfigdata
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiconfigdata{}, err
	}
	return c, nil
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
		w.Header().Set("contenet-type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})
	http.ListenAndServe(":9090", nil)

}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go\n"))

}
func query(city string) (weatherdata, error) {
	apiconfig, err := LoadApiConfig(".apiconfig")
	if err != nil {
		return weatherdata{}, err
	}
	resp, err := http.Get("api.openweathermap.org/data/2.5/weather?APPID=}" + apiconfig.OpenWeatherMapApi + "&q=" + city)
	if err != nil {
		return weatherdata{}, err
	}
	defer resp.Body.Close()
	var d weatherdata
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherdata{}, err
	}
	return d, nil
}
