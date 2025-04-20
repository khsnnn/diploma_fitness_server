package models

type ScheduleItem struct {
	Time       string `json:"time"`
	Activity   string `json:"activity"`
	Instructor string `json:"instructor"`
}

type Schedule struct {
	Tuesday   []ScheduleItem `json:"Tuesday,omitempty"`
	Monday    []ScheduleItem `json:"Monday,omitempty"`
	Wednesday []ScheduleItem `json:"Wednesday,omitempty"`
	Thursday  []ScheduleItem `json:"Thursday,omitempty"`
	Friday    []ScheduleItem `json:"Friday,omitempty"`
	Saturday  []ScheduleItem `json:"Saturday,omitempty"`
}

// UniClub представляет университетский спортивный клуб
type UniClub struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Rating      float64     `json:"rating"`
	Coordinates Coordinates `json:"coordinates"`
	Clusters    []string    `json:"clusters"`
	Schedule    Schedule    `json:"schedule"`
	Status      string      `json:"status"`
}
