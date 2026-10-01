package handler

import (
	"defect-management-system/backend/internal/model"
	"defect-management-system/backend/internal/repository"
	"defect-management-system/backend/pkg/response"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	repo      *repository.UploadRepository
	uploadDir string
	baseURL   string
}

func NewUploadHandler(repo *repository.UploadRepository, uploadDir, baseURL string) *UploadHandler {
	_ = os.MkdirAll(uploadDir, 0755)
	return &UploadHandler{
		repo:      repo,
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "File is required", err)
		return
	}

	ext := filepath.Ext(file.Filename)
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dst := filepath.Join(h.uploadDir, uniqueName)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.Internal(c, "Failed to save file", err)
		return
	}

	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	fileURL := fmt.Sprintf("%s/api/uploads/files/%s", h.baseURL, uniqueName)

	upload := model.UploadFile{
		OriginalName: file.Filename,
		StoredName:   uniqueName,
		FilePath:     dst,
		URL:          fileURL,
		Bytes:        file.Size,
		Format:       mimeType,
		MimeType:     mimeType,
	}

	saved, err := h.repo.Save(&upload)
	if err != nil {
		response.Internal(c, "Failed to record upload in database", err)
		return
	}

	c.JSON(http.StatusOK, saved)
}

func (h *UploadHandler) GetFiles(c *gin.Context) {
	idParams := c.QueryArray("id")
	if len(idParams) == 0 {
		singleID := c.Query("id")
		if singleID != "" {
			idParams = strings.Split(singleID, ",")
		}
	}

	var ids []int
	for _, idStr := range idParams {
		if id, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
			ids = append(ids, id)
		}
	}

	files, err := h.repo.GetByIDs(ids)
	if err != nil {
		response.Internal(c, "Failed to fetch files", err)
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *UploadHandler) ServeFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(h.uploadDir, filepath.Clean(filename))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.NotFound(c, "File not found")
		return
	}

	c.File(filePath)
}
