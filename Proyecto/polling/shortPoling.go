package polling

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"proyecto/models"
)

func ShortPolling() {
	for {
		resp, err := http.Get("http://localhost:8080/getProducts")
		if err == nil {
			var productos []models.Producto 
			err := json.NewDecoder(resp.Body).Decode(&productos)
			if err != nil {
				fmt.Println("Error decodificando JSON:", err)
			} else {
				fmt.Println("Short Polling - Datos obtenidos:", productos)
			}
		} else {
			fmt.Println("Error en la solicitud HTTP:", err)
		}
		time.Sleep(5 * time.Second)
	}
}
