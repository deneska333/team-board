package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

type Claims struct {
	BoardID string `json:"board_id"`
	jwt.RegisteredClaims
}

// GenerateToken создает JWT токен для доступа к доске
func GenerateToken(boardID string) (string, error) {
	claims := Claims{
		BoardID: boardID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// AuthMiddleware проверяет JWT токен из HTTP-only cookie
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем токен из HTTP-only cookie
		tokenString := c.Cookies("auth_token")
		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Требуется авторизация",
			})
		}

		// Парсим и валидируем токен
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"error": "Недействительный токен",
			})
		}

		// Сохраняем board_id в контексте
		c.Locals("board_id", claims.BoardID)
		return c.Next()
	}
}
