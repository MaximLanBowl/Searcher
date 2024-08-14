package models

type Driver struct {
	DriverID  int     `json:"driver_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RouteResponse struct {
	Routes []struct {
		Duration float64 `json:"duration"`
	} `json:"routes"`
}
