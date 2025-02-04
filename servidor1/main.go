package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
	"fmt"
)

type Registro struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

var registros []Registro
var lastID = 0
var lock sync.Mutex
var waitingClients []chan []Registro


func agregarRegistro(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	nombre := r.URL.Query().Get("nombre")
	apellido := r.URL.Query().Get("apellido")

	if nombre == "" || apellido == "" {
		http.Error(w, "Nombre y apellido son requeridos", http.StatusBadRequest)
		return
	}

	lastID++
	nuevoRegistro := Registro{ID: lastID, Nombre: nombre, Apellido: apellido}
	registros = append(registros, nuevoRegistro)


	fmt.Println("Registros actuales:", registros)


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


func editarRegistro(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	idStr := r.URL.Query().Get("id")
	nuevoNombre := r.URL.Query().Get("nombre")
	nuevoApellido := r.URL.Query().Get("apellido")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	encontrado := false
	for i, reg := range registros {
		if reg.ID == id {
			if nuevoNombre != "" {
				registros[i].Nombre = nuevoNombre
			}
			if nuevoApellido != "" {
				registros[i].Apellido = nuevoApellido
			}
			encontrado = true
			break
		}
	}

	if !encontrado {
		http.Error(w, "Registro no encontrado", http.StatusNotFound)
		return
	}

	for _, ch := range waitingClients {
		ch <- registros
	}
	waitingClients = nil

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registro editado"})
}


func eliminarRegistro(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	index := -1
	for i, reg := range registros {
		if reg.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Registro no encontrado", http.StatusNotFound)
		return
	}

	registros = append(registros[:index], registros[index+1:]...)

	for _, ch := range waitingClients {
		ch <- registros
	}
	waitingClients = nil

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registro eliminado"})
}

func main() {
	http.HandleFunc("/add", agregarRegistro)
	http.HandleFunc("/get", obtenerRegistros)
	http.HandleFunc("/longpoll", longPolling)
	http.HandleFunc("/edit", editarRegistro)
	http.HandleFunc("/delete", eliminarRegistro)

	http.ListenAndServe(":8080", nil)
}
