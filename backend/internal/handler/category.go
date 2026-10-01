package handler

import (
	"defect-management-system/backend/internal/repository"
	"defect-management-system/backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	repo *repository.CategoryRepository
}

func NewCategoryHandler(repo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.repo.GetAll()
	if err != nil {
		response.Internal(c, "Failed to fetch categories", err)
		return
	}
	c.JSON(http.StatusOK, categories)
}
