package models

// Sucursal defines the properties belonging to the Sucursal object in DynamoDB.
type Sucursal struct {
	// Used for DynamoDB lookup.
	ID string `json:"id"`
	// Address of Sucursal.
	Address string `json:"address"`
	// Latitude of Sucursal in decimal degrees.
	Latitude float64 `json:"latitude"`
	// Longitude of Sucursal in decimal degrees.
	Longitude float64 `json:"longitude"`
}

type SucursalKey struct {
	ID string `json:"id"`
}
