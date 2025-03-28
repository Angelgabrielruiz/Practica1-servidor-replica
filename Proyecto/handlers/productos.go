package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"proyecto/models"
	"proyecto/polling"
)

var productos []models.Producto
var lastID = 0
var lock sync.Mutex

func AgregarProducto(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	nombre := r.URL.Query().Get("nombre")
	precioStr := r.URL.Query().Get("precio")
	codigo := r.URL.Query().Get("codigo")
	descuentoStr := r.URL.Query().Get("descuento")

	if nombre == "" || precioStr == "" || codigo == "" {
		http.Error(w, "Nombre, precio y código son requeridos", http.StatusBadRequest)
		return
	}

	precio, err := strconv.Atoi(precioStr)
	if err != nil {
		http.Error(w, "Precio debe ser un número entero", http.StatusBadRequest)
		return
	}

	descuento, err := strconv.ParseBool(descuentoStr)
	if err != nil {
		descuento = false
	}

	lastID++
	nuevoProducto := models.Producto{ID: lastID, Nombre: nombre, Precio: precio, Codigo: codigo, Descuento: descuento}
	productos = append(productos, nuevoProducto)

	fmt.Println("Productos actuales:", productos)

	polling.NotifyClients(productos)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevoProducto)
}

func ObtenerProductos(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()
	json.NewEncoder(w).Encode(productos)
}


func CountProductsInDiscount(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	count := 0
	for _, producto := range productos {
		if producto.Descuento {
			count++
		}
	}

	json.NewEncoder(w).Encode(count)
}
