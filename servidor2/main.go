
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Registro struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

func shortPolling() {
	for {
		resp, err := http.Get("http://localhost:8080/get")
		if err == nil {
			var registros []Registro
			json.NewDecoder(resp.Body).Decode(&registros)
			fmt.Println("Short Polling - Datos obtenidos:", registros)
		}
		time.Sleep(5 * time.Second) 
	}
}

func longPolling() {
	for {
		resp, err := http.Get("http://localhost:8080/longpoll")
		if err == nil && resp.StatusCode == http.StatusOK {
			var registros []Registro
			json.NewDecoder(resp.Body).Decode(&registros)
			fmt.Println("Long Polling - Datos obtenidos:", registros)
		}
	}
}

func main() {
	go shortPolling() 
	

	select {} 
}
