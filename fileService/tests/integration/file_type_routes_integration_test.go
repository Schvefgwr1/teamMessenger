//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"fileService/internal/models"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFileTypeRoutes_Create_Get_Delete_Integration(t *testing.T) {
	env := setupIntegration(t)
	router := newTestRouter(env)

	// Используем короткое имя (максимум 20 символов для VARCHAR(20))
	// application/ = 12 символов, остаётся 8 для типа
	// Используем только последние 3 цифры timestamp для уникальности
	suffix := fmt.Sprintf("%03d", time.Now().UnixNano()%1000)
	name := fmt.Sprintf("application/t%s", suffix)

	// Create
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/file-types/?name="+name, nil)
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusOK, createW.Code)

	var created models.FileType
	require.NoError(t, json.Unmarshal(createW.Body.Bytes(), &created))
	require.NotZero(t, created.ID)
	assert.Equal(t, name, created.Name)

	// Get by ID
	getReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/file-types/%d", created.ID), nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	require.Equal(t, http.StatusOK, getW.Code)

	// Get by name (controller добавляет application/)
	pathName := created.Name[len("application/"):]
	getNameReq := httptest.NewRequest(http.MethodGet, "/api/v1/file-types/name/"+pathName, nil)
	getNameW := httptest.NewRecorder()
	router.ServeHTTP(getNameW, getNameReq)
	require.Equal(t, http.StatusOK, getNameW.Code)

	// Delete
	delReq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/file-types/%d", created.ID), nil)
	delW := httptest.NewRecorder()
	router.ServeHTTP(delW, delReq)
	require.Equal(t, http.StatusOK, delW.Code)

	// Ensure removed
	getAfterDelReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/file-types/%d", created.ID), nil)
	getAfterDelW := httptest.NewRecorder()
	router.ServeHTTP(getAfterDelW, getAfterDelReq)
	require.Equal(t, http.StatusNotFound, getAfterDelW.Code)
}

func TestFileTypeRoutes_NotFound_And_BadRequest_Integration(t *testing.T) {
	env := setupIntegration(t)
	router := newTestRouter(env)

	// Bad request: missing name
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/file-types/?name=", nil)
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusBadRequest, createW.Code)

	// Not found by id
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/file-types/9999", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	require.Equal(t, http.StatusNotFound, getW.Code)
}
