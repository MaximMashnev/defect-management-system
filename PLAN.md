# План реализации проекта Defect Management System

> **Цель проекта:** Создать полноценную систему управления заявками строительной компании на основе существующего React-приложения, сохранив текущий репозиторий `test-task-react` в качестве эталонной истории, разработать с нуля backend на **Go (Golang)** и объединить их в новый монорепозиторий `defect-management-system`.

---

## 1. Архитектура и структура монорепозитория

Проект будет расположен в директории `e:\projects\defect-management-system` со следующей структурой:

```text
defect-management-system/
├── backend/                     # Go (Golang) REST API сервис
│   ├── cmd/
│   │   └── api/
│   │       └── main.go          # Точка входа в приложение
│   ├── internal/
│   │   ├── config/              # Загрузка конфигурации из .env / окружения
│   │   ├── handler/             # HTTP контроллеры (Gin handlers)
│   │   │   ├── auth.go
│   │   │   ├── building.go
│   │   │   ├── application.go
│   │   │   ├── analytics.go
│   │   │   └── upload.go
│   │   ├── service/             # Бизнес-логика (скоринг приоритета, валидация)
│   │   ├── repository/          # Слой доступа к данным (PostgreSQL + pgx/sqlx)
│   │   ├── model/               # Доменные структуры данных
│   │   └── middleware/          # JWT auth, CORS, логгер, обработка ошибок
│   ├── pkg/                     # Переиспользуемые утилиты (hash, jwt, response)
│   ├── migrations/              # SQL файлы миграций (goose / golang-migrate)
│   ├── uploads/                 # Локальное хранилище загруженных файлов
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/                    # React 19 + TypeScript + MUI + MobX
│   ├── public/
│   ├── src/                     # FSD архитектура (app, pages, features, entities, shared)
│   ├── .env.development
│   ├── .env.production
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yml           # Поднятие Go API + PostgreSQL + React + Nginx (опционально)
├── .gitignore
└── README.md
```

---

## 2. Стек технологий

| Компонент | Технологии |
|---|---|
| **Backend** | Go 1.23+, **Gin** (HTTP роутер), **pgx/sqlx** (доступ к PostgreSQL), **golang-jwt** (JWT аутентификация), **bcrypt** (хеширование паролей) |
| **Database** | **PostgreSQL 16** (хранение объектов, заявок, категорий, пользователей, метаданных файлов) |
| **Frontend** | React 19, TypeScript, **Material-UI (MUI)**, **MobX**, React Router v6, Axios |
| **DevOps** | Docker, Docker Compose, Git |

---

## 3. Этапы реализации

### Этап 1. Инициализация проекта и окружения
- [ ] Установка и верификация компилятора Go.
- [ ] Создание папки проекта `e:\projects\defect-management-system` и инициализация Git (`git init`).
- [ ] Копирование исходного кода фронтенда из `test-task-react` в папку `frontend/` (без `.git` и `node_modules` для чистоты).
- [ ] Настройка `.gitignore` для монорепозитория (Go бинарники, `node_modules`, `uploads/`, `.env`).
- [ ] Инициализация Go модуля: `go mod init defect-management-system/backend`.

### Этап 2. Проектирование базы данных и миграции
- [ ] Создание схемы БД:
  - `users` (id, login, password_hash, role, created_at)
  - `buildings` (id, name, address, date_registration, number_applications, created_at)
  - `categories` (id, name, score)
  - `applications` (id, name, description, email, date_submission, date_in_progress, date_result, status, priority, building_id, user_id, category_id, reason)
  - `files` (id, application_id, original_name, stored_name, file_path, size, mime_type, created_at)
- [ ] Написание начальной миграции `00001_init_schema.sql`.
- [ ] Скрипт сидирования данных (seed):
  - Дефолтный менеджер: `admin@a.a` / `admin123`.
  - Начальные категории и базовые тестовые объекты.

### Этап 3. Разработка Go бэкенда (Clean Architecture)
- [ ] **Конфигурация и инфраструктура:**
  - Загрузка переменных окружения (`PORT`, `DATABASE_URL`, `JWT_SECRET`, `UPLOAD_DIR`).
  - Подключение к PostgreSQL с пулом соединений (`pgxpool`).
  - CORS middleware для связи с фронтендом.
- [ ] **Модуль аутентификации:**
  - `POST /api/auth/login` (проверка bcrypt, генерация JWT токена).
  - `GET /api/auth/me` (проверка Bearer токена, отдача профиля).
- [ ] **Модуль зданий/объектов (`buildings`):**
  - CRUD: список с пагинацией, поиском по адресу/названию, создание, редактирование, удаление.
- [ ] **Модуль категорий (`category`):**
  - Получение категорий с их скорингом.
- [ ] **Модуль заявок (`applications`):**
  - Создание заявки клиентом (без авторизации).
  - Автоматический расчет скоринга приоритета (1–10) на бэкенде.
  - Получение одной заявки по ID / UUID (публичный трекер).
  - Получение списка для менеджера: фильтрация (статус, критичность, объект), сортировка, пагинация (совместимо с MUI DataGrid).
  - Смена статуса заявки менеджером (новый -> в работе -> выполнено / отклонено с фиксацией дат и причины).
- [ ] **Модуль файлов (`uploads`):**
  - Загрузка файлов (изображения, документы) с валидацией размера и типа.
  - Скачивание / отдача файлов.
- [ ] **Модуль аналитики (`analytics`):**
  - Эндпоинт для дашборда: агрегация количества заявок по статусам, распределение по объектам, динамика по датам.

### Этап 4. Адаптация и улучшение Frontend
- [ ] Переключение `config.json` / `.env` на локальный Go бэкенд (`http://localhost:8080/api`).
- [ ] Проверка и адаптация DTO: согласование форматов пагинации (`page`, `pageSize`, `total`) и структуры ответов.
- [ ] Улучшение UX:
  - Индикация загрузки и алерты (Snackbar) при ошибках и успешных операциях.
  - Проверка валидации формы отправки заявки.
  - Отображение превью файлов и скачивание.

### Этап 5. Docker Compose и финализация
- [ ] Написание `Dockerfile` для Go сервиса (многоэтапная сборка `golang:alpine`).
- [ ] Написание `Dockerfile` для React фронтенда (или dev-сервер через Compose).
- [ ] Настройка `docker-compose.yml` (Postgres + Backend + Frontend).
- [ ] Создание нового репозитория на GitHub и первоначальный коммит.
