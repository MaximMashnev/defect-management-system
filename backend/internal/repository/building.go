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

type BuildingRepository struct {
	db *DB
}

func NewBuildingRepository(db *DB) *BuildingRepository {
	return &BuildingRepository{db: db}
}

func (r *BuildingRepository) GetAll() ([]model.Building, error) {
	rows, err := r.db.Query("SELECT id, name, address, date_registration, number_applications FROM buildings ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []model.Building
	for rows.Next() {
		var b model.Building
		if err := rows.Scan(&b.ID, &b.Name, &b.Address, &b.DateRegistration, &b.NumberApplications); err != nil {
			return nil, err
		}
		buildings = append(buildings, b)
	}
	return buildings, nil
}

func (r *BuildingRepository) GetByID(id int) (*model.Building, error) {
	query := "SELECT id, name, address, date_registration, number_applications FROM buildings WHERE id = ?"
	if r.db.Driver != "sqlite" {
		query = "SELECT id, name, address, date_registration, number_applications FROM buildings WHERE id = $1"
	}

	var b model.Building
	err := r.db.QueryRow(query, id).Scan(&b.ID, &b.Name, &b.Address, &b.DateRegistration, &b.NumberApplications)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

func (r *BuildingRepository) GetPaginated(page, limit int, search, sortBy string) (*model.PaginatedBuildings, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	whereClause := ""
	var args []interface{}
	if search != "" {
		if r.db.Driver == "sqlite" {
			whereClause = "WHERE name LIKE ? OR address LIKE ?"
			args = append(args, "%"+search+"%", "%"+search+"%")
		} else {
			whereClause = "WHERE name ILIKE $1 OR address ILIKE $2"
			args = append(args, "%"+search+"%", "%"+search+"%")
		}
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM buildings %s", whereClause)
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
		case "address":
			orderClause = fmt.Sprintf("ORDER BY address %s", dir)
		case "numberApplications":
			orderClause = fmt.Sprintf("ORDER BY number_applications %s", dir)
		case "dateRegistration":
			orderClause = fmt.Sprintf("ORDER BY date_registration %s", dir)
		}
	}

	var selectQuery string
	if r.db.Driver == "sqlite" {
		selectQuery = fmt.Sprintf("SELECT id, name, address, date_registration, number_applications FROM buildings %s %s LIMIT ? OFFSET ?", whereClause, orderClause)
		args = append(args, limit, offset)
	} else {
		n := len(args)
		selectQuery = fmt.Sprintf("SELECT id, name, address, date_registration, number_applications FROM buildings %s %s LIMIT $%d OFFSET $%d", whereClause, orderClause, n+1, n+2)
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(selectQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []model.Building
	for rows.Next() {
		var b model.Building
		if err := rows.Scan(&b.ID, &b.Name, &b.Address, &b.DateRegistration, &b.NumberApplications); err != nil {
			return nil, err
		}
		buildings = append(buildings, b)
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	remaining := totalItems - (page * limit)
	if remaining < 0 {
		remaining = 0
	}

	return &model.PaginatedBuildings{
		Meta: model.Meta{
			CurrentPage:    page,
			PerPage:        limit,
			RemainingCount: remaining,
			TotalItems:     totalItems,
			TotalPages:     totalPages,
		},
		Items: buildings,
	}, nil
}

func (r *BuildingRepository) Create(b *model.Building) (*model.Building, error) {
	if b.DateRegistration.IsZero() {
		b.DateRegistration = time.Now()
	}

	if r.db.Driver == "sqlite" {
		res, err := r.db.Exec("INSERT INTO buildings (name, address, date_registration, number_applications) VALUES (?, ?, ?, ?)",
			b.Name, b.Address, b.DateRegistration, b.NumberApplications)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		b.ID = int(id)
	} else {
		err := r.db.QueryRow("INSERT INTO buildings (name, address, date_registration, number_applications) VALUES ($1, $2, $3, $4) RETURNING id",
			b.Name, b.Address, b.DateRegistration, b.NumberApplications).Scan(&b.ID)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (r *BuildingRepository) Update(b *model.Building) (*model.Building, error) {
	if r.db.Driver == "sqlite" {
		_, err := r.db.Exec("UPDATE buildings SET name = ?, address = ?, number_applications = ? WHERE id = ?",
			b.Name, b.Address, b.NumberApplications, b.ID)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := r.db.Exec("UPDATE buildings SET name = $1, address = $2, number_applications = $3 WHERE id = $4",
			b.Name, b.Address, b.NumberApplications, b.ID)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (r *BuildingRepository) Delete(id int) error {
	query := "DELETE FROM buildings WHERE id = ?"
	if r.db.Driver != "sqlite" {
		query = "DELETE FROM buildings WHERE id = $1"
	}
	_, err := r.db.Exec(query, id)
	return err
}

func (r *BuildingRepository) IncrementApplications(id int) error {
	query := "UPDATE buildings SET number_applications = number_applications + 1 WHERE id = ?"
	if r.db.Driver != "sqlite" {
		query = "UPDATE buildings SET number_applications = number_applications + 1 WHERE id = $1"
	}
	_, err := r.db.Exec(query, id)
	return err
}
