package models



type Club struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Rating      float64     `json:"rating"`
	Coordinates Coordinates `json:"coordinates"`
	Clusters    []string    `json:"clusters"`
	Schedule    string      `json:"schedule"`
	Status      string      `json:"status"`
}
