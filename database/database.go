package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"task-board/models"
)

var DB *gorm.DB

// Config содержит параметры подключения к БД
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetConfigFromEnv получает конфигурацию из переменных окружения
func GetConfigFromEnv() Config {
	return Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "taskboard"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// Connect подключается к базе данных PostgreSQL через GORM
func Connect() error {
	config := GetConfigFromEnv()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Получаем базовое соединение для проверки
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("ошибка получения базового соединения: %w", err)
	}

	// Проверяем соединение
	if err = sqlDB.Ping(); err != nil {
		return fmt.Errorf("ошибка проверки соединения с БД: %w", err)
	}

	log.Println("Успешно подключились к PostgreSQL через GORM")
	return nil
}

// Close закрывает соединение с базой данных
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// MigrateModels выполняет автомиграцию моделей
func MigrateModels() error {
	err := DB.AutoMigrate(&models.Board{}, &models.Column{}, &models.Card{})
	if err != nil {
		return fmt.Errorf("ошибка миграции: %w", err)
	}

	log.Println("Миграция моделей выполнена успешно")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
