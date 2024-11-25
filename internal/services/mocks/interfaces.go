// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "gitlab.fast-go.ru/fast-go-team/project/internal/dto"
	models "gitlab.fast-go.ru/fast-go-team/project/internal/models"
	repositories "gitlab.fast-go.ru/fast-go-team/project/internal/repositories"
)

// MockProjectService is a mock of ProjectService interface.
type MockProjectService struct {
	ctrl     *gomock.Controller
	recorder *MockProjectServiceMockRecorder
}

// MockProjectServiceMockRecorder is the mock recorder for MockProjectService.
type MockProjectServiceMockRecorder struct {
	mock *MockProjectService
}

// NewMockProjectService creates a new mock instance.
func NewMockProjectService(ctrl *gomock.Controller) *MockProjectService {
	mock := &MockProjectService{ctrl: ctrl}
	mock.recorder = &MockProjectServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectService) EXPECT() *MockProjectServiceMockRecorder {
	return m.recorder
}

// CreateProject mocks base method.
func (m *MockProjectService) CreateProject(projectDTO *dto.ProjectDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProject", projectDTO)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProject indicates an expected call of CreateProject.
func (mr *MockProjectServiceMockRecorder) CreateProject(projectDTO interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProject", reflect.TypeOf((*MockProjectService)(nil).CreateProject), projectDTO)
}

// DeleteProject mocks base method.
func (m *MockProjectService) DeleteProject(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProject", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProject indicates an expected call of DeleteProject.
func (mr *MockProjectServiceMockRecorder) DeleteProject(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProject", reflect.TypeOf((*MockProjectService)(nil).DeleteProject), id)
}

// GetProjectByID mocks base method.
func (m *MockProjectService) GetProjectByID(id uint) (*models.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectByID", id)
	ret0, _ := ret[0].(*models.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectByID indicates an expected call of GetProjectByID.
func (mr *MockProjectServiceMockRecorder) GetProjectByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectByID", reflect.TypeOf((*MockProjectService)(nil).GetProjectByID), id)
}

// GetProjects mocks base method.
func (m *MockProjectService) GetProjects() ([]models.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjects")
	ret0, _ := ret[0].([]models.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjects indicates an expected call of GetProjects.
func (mr *MockProjectServiceMockRecorder) GetProjects() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjects", reflect.TypeOf((*MockProjectService)(nil).GetProjects))
}

// GetRepo mocks base method.
func (m *MockProjectService) GetRepo() repositories.ProjectRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepo")
	ret0, _ := ret[0].(repositories.ProjectRepository)
	return ret0
}

// GetRepo indicates an expected call of GetRepo.
func (mr *MockProjectServiceMockRecorder) GetRepo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepo", reflect.TypeOf((*MockProjectService)(nil).GetRepo))
}

// UpdateProject mocks base method.
func (m *MockProjectService) UpdateProject(id uint, projectDTO *dto.ProjectDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProject", id, projectDTO)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProject indicates an expected call of UpdateProject.
func (mr *MockProjectServiceMockRecorder) UpdateProject(id, projectDTO interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProject", reflect.TypeOf((*MockProjectService)(nil).UpdateProject), id, projectDTO)
}
