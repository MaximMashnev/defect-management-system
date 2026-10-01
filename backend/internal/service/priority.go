package service

import (
	"defect-management-system/backend/internal/repository"
	"strings"
	"time"
)

type PriorityService struct {
	appRepo      *repository.ApplicationRepository
	buildingRepo *repository.BuildingRepository
	catRepo      *repository.CategoryRepository
	uploadRepo   *repository.UploadRepository
}

func NewPriorityService(
	appRepo *repository.ApplicationRepository,
	buildingRepo *repository.BuildingRepository,
	catRepo *repository.CategoryRepository,
	uploadRepo *repository.UploadRepository,
) *PriorityService {
	return &PriorityService{
		appRepo:      appRepo,
		buildingRepo: buildingRepo,
		catRepo:      catRepo,
		uploadRepo:   uploadRepo,
	}
}

type PriorityInput struct {
	LevelCriticalityScore int
	CategoryID            int
	BuildingID            int
	Description           string
	UploadIDs             []int
	SubmissionTime        time.Time
}

func (s *PriorityService) CalculatePriority(input PriorityInput) int {
	priority := input.LevelCriticalityScore

	if input.CategoryID > 0 {
		cat, err := s.catRepo.GetByID(input.CategoryID)
		if err == nil && cat != nil {
			priority += cat.Score
		}
	}

	descLen := len(strings.TrimSpace(input.Description))
	if descLen >= 200 {
		priority += 2
	} else if descLen >= 50 {
		priority += 1
	}

	if len(input.UploadIDs) > 0 {
		if len(input.UploadIDs) >= 3 {
			priority += 2
		} else {
			priority += 1
		}

		files, err := s.uploadRepo.GetByIDs(input.UploadIDs)
		if err == nil {
			for _, f := range files {
				if strings.HasPrefix(f.MimeType, "video/") {
					priority += 1
					break
				}
			}
		}
	}

	t := input.SubmissionTime
	if t.IsZero() {
		t = time.Now()
	}
	hour := t.Hour()
	if hour < 6 || hour >= 22 {
		priority += 2
	}

	allApps, err := s.appRepo.GetAll()
	if err == nil {
		for _, a := range allApps {
			if a.BuildingID == input.BuildingID && a.CategoryID != nil && *a.CategoryID == input.CategoryID {
				priority += 2
				break
			}
		}
	}

	currentBuilding, err := s.buildingRepo.GetByID(input.BuildingID)
	if err == nil && currentBuilding != nil {
		allBuildings, bErr := s.buildingRepo.GetAll()
		if bErr == nil {
			hasHigher := false
			for _, b := range allBuildings {
				if b.NumberApplications > currentBuilding.NumberApplications {
					hasHigher = true
					break
				}
			}
			if !hasHigher {
				priority += 1
			}
		}
	}

	if priority > 10 {
		priority = 10
	}
	if priority < 1 {
		priority = 1
	}

	return priority
}
