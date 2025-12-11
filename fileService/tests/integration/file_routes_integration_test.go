//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fileService/internal/dto"
	"fileService/internal/models"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
	"time"
)

func TestFileRoutes_Upload_Get_Rename_List_Integration(t *testing.T) {
	env := setupIntegration(t)
	router := newTestRouter(env)

	// Upload
	body, contentType := mustMultipartBody(t, "upload.png", "image/png", []byte("route-upload"))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/files/upload", body)
	req.Header.Set("Content-Type", contentType)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var uploaded models.File
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &uploaded))
	require.NotZero(t, uploaded.ID)

	// Get
	getReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/files/%d", uploaded.ID), nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	require.Equal(t, http.StatusOK, getW.Code)

	var fetched models.File
	require.NoError(t, json.Unmarshal(getW.Body.Bytes(), &fetched))
	assert.Equal(t, uploaded.ID, fetched.ID)
	assert.Contains(t, fetched.URL, env.Config.MinIO.Bucket)

	// Rename
	newName := fmt.Sprintf("renamed-%d", time.Now().UnixNano())
	renameReq := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/files/%d?name=%s", uploaded.ID, newName), nil)
	renameW := httptest.NewRecorder()
	router.ServeHTTP(renameW, renameReq)
	require.Equal(t, http.StatusOK, renameW.Code)

	var renamed models.File
	require.NoError(t, json.Unmarshal(renameW.Body.Bytes(), &renamed))
	assert.Contains(t, renamed.Name, newName)
	assert.Contains(t, renamed.URL, renamed.Name)

	// List names with pagination
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/files/names?limit=1&offset=0", nil)
	listW := httptest.NewRecorder()
	router.ServeHTTP(listW, listReq)
	require.Equal(t, http.StatusOK, listW.Code)

	var names []dto.FileInformation
	require.NoError(t, json.Unmarshal(listW.Body.Bytes(), &names))
	require.Len(t, names, 1)
}

func TestFileRoutes_Upload_UnsupportedType_Integration(t *testing.T) {
	env := setupIntegration(t)
	router := newTestRouter(env)

	body, contentType := mustMultipartBody(t, "bad.bin", "application/unknown", []byte("bad type"))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/files/upload", body)
	req.Header.Set("Content-Type", contentType)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnsupportedMediaType, w.Code)
	assert.Contains(t, w.Body.String(), "Unsupported file type")
}

func TestFileRoutes_Get_NotFound_Integration(t *testing.T) {
	env := setupIntegration(t)
	router := newTestRouter(env)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/files/9999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestFileRoutes_Rename_BadRequest_Integration(t *testing.T) {
	env := setupIntegration(t)
	router := newTestRouter(env)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/files/abc?name=", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFileRoutes_GetNames_InvalidPagination_Integration(t *testing.T) {
	env := setupIntegration(t)
	router := newTestRouter(env)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/files/names?limit=0&offset=-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

// mustMultipartBody формирует multipart тело с явным Content-Type части.
func mustMultipartBody(t *testing.T, filename, contentType string, content []byte) (*bytes.Buffer, string) {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	return &body, writer.FormDataContentType()
}
