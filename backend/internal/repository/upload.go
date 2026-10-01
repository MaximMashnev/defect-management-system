package repository

import (
	"database/sql"
	"defect-management-system/backend/internal/model"
	"errors"
	"fmt"
	"strings"
)

type UploadRepository struct {
	db *DB
}

func NewUploadRepository(db *DB) *UploadRepository {
	return &UploadRepository{db: db}
}

func (r *UploadRepository) Save(u *model.UploadFile) (*model.UploadFile, error) {
	if r.db.Driver == "sqlite" {
		res, err := r.db.Exec(`
			INSERT INTO uploads (original_name, stored_name, file_path, url, size, mime_type)
			VALUES (?, ?, ?, ?, ?, ?)
		`, u.OriginalName, u.StoredName, u.FilePath, u.URL, u.Bytes, u.MimeType)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		u.ID = int(id)
	} else {
		err := r.db.QueryRow(`
			INSERT INTO uploads (original_name, stored_name, file_path, url, size, mime_type)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, u.OriginalName, u.StoredName, u.FilePath, u.URL, u.Bytes, u.MimeType).Scan(&u.ID)
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (r *UploadRepository) GetByID(id int) (*model.UploadFile, error) {
	query := `
		SELECT id, original_name, stored_name, file_path, url, size, mime_type, created_at
		FROM uploads WHERE id = ?
	`
	if r.db.Driver != "sqlite" {
		query = `
			SELECT id, original_name, stored_name, file_path, url, size, mime_type, created_at
			FROM uploads WHERE id = $1
		`
	}

	var u model.UploadFile
	err := r.db.QueryRow(query, id).Scan(
		&u.ID, &u.OriginalName, &u.StoredName, &u.FilePath, &u.URL, &u.Bytes, &u.MimeType, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	u.Format = u.MimeType
	return &u, nil
}

func (r *UploadRepository) GetByIDs(ids []int) ([]model.UploadFile, error) {
	if len(ids) == 0 {
		return []model.UploadFile{}, nil
	}

	var placeholders []string
	var args []interface{}
	for i, id := range ids {
		if r.db.Driver == "sqlite" {
			placeholders = append(placeholders, "?")
		} else {
			placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		}
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT id, original_name, stored_name, file_path, url, size, mime_type, created_at
		FROM uploads WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uploads []model.UploadFile
	for rows.Next() {
		var u model.UploadFile
		if err := rows.Scan(
			&u.ID, &u.OriginalName, &u.StoredName, &u.FilePath, &u.URL, &u.Bytes, &u.MimeType, &u.CreatedAt,
		); err != nil {
			return nil, err
		}
		u.Format = u.MimeType
		uploads = append(uploads, u)
	}
	return uploads, nil
}
