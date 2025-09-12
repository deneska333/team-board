package models

import (
	"time"
)

// Board представляет доску задач
type Board struct {
	ID           string    `json:"id" gorm:"primaryKey;size:32"`
	Name         string    `json:"name" gorm:"not null;size:255"`
	PasswordHash string    `json:"-" gorm:"not null;size:255"`
	CreatedAt    time.Time `json:"created" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated" gorm:"autoUpdateTime"`

	// Связь с колонками
	Columns []Column `json:"columns" gorm:"foreignKey:BoardID"`
}

// Column представляет колонку на доске (теперь динамические)
type Column struct {
	ID        string    `json:"id" gorm:"primaryKey;size:32"`
	BoardID   string    `json:"board_id" gorm:"not null;size:32;index"`
	Name      string    `json:"name" gorm:"not null;size:100"`
	OrderNum  int       `json:"order" gorm:"not null"`
	CreatedAt time.Time `json:"created" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated" gorm:"autoUpdateTime"`

	// Связи
	Board Board  `json:"-" gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE"`
	Cards []Card `json:"cards" gorm:"foreignKey:ColumnID"`
}

// Card представляет карточку задачи
type Card struct {
	ID          string     `json:"id" gorm:"primaryKey;size:32"`
	BoardID     string     `json:"board_id" gorm:"not null;size:32;index"`
	Title       string     `json:"title" gorm:"not null;size:500"`
	Description string     `json:"description" gorm:"type:text"`
	Assignee    string     `json:"assignee" gorm:"size:255"`
	Deadline    *time.Time `json:"deadline" gorm:"type:timestamp"`
	ColumnID    string     `json:"column_id" gorm:"not null;size:32;index"`
	OrderNum    int        `json:"order" gorm:"not null;default:1"`
	CreatedAt   time.Time  `json:"created" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated" gorm:"autoUpdateTime"`

	// Связи
	Board  Board  `json:"-" gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE"`
	Column Column `json:"-" gorm:"foreignKey:ColumnID;constraint:OnDelete:CASCADE"`
}

// TableName указывает имя таблицы для модели Card
func (Card) TableName() string {
	return "cards"
}

// TableName указывает имя таблицы для модели Board
func (Board) TableName() string {
	return "boards"
}

// TableName указывает имя таблицы для модели Column
func (Column) TableName() string {
	return "columns"
}

// Запросы для API
type CreateBoardRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Password string `json:"password" validate:"required"`
}

type CreateCardRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description"`
	Assignee    string     `json:"assignee"`
	Deadline    *time.Time `json:"deadline"`
	ColumnID    string     `json:"column_id" validate:"required"`
}

type UpdateCardRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Assignee    string     `json:"assignee"`
	Deadline    *time.Time `json:"deadline"`
}

type MoveCardRequest struct {
	ColumnID string `json:"column_id" validate:"required"`
	Order    int    `json:"order"`
}

// Запросы для управления колонками
type CreateColumnRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateColumnRequest struct {
	Name string `json:"name" validate:"required"`
}

type MoveColumnRequest struct {
	Order int `json:"order" validate:"required"`
}

// Ответы API
type LoginResponse struct {
	Message string `json:"message"`
	BoardID string `json:"board_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Предустановленные колонки
var DefaultColumns = []Column{
	{
		ID:       "todo",
		Name:     "Актуальные задачи",
		OrderNum: 1,
		Cards:    []Card{},
	},
	{
		ID:       "in-progress",
		Name:     "В работе",
		OrderNum: 2,
		Cards:    []Card{},
	},
	{
		ID:       "done",
		Name:     "Выполн��но",
		OrderNum: 3,
		Cards:    []Card{},
	},
}
