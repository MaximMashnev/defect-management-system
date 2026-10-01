package handler

import (
	"defect-management-system/backend/internal/model"
	"defect-management-system/backend/internal/repository"
	"defect-management-system/backend/internal/service"
	"defect-management-system/backend/pkg/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	appRepo         *repository.ApplicationRepository
	buildingRepo    *repository.BuildingRepository
	priorityService *service.PriorityService
}

func NewApplicationHandler(
	appRepo *repository.ApplicationRepository,
	buildingRepo *repository.BuildingRepository,
	priorityService *service.PriorityService,
) *ApplicationHandler {
	return &ApplicationHandler{
		appRepo:         appRepo,
		buildingRepo:    buildingRepo,
		priorityService: priorityService,
	}
}

func (h *ApplicationHandler) GetApplications(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	if pageStr == "" && limitStr == "" {
		apps, err := h.appRepo.GetAll()
		if err != nil {
			response.Internal(c, "Failed to fetch applications", err)
			return
		}
		c.JSON(http.StatusOK, apps)
		return
	}

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	sortBy := c.Query("sortBy")

	filters := make(map[string]string)
	queryParams := []string{"status", "building_id", "category", "category_id", "email", "priority", "name"}
	for _, qp := range queryParams {
		if val := c.Query(qp); val != "" {
			filters[qp] = val
		}
	}

	paginated, err := h.appRepo.GetPaginated(page, limit, filters, sortBy)
	if err != nil {
		response.Internal(c, "Failed to fetch paginated applications", err)
		return
	}

	c.JSON(http.StatusOK, paginated)
}

func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid application ID", err)
		return
	}

	app, err := h.appRepo.GetByID(id)
	if err != nil {
		response.Internal(c, "Database error", err)
		return
	}
	if app == nil {
		response.NotFound(c, "Application not found")
		return
	}

	c.JSON(http.StatusOK, app)
}

type CreateAppRequest struct {
	Name                  string `json:"name" binding:"required"`
	Description           string `json:"description" binding:"required"`
	Email                 string `json:"email" binding:"required"`
	BuildingID            int    `json:"building_id" binding:"required"`
	Category              *int   `json:"category"`
	UploadID              []int  `json:"upload_id"`
	LevelCriticalityScore int    `json:"levelCriticalityScore"`
	Priority              *int   `json:"priority"`
}

func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	var req CreateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid application data", err)
		return
	}

	now := time.Now()
	var catID int
	if req.Category != nil {
		catID = *req.Category
	}

	calcPriority := h.priorityService.CalculatePriority(service.PriorityInput{
		LevelCriticalityScore: req.LevelCriticalityScore,
		CategoryID:            catID,
		BuildingID:            req.BuildingID,
		Description:           req.Description,
		UploadIDs:             req.UploadID,
		SubmissionTime:        now,
	})

	if req.Priority != nil && *req.Priority > 0 {
		calcPriority = *req.Priority
	}

	app := model.Application{
		Name:           req.Name,
		Description:    req.Description,
		Email:          req.Email,
		DateSubmission: now,
		Status:         model.StatusNew,
		Priority:       calcPriority,
		BuildingID:     req.BuildingID,
		UploadID:       req.UploadID,
		CategoryID:     req.Category,
	}

	created, err := h.appRepo.Create(&app)
	if err != nil {
		response.Internal(c, "Failed to create application", err)
		return
	}

	_ = h.buildingRepo.IncrementApplications(req.BuildingID)

	c.JSON(http.StatusCreated, created)
}

func (h *ApplicationHandler) UpdateApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid application ID", err)
		return
	}

	existing, err := h.appRepo.GetByID(id)
	if err != nil {
		response.Internal(c, "Database error", err)
		return
	}
	if existing == nil {
		response.NotFound(c, "Application not found")
		return
	}

	var updateReq struct {
		Name           *string                  `json:"name"`
		Description    *string                  `json:"description"`
		Email          *string                  `json:"email"`
		Status         *model.ApplicationStatus `json:"status"`
		Priority       *int                     `json:"priority"`
		BuildingID     *int                     `json:"building_id"`
		UploadID       []int                    `json:"upload_id"`
		UserID         *int                     `json:"user_id"`
		Category       *int                     `json:"category"`
		Reason         *string                  `json:"reason"`
		DateInProgress *time.Time               `json:"dateInProgress"`
		DateResult     *time.Time               `json:"dateResult"`
	}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		response.BadRequest(c, "Invalid update data", err)
		return
	}

	if updateReq.Name != nil {
		existing.Name = *updateReq.Name
	}
	if updateReq.Description != nil {
		existing.Description = *updateReq.Description
	}
	if updateReq.Email != nil {
		existing.Email = *updateReq.Email
	}
	if updateReq.Priority != nil {
		existing.Priority = *updateReq.Priority
	}
	if updateReq.BuildingID != nil {
		existing.BuildingID = *updateReq.BuildingID
	}
	if updateReq.UploadID != nil {
		existing.UploadID = updateReq.UploadID
	}
	if updateReq.UserID != nil {
		existing.UserID = updateReq.UserID
	}
	if updateReq.Category != nil {
		existing.CategoryID = updateReq.Category
	}
	if updateReq.Reason != nil {
		existing.Reason = updateReq.Reason
	}

	if updateReq.Status != nil {
		newStatus := *updateReq.Status
		existing.Status = newStatus
		now := time.Now()
		if newStatus == model.StatusInProgress && existing.DateInProgress == nil {
			existing.DateInProgress = &now
		} else if (newStatus == model.StatusCompleted || newStatus == model.StatusRejected) && existing.DateResult == nil {
			existing.DateResult = &now
		}
	}

	if updateReq.DateInProgress != nil {
		existing.DateInProgress = updateReq.DateInProgress
	}
	if updateReq.DateResult != nil {
		existing.DateResult = updateReq.DateResult
	}

	updated, err := h.appRepo.Update(existing)
	if err != nil {
		response.Internal(c, "Failed to update application", err)
		return
	}

	c.JSON(http.StatusOK, updated)
}
