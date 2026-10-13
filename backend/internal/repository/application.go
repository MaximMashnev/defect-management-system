package repository

import (
	"database/sql"
	"defect-management-system/backend/internal/model"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type ApplicationRepository struct {
	db *DB
}

func NewApplicationRepository(db *DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) GetAll() ([]model.Application, error) {
	rows, err := r.db.Query(`
		SELECT id, name, description, email, date_submission, date_in_progress, date_result,
		       status, priority, building_id, user_id, category_id, reason
		FROM applications ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []model.Application
	for rows.Next() {
		var a model.Application
		var st string
		if err := rows.Scan(
			&a.ID, &a.Name, &a.Description, &a.Email, &a.DateSubmission, &a.DateInProgress, &a.DateResult,
			&st, &a.Priority, &a.BuildingID, &a.UserID, &a.CategoryID, &a.Reason,
		); err != nil {
			return nil, err
		}
		a.Status = model.ApplicationStatus(st)
		a.UploadID = r.getUploadIDs(a.ID)
		apps = append(apps, a)
	}
	return apps, nil
}

func (r *ApplicationRepository) GetByID(id int) (*model.Application, error) {
	query := `
		SELECT id, name, description, email, date_submission, date_in_progress, date_result,
		       status, priority, building_id, user_id, category_id, reason
		FROM applications WHERE id = ?
	`
	if r.db.Driver != "sqlite" {
		query = `
			SELECT id, name, description, email, date_submission, date_in_progress, date_result,
			       status, priority, building_id, user_id, category_id, reason
			FROM applications WHERE id = $1
		`
	}

	var a model.Application
	var st string
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.Name, &a.Description, &a.Email, &a.DateSubmission, &a.DateInProgress, &a.DateResult,
		&st, &a.Priority, &a.BuildingID, &a.UserID, &a.CategoryID, &a.Reason,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	a.Status = model.ApplicationStatus(st)
	a.UploadID = r.getUploadIDs(a.ID)
	return &a, nil
}

func (r *ApplicationRepository) GetPaginated(page, limit int, filters map[string]string, sortBy string) (*model.PaginatedApplications, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var whereParts []string
	var args []interface{}

	for k, v := range filters {
		if v == "" {
			continue
		}
		var col string
		switch k {
		case "status":
			col = "status"
		case "building_id":
			col = "building_id"
		case "category", "category_id":
			col = "category_id"
		case "email":
			col = "email"
		case "priority":
			col = "priority"
		case "name":
			col = "name"
		default:
			continue
		}

		if r.db.Driver == "sqlite" {
			whereParts = append(whereParts, fmt.Sprintf("%s = ?", col))
			args = append(args, v)
		} else {
			whereParts = append(whereParts, fmt.Sprintf("%s = $%d", col, len(args)+1))
			args = append(args, v)
		}
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM applications %s", whereClause)
	var totalItems int
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalItems); err != nil {
		return nil, err
	}

	orderClause := "ORDER BY id DESC"
	if sortBy != "" {
		dir := "ASC"
		field := sortBy
		if strings.HasPrefix(sortBy, "-") {
			dir = "DESC"
			field = sortBy[1:]
		}
		switch field {
		case "id":
			orderClause = fmt.Sprintf("ORDER BY id %s", dir)
		case "name":
			orderClause = fmt.Sprintf("ORDER BY name %s", dir)
		case "status":
			orderClause = fmt.Sprintf("ORDER BY status %s", dir)
		case "priority":
			orderClause = fmt.Sprintf("ORDER BY priority %s", dir)
		case "dateSubmission":
			orderClause = fmt.Sprintf("ORDER BY date_submission %s", dir)
		case "building_id":
			orderClause = fmt.Sprintf("ORDER BY building_id %s", dir)
		case "category":
			orderClause = fmt.Sprintf("ORDER BY category_id %s", dir)
		}
	}

	var selectQuery string
	if r.db.Driver == "sqlite" {
		selectQuery = fmt.Sprintf(`
			SELECT id, name, description, email, date_submission, date_in_progress, date_result,
			       status, priority, building_id, user_id, category_id, reason
			FROM applications %s %s LIMIT ? OFFSET ?
		`, whereClause, orderClause)
		args = append(args, limit, offset)
	} else {
		n := len(args)
		selectQuery = fmt.Sprintf(`
			SELECT id, name, description, email, date_submission, date_in_progress, date_result,
			       status, priority, building_id, user_id, category_id, reason
			FROM applications %s %s LIMIT $%d OFFSET $%d
		`, whereClause, orderClause, n+1, n+2)
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(selectQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []model.Application
	for rows.Next() {
		var a model.Application
		var st string
		if err := rows.Scan(
			&a.ID, &a.Name, &a.Description, &a.Email, &a.DateSubmission, &a.DateInProgress, &a.DateResult,
			&st, &a.Priority, &a.BuildingID, &a.UserID, &a.CategoryID, &a.Reason,
		); err != nil {
			return nil, err
		}
		a.Status = model.ApplicationStatus(st)
		a.UploadID = r.getUploadIDs(a.ID)
		apps = append(apps, a)
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	remaining := totalItems - (page * limit)
	if remaining < 0 {
		remaining = 0
	}

	return &model.PaginatedApplications{
		Meta: model.Meta{
			CurrentPage:    page,
			PerPage:        limit,
			RemainingCount: remaining,
			TotalItems:     totalItems,
			TotalPages:     totalPages,
		},
		Items: apps,
	}, nil
}

func (r *ApplicationRepository) Create(a *model.Application) (*model.Application, error) {
	if a.DateSubmission.IsZero() {
		a.DateSubmission = time.Now()
	}
	if a.Status == "" {
		a.Status = model.StatusNew
	}

	if r.db.Driver == "sqlite" {
		res, err := r.db.Exec(`
			INSERT INTO applications (name, description, email, date_submission, date_in_progress, date_result,
			                          status, priority, building_id, user_id, category_id, reason)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, a.Name, a.Description, a.Email, a.DateSubmission, a.DateInProgress, a.DateResult,
			string(a.Status), a.Priority, a.BuildingID, a.UserID, a.CategoryID, a.Reason)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		a.ID = int(id)
	} else {
		err := r.db.QueryRow(`
			INSERT INTO applications (name, description, email, date_submission, date_in_progress, date_result,
			                          status, priority, building_id, user_id, category_id, reason)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			RETURNING id
		`, a.Name, a.Description, a.Email, a.DateSubmission, a.DateInProgress, a.DateResult,
			string(a.Status), a.Priority, a.BuildingID, a.UserID, a.CategoryID, a.Reason).Scan(&a.ID)
		if err != nil {
			return nil, err
		}
	}

	if a.UploadID == nil {
		a.UploadID = []int{}
	}
	if len(a.UploadID) > 0 {
		_ = r.saveUploadIDs(a.ID, a.UploadID)
	}

	return a, nil
}

func (r *ApplicationRepository) Update(a *model.Application) (*model.Application, error) {
	if r.db.Driver == "sqlite" {
		_, err := r.db.Exec(`
			UPDATE applications
			SET name = ?, description = ?, email = ?, date_in_progress = ?, date_result = ?,
			    status = ?, priority = ?, building_id = ?, user_id = ?, category_id = ?, reason = ?
			WHERE id = ?
		`, a.Name, a.Description, a.Email, a.DateInProgress, a.DateResult,
			string(a.Status), a.Priority, a.BuildingID, a.UserID, a.CategoryID, a.Reason, a.ID)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := r.db.Exec(`
			UPDATE applications
			SET name = $1, description = $2, email = $3, date_in_progress = $4, date_result = $5,
			    status = $6, priority = $7, building_id = $8, user_id = $9, category_id = $10, reason = $11
			WHERE id = $12
		`, a.Name, a.Description, a.Email, a.DateInProgress, a.DateResult,
			string(a.Status), a.Priority, a.BuildingID, a.UserID, a.CategoryID, a.Reason, a.ID)
		if err != nil {
			return nil, err
		}
	}

	if a.UploadID == nil {
		a.UploadID = []int{}
	}
	if len(a.UploadID) > 0 {
		_ = r.saveUploadIDs(a.ID, a.UploadID)
	}

	return a, nil
}

func (r *ApplicationRepository) getUploadIDs(appID int) []int {
	ids := make([]int, 0)
	query := "SELECT upload_id FROM application_uploads WHERE application_id = ?"
	if r.db.Driver != "sqlite" {
		query = "SELECT upload_id FROM application_uploads WHERE application_id = $1"
	}
	rows, err := r.db.Query(query, appID)
	if err != nil {
		return ids
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

func (r *ApplicationRepository) saveUploadIDs(appID int, uploadIDs []int) error {
	delQuery := "DELETE FROM application_uploads WHERE application_id = ?"
	if r.db.Driver != "sqlite" {
		delQuery = "DELETE FROM application_uploads WHERE application_id = $1"
	}
	_, _ = r.db.Exec(delQuery, appID)

	for _, uid := range uploadIDs {
		if r.db.Driver == "sqlite" {
			_, _ = r.db.Exec("INSERT INTO application_uploads (application_id, upload_id) VALUES (?, ?)", appID, uid)
			_, _ = r.db.Exec("UPDATE uploads SET application_id = ? WHERE id = ?", appID, uid)
		} else {
			_, _ = r.db.Exec("INSERT INTO application_uploads (application_id, upload_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", appID, uid)
			_, _ = r.db.Exec("UPDATE uploads SET application_id = $1 WHERE id = $2", appID, uid)
		}
	}
	return nil
}
