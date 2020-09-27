package models

// Sucursal defines the properties belonging to the Sucursal object in DynamoDB.
type Sucursal struct {
	// Used for DynamoDB lookup.
	ID string
	// Address of Sucursal.
	Address string
	// Latitude of Sucursal in decimal degrees.
	Latitude float64
	// Longitude of Sucursal in decimal degrees.
	Longitude float64
}
