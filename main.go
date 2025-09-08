package main

import (
	"log"

	"task-board/database"
	"task-board/handlers"
	"task-board/middleware"
	"task-board/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем переменные окружения системы")
	}

	// Подключаемся к базе данных
	if err := database.Connect(); err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer database.Close()

	// Выполняем автомиграцию моделей
	if err := database.MigrateModels(); err != nil {
		log.Fatal("Ошибка миграции моделей:", err)
	}

	// Создаем экземпляр Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Сервисы
	boardService := services.NewBoardService()
	boardHandler := handlers.NewBoardHandler(boardService)

	// Публичные маршруты
	api := app.Group("/api")

	// Создание доски
	api.Post("/boards", boardHandler.CreateBoard)

	// Вход в доску
	api.Post("/boards/:id/login", boardHandler.Login)

	// Защищенные маршруты (требуют аутентификации)
	protected := api.Use(middleware.AuthMiddleware())

	// Получение данных доски
	protected.Get("/board", boardHandler.GetBoard)

	// Работа с колонками
	protected.Post("/columns", boardHandler.CreateColumn)
	protected.Put("/columns/:columnId", boardHandler.UpdateColumn)
	protected.Delete("/columns/:columnId", boardHandler.DeleteColumn)

	// Работа с карточками
	protected.Post("/cards", boardHandler.CreateCard)
	protected.Put("/cards/:cardId", boardHandler.UpdateCard)
	protected.Put("/cards/:cardId/move", boardHandler.MoveCard)
	protected.Delete("/cards/:cardId", boardHandler.DeleteCard)

	// Выход
	protected.Post("/logout", boardHandler.Logout)

	// Статические файлы (фронтенд)
	app.Static("/", "./frontend")

	// Запуск сервера
	log.Println("Сервер запущен на порту :3000")
	log.Fatal(app.Listen(":3000"))
}
