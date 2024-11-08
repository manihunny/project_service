package services

import (
	"errors"
	"io"
	"log/slog"
	"reflect"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"gitlab.fast-go.ru/fast-go-team/project/internal/dto"
	"gitlab.fast-go.ru/fast-go-team/project/internal/models"
	mock_repositories "gitlab.fast-go.ru/fast-go-team/project/internal/repositories/mocks"
)

func ptr[T any](obj T) *T {
	return &obj
}

func TestProjectService_CreateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectRepo := mock_repositories.NewMockProjectRepository(ctrl)
	projectService := NewProjectServiceGORM(mockProjectRepo, slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})))

	type args struct {
		projectDTO *dto.ProjectDTO
	}
	tests := []struct {
		name           string
		projectService ProjectService
		prepare        func(mpri mock_repositories.MockProjectRepository)
		args           args
		wantErr        bool
	}{
		{
			name:           "success",
			projectService: projectService,
			prepare: func(mpri mock_repositories.MockProjectRepository) {
				mpri.EXPECT().Create(&models.Project{Name: ptr("Project 1"), HourlyRate: ptr(float64(500)), StartDate: ptr(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local))}).Return(nil)
			},
			args:    args{&dto.ProjectDTO{Name: ptr("Project 1"), HourlyRate: ptr(float64(500)), StartDate: ptr(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local))}},
			wantErr: false,
		},
		{
			name:           "without Name",
			projectService: projectService,
			prepare: func(mpri mock_repositories.MockProjectRepository) {
				mpri.EXPECT().Create(&models.Project{HourlyRate: ptr(float64(500)), StartDate: ptr(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local))}).Return(errors.New("An error occurred"))
			},
			args:    args{&dto.ProjectDTO{HourlyRate: ptr(float64(500)), StartDate: ptr(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local))}},
			wantErr: true,
		},
		{
			name:           "without HourlyRate",
			projectService: projectService,
			prepare: func(mpri mock_repositories.MockProjectRepository) {
				mpri.EXPECT().Create(&models.Project{Name: ptr("Project 1"), StartDate: ptr(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local))}).Return(errors.New("An error occurred"))
			},
			args:    args{&dto.ProjectDTO{Name: ptr("Project 1"), StartDate: ptr(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local))}},
			wantErr: true,
		},
		{
			name:           "without StartDate",
			projectService: projectService,
			prepare: func(mpri mock_repositories.MockProjectRepository) {
				mpri.EXPECT().Create(&models.Project{Name: ptr("Project 1"), HourlyRate: ptr(float64(500))}).Return(errors.New("An error occurred"))
			},
			args:    args{&dto.ProjectDTO{Name: ptr("Project 1"), HourlyRate: ptr(float64(500))}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt.prepare(*mockProjectRepo)
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.projectService.CreateProject(tt.args.projectDTO); (err != nil) != tt.wantErr {
				t.Errorf("ProjectService.CreateProject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectService_DeleteProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectRepo := mock_repositories.NewMockProjectRepository(ctrl)
	projectService := NewProjectServiceGORM(mockProjectRepo, slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})))

	type args struct {
		id uint
	}
	tests := []struct {
		name           string
		projectService ProjectService
		prepare        func(mock_repositories.MockProjectRepository)
		args           args
		wantErr        bool
	}{
		{
			name:           "success",
			projectService: projectService,
			prepare: func(mpri mock_repositories.MockProjectRepository) {
				mpri.EXPECT().Delete(uint(1))
			},
			args:    args{1},
			wantErr: false,
		},
		{
			name:           "zero ID",
			projectService: projectService,
			prepare: func(mpri mock_repositories.MockProjectRepository) {
				mpri.EXPECT().Delete(uint(0)).Return(errors.New("An error occurred"))
			},
			args:    args{0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt.prepare(*mockProjectRepo)
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.projectService.DeleteProject(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("ProjectService.DeleteProject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectService_GetProjectByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectRepo := mock_repositories.NewMockProjectRepository(ctrl)
	projectService := NewProjectServiceGORM(mockProjectRepo, slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})))

	type args struct {
		id uint
	}
	tests := []struct {
		name           string
		projectService ProjectService
		prepare        func(mock_repositories.MockProjectRepository)
		args           args
		want           *models.Project
		wantErr        bool
	}{
		{
			name:           "success",
			projectService: projectService,
			prepare: func(mpri mock_repositories.MockProjectRepository) {
				mpri.EXPECT().FindByID(uint(1)).Return(&models.Project{Name: ptr("Project 1"), HourlyRate: ptr(float64(500)), StartDate: ptr(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local))}, nil)
			},
			args:    args{1},
			want:    &models.Project{Name: ptr("Project 1"), HourlyRate: ptr(float64(500)), StartDate: ptr(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local))},
			wantErr: false,
		},
		{
			name:           "zero ID",
			projectService: projectService,
			prepare: func(mpri mock_repositories.MockProjectRepository) {
				mpri.EXPECT().FindByID(uint(0)).Return(nil, errors.New("An error occurred"))
			},
			args:    args{0},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(*mockProjectRepo)
			got, err := tt.projectService.GetProjectByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectService.GetProjectByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProjectService.GetProjectByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProjectService_GetProjects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectRepo := mock_repositories.NewMockProjectRepository(ctrl)
	projectService := NewProjectServiceGORM(mockProjectRepo, slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})))

	tests := []struct {
		name           string
		projectService ProjectService
		prepare        func(mock_repositories.MockProjectRepository)
		want           []models.Project
		wantErr        bool
	}{
		{
			name:           "success",
			projectService: projectService,
			prepare: func(mpri mock_repositories.MockProjectRepository) {
				mpri.EXPECT().FindAll().Return([]models.Project{{Name: ptr("Project 1"), HourlyRate: ptr(float64(500)), StartDate: ptr(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local))}, {Name: ptr("Project 2"), HourlyRate: ptr(float64(1000)), StartDate: ptr(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local))}}, nil)
			},
			want:    []models.Project{{Name: ptr("Project 1"), HourlyRate: ptr(float64(500)), StartDate: ptr(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local))}, {Name: ptr("Project 2"), HourlyRate: ptr(float64(1000)), StartDate: ptr(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local))}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.prepare(*mockProjectRepo)
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.projectService.GetProjects()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectService.GetProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProjectService.GetProjects() = %v, want %v", got, tt.want)
			}
		})
	}
}
