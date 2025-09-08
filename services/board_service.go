package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"task-board/database"
	"task-board/models"
)

type BoardService struct {
	db *gorm.DB
}

func NewBoardService() *BoardService {
	return &BoardService{
		db: database.DB,
	}
}

// CreateBoard создает новую доску с тремя колонками по умолчанию
func (s *BoardService) CreateBoard(name, password string) (*models.Board, error) {
	id := generateID()

	// Хешируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("ошибка хеширования пароля")
	}

	// Создаем доску
	board := &models.Board{
		ID:           id,
		Name:         name,
		PasswordHash: string(passwordHash),
	}

	// Используем транзакцию для создания доски и колонок
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Сохраняем доску
		if err := tx.Create(board).Error; err != nil {
			return err
		}

		// Создаем три колонки по умолчанию
		defaultColumns := []models.Column{
			{ID: generateID(), BoardID: id, Name: "Актуальные задачи", OrderNum: 1},
			{ID: generateID(), BoardID: id, Name: "В работе", OrderNum: 2},
			{ID: generateID(), BoardID: id, Name: "Выполнено", OrderNum: 3},
		}

		for _, col := range defaultColumns {
			if err := tx.Create(&col).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.New("ошибка создания доски")
	}

	// Загружаем доску с колонками для ответа
	return s.GetBoard(id)
}

// GetBoard получает доску по ID с колонками и карточками
func (s *BoardService) GetBoard(id string) (*models.Board, error) {
	var board models.Board

	// Получаем доску с колонками и карточками
	if err := s.db.Preload("Columns.Cards", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_num ASC")
	}).Preload("Columns", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_num ASC")
	}).First(&board, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("доска не найдена")
		}
		return nil, errors.New("ошибка получения доски")
	}

	return &board, nil
}

// CreateColumn создает новую колонку
func (s *BoardService) CreateColumn(boardID string, req models.CreateColumnRequest) (*models.Column, error) {
	// Проверяем, что доска существует
	var board models.Board
	if err := s.db.First(&board, "id = ?", boardID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("доска не найдена")
		}
		return nil, errors.New("ошибка получения доски")
	}

	// Получаем следующий порядковый номер
	var maxOrder int64
	s.db.Model(&models.Column{}).Where("board_id = ?", boardID).
		Select("COALESCE(MAX(order_num), 0)").Scan(&maxOrder)

	column := &models.Column{
		ID:       generateID(),
		BoardID:  boardID,
		Name:     req.Name,
		OrderNum: int(maxOrder) + 1,
	}

	if err := s.db.Create(column).Error; err != nil {
		return nil, errors.New("ошибка создания колонки")
	}

	return column, nil
}

// UpdateColumn обновляет колонку
func (s *BoardService) UpdateColumn(boardID, columnID string, req models.UpdateColumnRequest) (*models.Column, error) {
	var column models.Column

	if err := s.db.Where("id = ? AND board_id = ?", columnID, boardID).First(&column).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("колонка не найдена")
		}
		return nil, errors.New("ошибка получения колонки")
	}

	if err := s.db.Model(&column).Update("name", req.Name).Error; err != nil {
		return nil, errors.New("ошибка обновления колонки")
	}

	return &column, nil
}

// DeleteColumn удаляет колонку (и все её карточки)
func (s *BoardService) DeleteColumn(boardID, columnID string) error {
	result := s.db.Where("id = ? AND board_id = ?", columnID, boardID).Delete(&models.Column{})

	if result.Error != nil {
		return errors.New("ошибка удаления колонки")
	}

	if result.RowsAffected == 0 {
		return errors.New("колонка не найдена")
	}

	return nil
}

// ValidatePassword проверяет пароль ��оски
func (s *BoardService) ValidatePassword(boardID, password string) error {
	var board models.Board

	if err := s.db.Select("password_hash").First(&board, "id = ?", boardID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("доска не найдена")
		}
		return errors.New("ошибка проверки пароля")
	}

	// Сравниваем хеш пароля
	if err := bcrypt.CompareHashAndPassword([]byte(board.PasswordHash), []byte(password)); err != nil {
		return errors.New("неверный пароль")
	}

	return nil
}

// CreateCard создает новую карточку в указанной колонке
func (s *BoardService) CreateCard(boardID string, req models.CreateCardRequest) (*models.Card, error) {
	// Проверяем, что доска существует
	var board models.Board
	if err := s.db.First(&board, "id = ?", boardID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("доска не найдена")
		}
		return nil, errors.New("ошибка получения доски")
	}

	// Получаем сле��ующий порядковый номер для колонки
	var maxOrder int64
	s.db.Model(&models.Card{}).Where("board_id = ? AND column_id = ?", boardID, req.ColumnID).
		Select("COALESCE(MAX(order_num), 0)").Scan(&maxOrder)

	card := &models.Card{
		ID:          generateID(),
		BoardID:     boardID,
		Title:       req.Title,
		Description: req.Description,
		Assignee:    req.Assignee,
		ColumnID:    req.ColumnID,
		OrderNum:    int(maxOrder) + 1,
	}

	// Сохраняем карточку в БД
	if err := s.db.Create(card).Error; err != nil {
		return nil, errors.New("ошибка создания карточки")
	}

	return card, nil
}

// UpdateCard обновляет карточку
func (s *BoardService) UpdateCard(boardID, cardID string, req models.UpdateCardRequest) (*models.Card, error) {
	var card models.Card

	// Проверяем, что карточка существует и принадлежлежи�� доске
	if err := s.db.Where("id = ? AND board_id = ?", cardID, boardID).First(&card).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("карточка не найдена")
		}
		return nil, errors.New("ошибка получения карточки")
	}

	// Подготавливаем данные для обновления
	updates := make(map[string]interface{})

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Assignee != "" {
		updates["assignee"] = req.Assignee
	}

	// Если есть что обновлять
	if len(updates) > 0 {
		if err := s.db.Model(&card).Updates(updates).Error; err != nil {
			return nil, errors.New("ошибка обновления карточки")
		}
	}

	// Перезагружаем карточку с обновленными данными
	if err := s.db.First(&card, "id = ?", cardID).Error; err != nil {
		return nil, errors.New("ошибка перезагрузки карточки")
	}

	return &card, nil
}

// MoveCard перемещает карточку между колонками
func (s *BoardService) MoveCard(boardID, cardID string, req models.MoveCardRequest) (*models.Card, error) {
	var card models.Card

	// Начинаем транзакцию
	return &card, s.db.Transaction(func(tx *gorm.DB) error {
		// Получаем текущую карточку
		if err := tx.Where("id = ? AND board_id = ?", cardID, boardID).First(&card).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("карточка не найдена")
			}
			return errors.New("ошибка получения карточки")
		}

		// Если колонка не изменилась, просто возвращаем карточку
		if card.ColumnID == req.ColumnID {
			return nil
		}

		// Получаем следующий порядковый номер для целевой колонки
		var maxOrder int64
		tx.Model(&models.Card{}).Where("board_id = ? AND column_id = ?", boardID, req.ColumnID).
			Select("COALESCE(MAX(order_num), 0)").Scan(&maxOrder)

		// Обновляем карточку
		updates := map[string]interface{}{
			"column_id": req.ColumnID,
			"order_num": int(maxOrder) + 1,
		}

		if err := tx.Model(&card).Updates(updates).Error; err != nil {
			return errors.New("ошибка перемещения карточки")
		}

		// Перезагружаем карточку
		if err := tx.First(&card, "id = ?", cardID).Error; err != nil {
			return errors.New("ошибка перезагрузки карточки")
		}

		return nil
	})
}

// DeleteCard удаляет карточку
func (s *BoardService) DeleteCard(boardID, cardID string) error {
	result := s.db.Where("id = ? AND board_id = ?", cardID, boardID).Delete(&models.Card{})

	if result.Error != nil {
		return errors.New("ошибка удаления карточки")
	}

	if result.RowsAffected == 0 {
		return errors.New("карточка не найдена")
	}

	return nil
}

func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
