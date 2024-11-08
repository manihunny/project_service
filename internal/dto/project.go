package dto

import (
	"project-service/internal/models"
	"time"

	"github.com/lib/pq"
)

// ProjectDTO представляет данные проекта, которые может установить пользователь.
type ProjectDTO struct {
	Name            *string        `json:"name"`
	HourlyRate      *float64       `json:"hourlyRate"`
	StartDate       *time.Time     `json:"startDate"`
	EndDate         *time.Time     `json:"endDate,omitempty"`
	Description     *string        `json:"description,omitempty"`
	IsInternational *bool          `json:"isInternational,omitempty"`
	ContractType    *string        `json:"contractType,omitempty"`
	TechStack       pq.StringArray `json:"techStack"`
	Type            *string        `json:"type"`
	UserID          *uint          `json:"userId"`
	Logo            *string        `json:"logo,omitempty"`
	UniqueColor     *string        `json:"uniqueColor,omitempty"`
}

// Вспомогательная функция для копирования строковых значений
func copyIfNotNil[T any](dst **T, src *T) {
	if src != nil {
		*dst = src
	}
}

// Parse достает данные из модели и вставляет их в ProjectDTO
func (dto *ProjectDTO) Parse(project *models.Project) error {
	copyIfNotNil(&dto.Name, project.Name)
	copyIfNotNil(&dto.HourlyRate, project.HourlyRate)
	copyIfNotNil(&dto.StartDate, project.StartDate)
	copyIfNotNil(&dto.EndDate, project.EndDate)
	copyIfNotNil(&dto.Description, project.Description)
	copyIfNotNil(&dto.IsInternational, project.IsInternational)
	copyIfNotNil(&dto.ContractType, project.ContractType)
	dto.TechStack = project.TechStack
	copyIfNotNil(&dto.Type, project.Type)
	copyIfNotNil(&dto.UserID, project.UserID)
	copyIfNotNil(&dto.Logo, project.Logo)
	copyIfNotNil(&dto.UniqueColor, project.UniqueColor)

	return nil
}

// Map обновляет данные модели значениями из ProjectDTO
func (dto *ProjectDTO) Map(project *models.Project) error {
	copyIfNotNil(&project.Name, dto.Name)
	copyIfNotNil(&project.HourlyRate, dto.HourlyRate)
	copyIfNotNil(&project.StartDate, dto.StartDate)
	copyIfNotNil(&project.EndDate, dto.EndDate)
	copyIfNotNil(&project.Description, dto.Description)
	copyIfNotNil(&project.IsInternational, dto.IsInternational)
	copyIfNotNil(&project.ContractType, dto.ContractType)
	project.TechStack = dto.TechStack
	copyIfNotNil(&project.Type, dto.Type)
	copyIfNotNil(&project.UserID, dto.UserID)
	copyIfNotNil(&project.Logo, dto.Logo)
	copyIfNotNil(&project.UniqueColor, dto.UniqueColor)

	return nil
}
