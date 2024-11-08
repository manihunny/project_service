package controllers

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	mock_service "project-service/internal/services/mocks"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestProjectHandler_DeleteProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectService := mock_service.NewMockProjectService(ctrl)
	projectHandler := NewProjectHandler(mockProjectService, slog.New(slog.NewTextHandler(io.Discard, nil)))

	router := gin.Default()
	router.DELETE("/projects/:id", projectHandler.DeleteProject)

	server := httptest.NewServer(router.Handler())
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	tests := []struct {
		name       string
		h          *ProjectHandler
		prepare    func(mock_service.MockProjectService)
		id         string
		statusCode int
	}{
		{
			name: "success",
			h:    projectHandler,
			prepare: func(mpsi mock_service.MockProjectService) {
				mpsi.EXPECT().DeleteProject(uint(1)).Return(nil)
			},
			id:         "1",
			statusCode: http.StatusOK,
		},
		{
			name:       "ID is not int",
			h:          projectHandler,
			prepare:    func(mpsi mock_service.MockProjectService) {},
			id:         "wrong_id",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "ID is negative int",
			h:          projectHandler,
			prepare:    func(mpsi mock_service.MockProjectService) {},
			id:         "-1",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "zero ID",
			h:          projectHandler,
			prepare:    func(mpsi mock_service.MockProjectService) {},
			id:         "0",
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt.prepare(*mockProjectService)
		t.Run(tt.name, func(t *testing.T) {
			e.DELETE("/projects/{id}").WithPath("id", tt.id).Expect().Status(tt.statusCode)
		})
	}
}
