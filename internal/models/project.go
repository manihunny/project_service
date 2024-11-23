package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Project struct {
	gorm.Model
	Name            *string        `json:"name" gorm:"type:varchar(255);not null"` // Название проекта
	HourlyRate      *float64       `json:"hourlyRate" gorm:"not null"`             // Рейт проекта в виде цены за час
	StartDate       *time.Time     `json:"startDate" gorm:"not null"`              // Дата начала проекта
	EndDate         *time.Time     `json:"endDate"`                                // Дата окончания проекта или не определено
	Description     *string        `json:"description" gorm:"type:text"`           // Описание проекта
	IsInternational *bool          `json:"isInternational" gorm:"default:false"`   // Зарубежный проект или российский проект
	ContractType    *string        `json:"contractType" gorm:"type:varchar(10)"`   // Взят на ООО или на ИП
	TechStack       pq.StringArray `json:"techStack" gorm:"type:text[]"`           // Технологический стек проекта
	Type            *string        `json:"type" gorm:"type:varchar(20)"`           // Бэкенд, фронтенд или фуллстек
	UserID          *uint          `json:"userId"`                                 //uint пользователя, который ведет проект
	Logo            *string        `json:"logo" gorm:"type:varchar(255)"`          // Логотип проекта (URL или путь к изображению)
	UniqueColor     *string        `json:"uniqueColor" gorm:"type:varchar(20)"`    // Уникальный цвет проекта (в формате HEX, RGB и т.д.)
}

func (Project) TableName() string {
	return "projects"
}
