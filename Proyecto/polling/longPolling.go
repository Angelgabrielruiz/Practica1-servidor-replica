package polling

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
	"proyecto/models"
)

var lock sync.Mutex
var waitingClients []chan []models.Producto

func LongPollingHandler(w http.ResponseWriter, r *http.Request) {
	ch := make(chan []models.Producto)
	lock.Lock()
	waitingClients = append(waitingClients, ch)
	lock.Unlock()

	
	productosConDescuento := obtenerProductosConDescuento() 
	json.NewEncoder(w).Encode(productosConDescuento)

	select {
	case data := <-ch:
		json.NewEncoder(w).Encode(data)
	case <-time.After(30 * time.Second):
		w.WriteHeader(http.StatusNoContent)
	}
}

func NotifyClients(productos []models.Producto) {
	lock.Lock()
	defer lock.Unlock()

	for _, ch := range waitingClients {
		ch <- productos
	}
	waitingClients = nil
}


func obtenerProductosConDescuento() []models.Producto {

	resp, err := http.Get("http://localhost:8080/getProducts")
	if err == nil {
		var productos []models.Producto 
		err := json.NewDecoder(resp.Body).Decode(&productos)
		if err == nil {
			lock.Lock()
			defer lock.Unlock()

			productosConDescuento := []models.Producto{}
			for _, p := range productos {
				if p.Descuento {
					productosConDescuento = append(productosConDescuento, p)
				}
			}
			return productosConDescuento
		}
	}
	return []models.Producto{}
}