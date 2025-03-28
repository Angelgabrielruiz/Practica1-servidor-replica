package server

import (
	"net/http"
	"proyecto/handlers"
	"proyecto/polling"
	"github.com/rs/cors"
)

func StartServer() {
	c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Accept", "Accept-Encoding", "Authorization", "Content-Type", "X-CSRF-Token"},
    })
	handler := c.Handler(http.DefaultServeMux)

	http.HandleFunc("/addProduct", handlers.AgregarProducto)
	http.HandleFunc("/getProducts", handlers.ObtenerProductos)
	http.HandleFunc("/longpoll", polling.LongPollingHandler)
	http.HandleFunc("/countProductsInDiscount", handlers.CountProductsInDiscount)

	http.ListenAndServe(":8080", handler)
}