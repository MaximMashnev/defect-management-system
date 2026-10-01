package handler

import (
	"defect-management-system/backend/internal/model"
	"defect-management-system/backend/internal/repository"
	"defect-management-system/backend/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	appRepo      *repository.ApplicationRepository
	buildingRepo *repository.BuildingRepository
}

func NewAnalyticsHandler(
	appRepo *repository.ApplicationRepository,
	buildingRepo *repository.BuildingRepository,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		appRepo:      appRepo,
		buildingRepo: buildingRepo,
	}
}

type AnalyticsResponse struct {
	TotalApplications int              `json:"totalApplications"`
	TotalBuildings    int              `json:"totalBuildings"`
	ByStatus          map[string]int   `json:"byStatus"`
	ByBuilding        []BuildingStat   `json:"byBuilding"`
	ByHour            map[string]int   `json:"byHour"`
}

type BuildingStat struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (h *AnalyticsHandler) GetAnalytics(c *gin.Context) {
	apps, err := h.appRepo.GetAll()
	if err != nil {
		response.Internal(c, "Failed to get applications", err)
		return
	}

	buildings, err := h.buildingRepo.GetAll()
	if err != nil {
		response.Internal(c, "Failed to get buildings", err)
		return
	}

	byStatus := map[string]int{
		string(model.StatusNew):        0,
		string(model.StatusInProgress): 0,
		string(model.StatusCompleted):  0,
		string(model.StatusRejected):   0,
	}

	byHour := make(map[string]int)
	for _, app := range apps {
		byStatus[string(app.Status)]++
		hourStr := strconv.Itoa(app.DateSubmission.Hour())
		byHour[hourStr]++
	}

	var byBuilding []BuildingStat
	for _, b := range buildings {
		byBuilding = append(byBuilding, BuildingStat{
			Name:  b.Name,
			Count: b.NumberApplications,
		})
	}

	c.JSON(http.StatusOK, AnalyticsResponse{
		TotalApplications: len(apps),
		TotalBuildings:    len(buildings),
		ByStatus:          byStatus,
		ByBuilding:        byBuilding,
		ByHour:            byHour,
	})
}
