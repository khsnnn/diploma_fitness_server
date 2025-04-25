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

type UniClub struct {
    ID           int                 `json:"id"`
    Name         string              `json:"name"`
    Address      string              `json:"address"`
    Description  string              `json:"description"`
    WorkingHours string              `json:"working_hours"`
    Categories   map[string][]string `json:"categories"`
    Schedule     Schedule            `json:"schedule"`
    Lat          string              `json:"lat"` // Новое поле для широты
    Lon          string              `json:"lon"` // Новое поле для долготы
    Status       string              `json:"status"`
}