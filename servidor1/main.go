package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Registro struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

var registros []Registro
var lastID = 0
var lock sync.Mutex
var waitingClients []chan []Registro

func agregarRegistro(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	lastID++
	nuevoRegistro := Registro{ID: lastID, Data: "Nuevo dato"}
	registros = append(registros, nuevoRegistro)


	for _, ch := range waitingClients {
		ch <- registros
	}
	waitingClients = nil

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevoRegistro)
}

func obtenerRegistros(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	json.NewEncoder(w).Encode(registros)
}

func longPolling(w http.ResponseWriter, r *http.Request) {
	ch := make(chan []Registro)
	lock.Lock()
	waitingClients = append(waitingClients, ch)
	lock.Unlock()

	select {
	case data := <-ch:
		json.NewEncoder(w).Encode(data)
	case <-time.After(30 * time.Second): 
		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {
	http.HandleFunc("/add", agregarRegistro)
	http.HandleFunc("/get", obtenerRegistros)
	http.HandleFunc("/longpoll", longPolling)

	http.ListenAndServe(":8080", nil)
}
