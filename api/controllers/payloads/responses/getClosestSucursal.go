package responses

import "github.com/NJRodriguez/shiny-waddle/api/models"

type ClosestSucursalResponse struct {
	Sucursal     models.Sucursal
	DistanceInKm float64
}
