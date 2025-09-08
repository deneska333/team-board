package handlers

import (
	"time"

	"task-board/middleware"
	"task-board/models"
	"task-board/services"

	"github.com/gofiber/fiber/v2"
)

type BoardHandler struct {
	boardService *services.BoardService
}

func NewBoardHandler(boardService *services.BoardService) *BoardHandler {
	return &BoardHandler{
		boardService: boardService,
	}
}

// CreateBoard создает новую доску
func (h *BoardHandler) CreateBoard(c *fiber.Ctx) error {
	var req models.CreateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Неверный формат запроса",
		})
	}

	if req.Name == "" || req.Password == "" {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Имя доски и пароль обязательны",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Пароль должен содержать минимум 6 символ����в",
		})
	}

	board, err := h.boardService.CreateBoard(req.Name, req.Password)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Ошибка создания доски",
		})
	}

	// Генерируем JWT токен
	token, err := middleware.GenerateToken(board.ID)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Ошибка генерации токена",
		})
	}

	// Устанавливаем HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false, // Установите true в production с HTTPS
		SameSite: "Lax",
	})

	return c.JSON(board)
}

// Login вход в доску по паролю
func (h *BoardHandler) Login(c *fiber.Ctx) error {
	boardID := c.Params("id")

	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Неверный формат запроса",
		})
	}

	if err := h.boardService.ValidatePassword(boardID, req.Password); err != nil {
		return c.Status(401).JSON(models.ErrorResponse{
			Error: "Неверный пароль",
		})
	}

	// Генерируем JWT токен
	token, err := middleware.GenerateToken(boardID)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Ошибка генерации токена",
		})
	}

	// Устанавливаем HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	return c.JSON(models.LoginResponse{
		Message: "Успешный вход",
		BoardID: boardID,
	})
}

// GetBoard получает данные доски
func (h *BoardHandler) GetBoard(c *fiber.Ctx) error {
	boardID := c.Locals("board_id").(string)

	board, err := h.boardService.GetBoard(boardID)
	if err != nil {
		return c.Status(404).JSON(models.ErrorResponse{
			Error: "Доска не найдена",
		})
	}

	return c.JSON(board)
}

// CreateCard создает новую карточку
func (h *BoardHandler) CreateCard(c *fiber.Ctx) error {
	boardID := c.Locals("board_id").(string)

	var req models.CreateCardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Неверный формат запроса",
		})
	}

	if req.Title == "" || req.ColumnID == "" {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Заголовок и ID колонки обязательны",
		})
	}

	card, err := h.boardService.CreateCard(boardID, req)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(201).JSON(card)
}

// UpdateCard обновляет карточку
func (h *BoardHandler) UpdateCard(c *fiber.Ctx) error {
	boardID := c.Locals("board_id").(string)
	cardID := c.Params("cardId")

	var req models.UpdateCardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Неверный формат запроса",
		})
	}

	card, err := h.boardService.UpdateCard(boardID, cardID, req)
	if err != nil {
		return c.Status(404).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(card)
}

// MoveCard перемещает карточку между колонками
func (h *BoardHandler) MoveCard(c *fiber.Ctx) error {
	boardID := c.Locals("board_id").(string)
	cardID := c.Params("cardId")

	var req models.MoveCardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Неверный формат запроса",
		})
	}

	if req.ColumnID == "" {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "ID колонки обязателен",
		})
	}

	card, err := h.boardService.MoveCard(boardID, cardID, req)
	if err != nil {
		return c.Status(404).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(card)
}

// DeleteCard удаляет карточку
func (h *BoardHandler) DeleteCard(c *fiber.Ctx) error {
	boardID := c.Locals("board_id").(string)
	cardID := c.Params("cardId")

	err := h.boardService.DeleteCard(boardID, cardID)
	if err != nil {
		return c.Status(404).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Карточка удалена",
	})
}

// Logout выход из системы
func (h *BoardHandler) Logout(c *fiber.Ctx) error {
	// Удаляем cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"message": "Выход выполнен",
	})
}

// CreateColumn создает новую колонку
func (h *BoardHandler) CreateColumn(c *fiber.Ctx) error {
	boardID := c.Locals("board_id").(string)

	var req models.CreateColumnRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Неверный формат запроса",
		})
	}

	if req.Name == "" {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Название колонки обязательно",
		})
	}

	column, err := h.boardService.CreateColumn(boardID, req)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(201).JSON(column)
}

// UpdateColumn обновляет колонку
func (h *BoardHandler) UpdateColumn(c *fiber.Ctx) error {
	boardID := c.Locals("board_id").(string)
	columnID := c.Params("columnId")

	var req models.UpdateColumnRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Неверный формат запроса",
		})
	}

	if req.Name == "" {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Название колонки обязательно",
		})
	}

	column, err := h.boardService.UpdateColumn(boardID, columnID, req)
	if err != nil {
		return c.Status(404).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(column)
}

// DeleteColumn удаляет колонку
func (h *BoardHandler) DeleteColumn(c *fiber.Ctx) error {
	boardID := c.Locals("board_id").(string)
	columnID := c.Params("columnId")

	err := h.boardService.DeleteColumn(boardID, columnID)
	if err != nil {
		return c.Status(404).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Колонка удалена",
	})
}
