package models

type Producto struct {
	ID        int    `json:"id"`
	Nombre    string `json:"nombre"`
	Precio    int    `json:"precio"`
	Codigo    string `json:"codigo"`
	Descuento bool   `json:"descuento"`
}
