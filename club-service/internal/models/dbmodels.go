package models


// Модели для работы с базой данных (GORM)
type DBClub struct {
    ID           int            `json:"id" gorm:"primaryKey"`
    Name         string         `json:"name"`
    Address      string         `json:"address"`
    Description  string         `json:"description"`
    WorkingHours string         `json:"working_hours"`
    Rating       float64        `json:"rating"`
    Lat          float64        `json:"lat"`
    Lon          float64        `json:"lon"`
    Type         string         `json:"type"`
    Status       string         `json:"status"`
    Categories   []DBCategory   `json:"categories" gorm:"many2many:club_categories;foreignKey:ID;joinForeignKey:ClubID;References:ID;joinReferences:CategoryID"`
    Schedules    []DBSchedule   `json:"schedule" gorm:"foreignKey:ClubID"`
}

func (DBClub) TableName() string {
    return "clubs"
}

type DBCategory struct {
    ID           int            `json:"-" gorm:"primaryKey"`
    Name         string         `json:"name"`
    Subcategories []DBSubcategory `json:"subcategories" gorm:"foreignKey:CategoryID"`
}

func (DBCategory) TableName() string {
    return "categories"
}

type DBSubcategory struct {
    ID         int    `json:"-" gorm:"primaryKey"`
    CategoryID int    `json:"-"`
    Name       string `json:"name"`
}

func (DBSubcategory) TableName() string {
    return "subcategories"
}

type DBSchedule struct {
    ID         int    `json:"-" gorm:"primaryKey"`
    ClubID     int    `json:"-"`
    Day        string `json:"day"`
    Time       string `json:"time"`
    Activity   string `json:"activity"`
    Instructor string `json:"instructor"`
}

func (DBSchedule) TableName() string {
    return "schedules"
}

// Структура для ответа API
type APICategory struct {
    Name         string   `json:"name"`
    Subcategories []string `json:"subcategories"`
}

type APIClub struct {
    ID           int           `json:"id"`
    Name         string        `json:"name"`
    Address      string        `json:"address"`
    Description  string        `json:"description"`
    WorkingHours string        `json:"working_hours"`
    Rating       float64       `json:"rating"`
    Lat          float64       `json:"lat"`
    Lon          float64       `json:"lon"`
    Type         string        `json:"type"`
    Status       string        `json:"status"`
    Categories   []APICategory `json:"categories"`
    Schedules    []DBSchedule  `json:"schedule"`
}