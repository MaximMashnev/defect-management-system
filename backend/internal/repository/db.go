package repository

import (
	"database/sql"
	"defect-management-system/backend/pkg/hash"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"

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

	// Seed sample uploads if empty
	var uploadCount int
	err = db.QueryRow("SELECT COUNT(*) FROM uploads").Scan(&uploadCount)
	if err == nil && uploadCount == 0 {
		_ = os.MkdirAll("uploads", 0755)
		sampleFiles := []struct {
			filename string
			name     string
			color    color.RGBA
		}{
			{"sample_pipe.png", "Дефект_трубы_отопления.png", color.RGBA{R: 52, G: 152, B: 219, A: 255}},
			{"sample_electric.png", "Щиток_искрение.png", color.RGBA{R: 231, G: 76, B: 60, A: 255}},
			{"sample_facade.png", "Трещина_фасада.png", color.RGBA{R: 243, G: 156, B: 18, A: 255}},
		}

		for _, sf := range sampleFiles {
			filePath := filepath.Join("uploads", sf.filename)
			createSampleImage(filePath, sf.color)
			fileURL := fmt.Sprintf("http://localhost:8080/api/uploads/files/%s", sf.filename)
			fi, fErr := os.Stat(filePath)
			var size int64 = 1024
			if fErr == nil {
				size = fi.Size()
			}

			if db.Driver == "sqlite" {
				_, _ = db.Exec(`
					INSERT INTO uploads (original_name, stored_name, file_path, url, size, mime_type)
					VALUES (?, ?, ?, ?, ?, ?)
				`, sf.name, sf.filename, filePath, fileURL, size, "image/png")
			} else {
				_, _ = db.Exec(`
					INSERT INTO uploads (original_name, stored_name, file_path, url, size, mime_type)
					VALUES ($1, $2, $3, $4, $5, $6)
				`, sf.name, sf.filename, filePath, fileURL, size, "image/png")
			}
		}
	}

	// Seed applications if empty
	var appCount int
	err = db.QueryRow("SELECT COUNT(*) FROM applications").Scan(&appCount)
	if err == nil && appCount == 0 {
		now := time.Now()
		tMinus7d := now.Add(-7 * 24 * time.Hour)
		tMinus6d := now.Add(-6 * 24 * time.Hour)
		tMinus5d := now.Add(-5 * 24 * time.Hour)
		tMinus4d := now.Add(-4 * 24 * time.Hour)
		tMinus3d := now.Add(-3 * 24 * time.Hour)
		tMinus2d := now.Add(-2 * 24 * time.Hour)
		tMinus1d := now.Add(-1 * 24 * time.Hour)
		tMinus12h := now.Add(-12 * time.Hour)
		tMinus4h := now.Add(-4 * time.Hour)

		rejectReason := "Повреждение возникло в результате нарушения жильцами правил переноса тяжелой мебели при переезде."

		type appSeed struct {
			Name           string
			Description    string
			Email          string
			DateSubmission time.Time
			DateInProgress *time.Time
			DateResult     *time.Time
			Status         string
			Priority       int
			BuildingID     int
			CategoryID     int
			Reason         *string
			UploadIDs      []int
		}

		seeds := []appSeed{
			{
				Name:           "Протечка стояка отопления на 5 этаже",
				Description:    "В районе межэтажного перекрытия между 5 и 6 этажами обнаружена активная течь трубы отопления. Вода капает на электрощиток. Требуется срочный вызов аварийной бригады сантехников.",
				Email:          "petrov.nord@example.com",
				DateSubmission: tMinus3d,
				Status:         "новый",
				Priority:       9,
				BuildingID:     1,
				CategoryID:     1,
				UploadIDs:      []int{1},
			},
			{
				Name:           "Искрит распределительный щиток в холле",
				Description:    "Периодически наблюдается треск и искрение в главном распределительном щите первого подъезда, чувствуется запах паленой изоляции.",
				Email:          "electric.guard@mail.ru",
				DateSubmission: tMinus2d,
				DateInProgress: &tMinus1d,
				Status:         "в работе",
				Priority:       8,
				BuildingID:     1,
				CategoryID:     2,
				UploadIDs:      []int{2},
			},
			{
				Name:           "Продувание оконного блока в квартире 42",
				Description:    "Створка панорамного окна в гостиной не прижимается плотно по верхнему контуру, при сильном ветре слышен свист и продувание.",
				Email:          "kuznetsov42@gmail.com",
				DateSubmission: tMinus5d,
				DateInProgress: &tMinus4d,
				DateResult:     &tMinus1d,
				Status:         "выполнено",
				Priority:       4,
				BuildingID:     1,
				CategoryID:     3,
			},
			{
				Name:           "Не работает приточно-вытяжная вентиляция",
				Description:    "Обратная тяга в вентиляционной шахте кухни и санузла. Запахи из соседних квартир поступают в жилое помещение.",
				Email:          "smirnov.vent@yandex.ru",
				DateSubmission: tMinus1d,
				Status:         "новый",
				Priority:       6,
				BuildingID:     2,
				CategoryID:     4,
			},
			{
				Name:           "Трещины в штукатурке входной группы",
				Description:    "Повреждение декоративного штукатурного слоя и сколы углов возле лифтового холла на первом этаже подъезда.",
				Email:          "volkov.info@mail.ru",
				DateSubmission: tMinus6d,
				DateInProgress: &tMinus5d,
				DateResult:     &tMinus4d,
				Status:         "отклонено",
				Reason:         &rejectReason,
				Priority:       3,
				BuildingID:     2,
				CategoryID:     5,
			},
			{
				Name:           "Затопление подвального помещения секции Б",
				Description:    "В подвале секции Б подъем грунтовых вод в связи с неисправностью дренажного насоса. Уровень воды около 10 см.",
				Email:          "spb.resident@gmail.com",
				DateSubmission: tMinus4h,
				Status:         "новый",
				Priority:       10,
				BuildingID:     3,
				CategoryID:     1,
				UploadIDs:      []int{1, 3},
			},
			{
				Name:           "Отслоение фасадной плитки на уровне 3 этажа",
				Description:    "Керамогранитная плитка вентилируемого фасада отошла от направляющей, есть риск падения плитки на пешеходную дорожку.",
				Email:          "facade.control@spb.ru",
				DateSubmission: tMinus2d,
				DateInProgress: &tMinus1d,
				Status:         "в работе",
				Priority:       8,
				BuildingID:     3,
				CategoryID:     6,
				UploadIDs:      []int{3},
			},
			{
				Name:           "Перебои с освещением подземного паркинга",
				Description:    "Мигают и гаснут светодиодные светильники в проезде между рядами 3 и 4 паркинга -1 уровня.",
				Email:          "parking.master@auto.org",
				DateSubmission: tMinus3d,
				DateInProgress: &tMinus2d,
				Status:         "в работе",
				Priority:       5,
				BuildingID:     3,
				CategoryID:     2,
			},
			{
				Name:           "Деформация балконной рамы после шторма",
				Description:    "После штормового ветра с залива перекосило створку холодного остекления балкона, не закрывается замок.",
				Email:          "sea.view15@spb.net",
				DateSubmission: tMinus7d,
				DateInProgress: &tMinus6d,
				DateResult:     &tMinus2d,
				Status:         "выполнено",
				Priority:       5,
				BuildingID:     3,
				CategoryID:     3,
			},
			{
				Name:           "Шум и вибрация в шахте центрального кондиционирования",
				Description:    "На 8 этаже бизнес-центра от воздуховода исходит непрерывный низкочастотный гул и вибрация, мешающая работе сотрудников.",
				Email:          "office.horizon@corp.biz",
				DateSubmission: tMinus12h,
				Status:         "новый",
				Priority:       7,
				BuildingID:     4,
				CategoryID:     4,
			},
		}

		for _, s := range seeds {
			var appID int
			if db.Driver == "sqlite" {
				res, insErr := db.Exec(`
					INSERT INTO applications (name, description, email, date_submission, date_in_progress, date_result,
					                          status, priority, building_id, category_id, reason)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				`, s.Name, s.Description, s.Email, s.DateSubmission, s.DateInProgress, s.DateResult,
					s.Status, s.Priority, s.BuildingID, s.CategoryID, s.Reason)
				if insErr == nil {
					lid, _ := res.LastInsertId()
					appID = int(lid)
				}
			} else {
				_ = db.QueryRow(`
					INSERT INTO applications (name, description, email, date_submission, date_in_progress, date_result,
					                          status, priority, building_id, category_id, reason)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
					RETURNING id
				`, s.Name, s.Description, s.Email, s.DateSubmission, s.DateInProgress, s.DateResult,
					s.Status, s.Priority, s.BuildingID, s.CategoryID, s.Reason).Scan(&appID)
			}

			if appID > 0 && len(s.UploadIDs) > 0 {
				for _, uid := range s.UploadIDs {
					if db.Driver == "sqlite" {
						_, _ = db.Exec("INSERT INTO application_uploads (application_id, upload_id) VALUES (?, ?)", appID, uid)
						_, _ = db.Exec("UPDATE uploads SET application_id = ? WHERE id = ?", appID, uid)
					} else {
						_, _ = db.Exec("INSERT INTO application_uploads (application_id, upload_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", appID, uid)
						_, _ = db.Exec("UPDATE uploads SET application_id = $1 WHERE id = $2", appID, uid)
					}
				}
			}
		}
	}

	// Synchronize number_applications in buildings with real count
	if db.Driver == "sqlite" {
		_, _ = db.Exec(`
			UPDATE buildings
			SET number_applications = (
				SELECT COUNT(*) FROM applications WHERE applications.building_id = buildings.id
			)
		`)
	} else {
		_, _ = db.Exec(`
			UPDATE buildings
			SET number_applications = (
				SELECT COUNT(*) FROM applications WHERE applications.building_id = buildings.id
			)
		`)
	}

	log.Println("Database schema initialized and seed verified.")
	return nil
}

func createSampleImage(filePath string, c color.Color) {
	if _, err := os.Stat(filePath); err == nil {
		return
	}
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: c}, image.Point{}, draw.Src)
	f, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	_ = png.Encode(f, img)
}

