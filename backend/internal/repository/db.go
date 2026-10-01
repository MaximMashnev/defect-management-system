package repository

import (
	"database/sql"
	"defect-management-system/backend/pkg/hash"
	"fmt"
	"log"

	_ "github.com/glebarez/go-sqlite"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	*sql.DB
	Driver string
}

func NewDB(driver, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	repoDB := &DB{DB: db, Driver: driver}
	if err := repoDB.InitSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repoDB, nil
}

func (db *DB) InitSchema() error {
	var schema string
	if db.Driver == "sqlite" {
		schema = `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			login TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'manager',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS buildings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			address TEXT NOT NULL,
			date_registration DATETIME DEFAULT CURRENT_TIMESTAMP,
			number_applications INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			score INTEGER NOT NULL DEFAULT 1
		);

		CREATE TABLE IF NOT EXISTS applications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			email TEXT NOT NULL,
			date_submission DATETIME DEFAULT CURRENT_TIMESTAMP,
			date_in_progress DATETIME,
			date_result DATETIME,
			status TEXT NOT NULL DEFAULT 'новый',
			priority INTEGER NOT NULL DEFAULT 1,
			building_id INTEGER NOT NULL,
			user_id INTEGER,
			category_id INTEGER,
			reason TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS uploads (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			application_id INTEGER,
			original_name TEXT NOT NULL,
			stored_name TEXT NOT NULL,
			file_path TEXT NOT NULL,
			url TEXT NOT NULL,
			size INTEGER NOT NULL,
			mime_type TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS application_uploads (
			application_id INTEGER NOT NULL,
			upload_id INTEGER NOT NULL,
			PRIMARY KEY (application_id, upload_id)
		);
		`
	} else {
		schema = `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'manager',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS buildings (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			address TEXT NOT NULL,
			date_registration TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			number_applications INT NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			score INT NOT NULL DEFAULT 1
		);

		CREATE TABLE IF NOT EXISTS applications (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			email VARCHAR(255) NOT NULL,
			date_submission TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			date_in_progress TIMESTAMPTZ,
			date_result TIMESTAMPTZ,
			status VARCHAR(50) NOT NULL DEFAULT 'новый',
			priority INT NOT NULL DEFAULT 1,
			building_id INT NOT NULL,
			user_id INT,
			category_id INT,
			reason TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS uploads (
			id SERIAL PRIMARY KEY,
			application_id INT,
			original_name VARCHAR(255) NOT NULL,
			stored_name VARCHAR(255) NOT NULL,
			file_path TEXT NOT NULL,
			url TEXT NOT NULL,
			size BIGINT NOT NULL,
			mime_type VARCHAR(100) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS application_uploads (
			application_id INT NOT NULL,
			upload_id INT NOT NULL,
			PRIMARY KEY (application_id, upload_id)
		);
		`
	}

	if _, err := db.Exec(schema); err != nil {
		return err
	}

	return db.seedData()
}

func (db *DB) seedData() error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE login = 'admin@a.a'").Scan(&count)
	if err == nil && count == 0 {
		pwdHash, _ := hash.HashPassword("admin123")
		if db.Driver == "sqlite" {
			_, _ = db.Exec("INSERT INTO users (login, password_hash, role) VALUES (?, ?, ?)", "admin@a.a", pwdHash, "manager")
		} else {
			_, _ = db.Exec("INSERT INTO users (login, password_hash, role) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", "admin@a.a", pwdHash, "manager")
		}
	} else if err == nil && count > 0 {
		// Update password hash to make sure it's valid admin123
		pwdHash, _ := hash.HashPassword("admin123")
		if db.Driver == "sqlite" {
			_, _ = db.Exec("UPDATE users SET password_hash = ? WHERE login = 'admin@a.a'", pwdHash)
		} else {
			_, _ = db.Exec("UPDATE users SET password_hash = $1 WHERE login = 'admin@a.a'", pwdHash)
		}
	}

	var catCount int
	err = db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&catCount)
	if err == nil && catCount == 0 {
		categories := []struct {
			name  string
			score int
		}{
			{"Водоснабжение и отопление", 3},
			{"Электрика", 3},
			{"Окна и двери", 2},
			{"Вентиляция", 2},
			{"Отделочные работы", 1},
			{"Фасад и кровля", 2},
		}
		for _, cat := range categories {
			if db.Driver == "sqlite" {
				_, _ = db.Exec("INSERT INTO categories (name, score) VALUES (?, ?)", cat.name, cat.score)
			} else {
				_, _ = db.Exec("INSERT INTO categories (name, score) VALUES ($1, $2)", cat.name, cat.score)
			}
		}
	}

	var buildCount int
	err = db.QueryRow("SELECT COUNT(*) FROM buildings").Scan(&buildCount)
	if err == nil && buildCount == 0 {
		buildings := []struct {
			name    string
			address string
			num     int
		}{
			{"ЖК \"Северное Сияние\", Корпус 1", "г. Москва, ул. Полярная, д. 14", 3},
			{"ЖК \"Северное Сияние\", Корпус 2", "г. Москва, ул. Полярная, д. 16", 2},
			{"ЖК \"Лазурный Берег\", Башня А", "г. Санкт-Петербург, Приморский пр., д. 45", 4},
			{"БЦ \"Горизонт\"", "г. Москва, Ленинградский пр-т, д. 80", 1},
		}
		for _, b := range buildings {
			if db.Driver == "sqlite" {
				_, _ = db.Exec("INSERT INTO buildings (name, address, number_applications) VALUES (?, ?, ?)", b.name, b.address, b.num)
			} else {
				_, _ = db.Exec("INSERT INTO buildings (name, address, number_applications) VALUES ($1, $2, $3)", b.name, b.address, b.num)
			}
		}
	}

	log.Println("Database schema initialized and seed verified.")
	return nil
}
