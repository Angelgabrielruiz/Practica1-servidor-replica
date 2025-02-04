
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Registro struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}


func shortPolling() {
	for {
		resp, err := http.Get("http://localhost:8080/get")
		if err == nil {
			var registros []Registro
			err := json.NewDecoder(resp.Body).Decode(&registros)
			if err != nil {
				fmt.Println("Error decodificando JSON:", err)
			} else {
				fmt.Println("Short Polling - Datos obtenidos:", registros)
			}
		} else {
			fmt.Println("Error en la solicitud HTTP:", err)
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
