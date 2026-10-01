package handler

import (
	"defect-management-system/backend/internal/model"
	"defect-management-system/backend/internal/repository"
	"defect-management-system/backend/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BuildingHandler struct {
	repo *repository.BuildingRepository
}

func NewBuildingHandler(repo *repository.BuildingRepository) *BuildingHandler {
	return &BuildingHandler{repo: repo}
}

func (h *BuildingHandler) GetBuildings(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	if pageStr == "" && limitStr == "" {
		buildings, err := h.repo.GetAll()
		if err != nil {
			response.Internal(c, "Failed to fetch buildings", err)
			return
		}
		c.JSON(http.StatusOK, buildings)
		return
	}

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	sortBy := c.Query("sortBy")
	search := c.Query("search")

	res, err := h.repo.GetPaginated(page, limit, search, sortBy)
	if err != nil {
		response.Internal(c, "Failed to fetch paginated buildings", err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *BuildingHandler) GetBuilding(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid building ID", err)
		return
	}

	building, err := h.repo.GetByID(id)
	if err != nil {
		response.Internal(c, "Database error", err)
		return
	}
	if building == nil {
		response.NotFound(c, "Building not found")
		return
	}

	c.JSON(http.StatusOK, building)
}

func (h *BuildingHandler) CreateBuilding(c *gin.Context) {
	var b model.Building
	if err := c.ShouldBindJSON(&b); err != nil {
		response.BadRequest(c, "Invalid building data", err)
		return
	}

	if b.Name == "" || b.Address == "" {
		response.BadRequest(c, "Name and address are required", nil)
		return
	}

	created, err := h.repo.Create(&b)
	if err != nil {
		response.Internal(c, "Failed to create building", err)
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *BuildingHandler) UpdateBuilding(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid building ID", err)
		return
	}

	var b model.Building
	if err := c.ShouldBindJSON(&b); err != nil {
		response.BadRequest(c, "Invalid building data", err)
		return
	}
	b.ID = id

	updated, err := h.repo.Update(&b)
	if err != nil {
		response.Internal(c, "Failed to update building", err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *BuildingHandler) DeleteBuilding(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid building ID", err)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		response.Internal(c, "Failed to delete building", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
