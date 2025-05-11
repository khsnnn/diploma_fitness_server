package models

type Club struct {
	ID           int                 `json:"id"`
	Name         string              `json:"name"`
	Address      string              `json:"address"`
	Description  string              `json:"description"`
	WorkingHours string              `json:"working_hours"`
	Rating       string              `json:"rating"` 
	Categories   map[string][]string `json:"categories"`
	Lat          string              `json:"lat"` 
	Lon          string              `json:"lon"` 
	Status       string              `json:"status"`
}
