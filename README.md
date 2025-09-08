# Доска задач с PostgreSQL и GORM

Система управления задачами для команды с тремя колонками: "Актуальные задачи", "В работе" и "Выполнено".

## Технологии

- **Backend**: Go + Fiber + GORM
- **Database**: PostgreSQL
- **Frontend**: HTML + CSS + JavaScript (Vanilla)
- **Auth**: JWT с HTTP-only cookies

## Настройка PostgreSQL

### 1. Установка PostgreSQL

**macOS (с Homebrew):**
```bash
brew install postgresql
brew services start postgresql
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

### 2. Создание базы данных

```bash
# Подключение к PostgreSQL
sudo -u postgres psql

# Создание базы данных и пользователя
CREATE DATABASE taskboard;
CREATE USER taskboard_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE taskboard TO taskboard_user;
\q
```

### 3. Настройка переменных окружения

Создайте файл `.env` в корне проекта:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=taskboard_user
DB_PASSWORD=secure_password
DB_NAME=taskboard
DB_SSLMODE=disable
```

**Для production использовайте более безопасные настройки:**
```env
DB_HOST=your-postgres-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-secure-password
DB_NAME=taskboard
DB_SSLMODE=require
```

## Запуск проекта

### 1. Установка зависимостей Go

```bash
go mod tidy
```

### 2. Запуск сервера

```bash
go run main.go
```

Сервер будет доступен по адресу: http://localhost:3000

## Возможности системы

### ✅ Управление досками
- Создание доски с уникальным паролем
- Вход в доску по ID и паролю
- Единый пароль для всей команды

### ✅ Три предустановленные колонки
- **"Актуальные задачи"** (todo)
- **"В работе"** (in-progress)  
- **"Выполнено"** (done)

### ✅ Управление карточками
- Создание карточек с заголовком, описанием и ответственным
- Редактирование карточек
- Удаление карточек
- Перетаскивание между колонками (drag & drop)

### ✅ Безопасность
- JWT токены в HTTP-only cookies
- Хеширование паролей с bcrypt
- Защищенные API endpoints

## API Endpoints

### Публичные маршруты
- `POST /api/boards` - создание доски
- `POST /api/boards/:id/login` - вход в доску

### Защищенные маршруты (требуют авторизации)
- `GET /api/board` - получение данных доски
- `POST /api/cards` - создание карточки
- `PUT /api/cards/:id` - редактирование карточки  
- `PUT /api/cards/:id/move` - перемещение карточки
- `DELETE /api/cards/:id` - удаление карточки
- `POST /api/logout` - выход

## Структура базы данных

### Таблица `boards`
- `id` - уникальный идентификатор доски
- `name` - название доски
- `password_hash` - хеш пароля
- `created_at`, `updated_at` - временные метки

### Таблица `columns` (предустановленные)
- `id` - идентификатор колонки (todo, in-progress, done)
- `name` - название колонки
- `order_num` - порядок отображения

### Таблица `cards`
- `id` - уникальный идентификатор карточки
- `board_id` - ссылка на доску
- `title` - заголовок карточки
- `description` - описание (опционально)
- `assignee` - ответственный (опционально)
- `column_id` - ссылка на колонку
- `order_num` - порядок в колонке
- `created_at`, `updated_at` - временные метки

## Особенности GORM интеграции

- **Автомиграция**: таблицы создаются автоматически при запуске
- **Связи**: настроены foreign key с cascade delete
- **Транзакции**: используются для обеспечения целостности данных
- **Индексы**: автоматическое создание индексов для производительности

## Разработка

### Структура проекта
```
task-board/
├── database/          # Подключение к БД и миграции
├── frontend/          # HTML, CSS, JS файлы
├── handlers/          # HTTP обработчики
├── middleware/        # JWT аутентификация
├── models/           # GORM модели данных
├── services/         # Бизнес-логика
├── go.mod           # Go зависимости
└── main.go          # Точка входа
```

### Добавление новых функций

1. Обновите модели в `models/models.go`
2. Добавьте методы в `services/board_service.go`
3. Создайте обработчики в `handlers/board_handler.go`
4. Добавьте маршруты в `main.go`

## Мониторинг и логи

Приложение использует встроенное логирование Fiber для отслеживания запросов и ошибок.

## Безопасность в production

1. Используйте сильные пароли для БД
2. Настройте SSL/TLS для PostgreSQL
3. Используйте HTTPS для веб-сервера
4. Регулярно обновляйте зависимости
5. Настройте firewall и ограничьте доступ к БД
