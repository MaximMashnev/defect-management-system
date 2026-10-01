# Defect Management System

Полнофункциональная система для приёма, автоматической оценки и управления заявками на устранение строительных недочетов.

## 🚀 Архитектура и стек

- **Backend (Go 1.24+):**
  - Фреймворк: **Gin** (высокопроизводительный REST API).
  - База данных: **PostgreSQL** (продакшн / Docker) и встроенный **SQLite** (быстрый локальный запуск без сторонних сервисов).
  - Безопасность: **JWT** токены (HMAC-SHA256) + **bcrypt** для паролей.
  - Бизнес-логика: автоматический расчет приоритета заявок (1–10) на основе критичности, категории, объёма описания, прикрепленных медиафайлов, времени подачи и истории объекта.
  - Загрузка и раздача файлов (`/api/uploads`).
- **Frontend (React 19 + TypeScript):**
  - **Material-UI (MUI)**, **MobX**, **React Router v6**, **Axios**.
  - Архитектура **Feature-Sliced Design (FSD)**.
  - Роли: Публичный пользователь (подача заявки, трекинг по ID/ссылке) и Менеджер (вкладки "Объекты", "Заявки", "Аналитика").

---

## 📁 Структура проекта

```text
defect-management-system/
├── backend/                     # Серверная часть на Go
│   ├── cmd/api/main.go          # Точка входа в API
│   ├── internal/                # Clean architecture (handler, service, repository, model, middleware)
│   ├── pkg/                     # Утилиты (jwt, hash, response)
│   ├── migrations/              # SQL-миграции и сиды
│   ├── uploads/                 # Хранилище загруженных файлов
│   ├── Dockerfile
│   └── go.mod
├── frontend/                    # Клиентская часть на React
│   ├── src/                     # Исходный код (FSD)
│   ├── .env                     # Конфигурация API URL
│   └── package.json
├── docker-compose.yml           # Запуск PostgreSQL
├── PLAN.md                      # Дорожная карта и план реализации
└── README.md
```

---

## ⚡ Быстрый запуск

### 1. Запуск Go бэкенда
```bash
cd backend
# Локальный запуск (по умолчанию использует встроенную SQLite БД)
go run cmd/api/main.go
```
Сервер запустится на `http://localhost:8080`.
База данных автоматически создастся и наполнится начальными данными:
- **Менеджер:** логин `admin@a.a`, пароль `admin123`
- **Категории дефектов:** Водоснабжение, Электрика, Окна и двери, Вентиляция и др.
- **Тестовые объекты:** ЖК "Северное Сияние", ЖК "Лазурный Берег", БЦ "Горизонт".

*(Опционально)* Для работы с PostgreSQL:
```bash
docker compose up -d postgres
# В backend/.env установите:
# DATABASE_URL=postgres://postgres:postgrespassword@localhost:5432/defect_management?sslmode=disable
```

### 2. Запуск React фронтенда
```bash
cd ../frontend
npm install
npm start
```
Приложение откроется на `http://localhost:3000`.
Фронтенд уже преднастроен на обращение к локальному Go API (`http://localhost:8080/api/`).
