package requests

type PostSucursal struct {
	ID        string  `json:"id" validate:"required,uuid4"`
	Address   string  `json:"address" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
}
