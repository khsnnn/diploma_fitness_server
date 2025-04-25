package models

type Club struct {
	ID           int                 `json:"id"`
	Name         string              `json:"name"`
	Address      string              `json:"address"`
	Description  string              `json:"description"`
	WorkingHours string              `json:"working_hours"`
	Rating       string              `json:"rating"` // String in JSON, will be normalized to float
	Categories   map[string][]string `json:"categories"`
	Lat          string              `json:"lat"` // Новое поле для широты
	Lon          string              `json:"lon"` // Новое поле для долготы
	Status       string              `json:"status"`
}
