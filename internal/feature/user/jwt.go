package user

import (
	"errors"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Claims struct {
	ID   uuid.UUID  `json:"id"`
	Role model.Role `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJwtToken(id uuid.UUID, role model.Role, expire time.Time) (string, error) {
	claim := &Claims{
		ID:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := tokens.SignedString(secretKey)

	if err != nil {
		slog.Error("error generate jwt token", err.Error())
		return "", err
	}
	return tokenString, nil
}

func AddCookieTokens(id uuid.UUID, Role model.Role, w http.ResponseWriter) error {
	expirationTimeAccess := time.Now().Add(15 * time.Minute)
	expirationTimeRefresh := time.Now().Add(14 * 24 * time.Hour)
	refreshToken, err := GenerateJwtToken(id, Role, expirationTimeRefresh)
	if err != nil {
		return err
	}
	accessToken, err := GenerateJwtToken(id, Role, expirationTimeAccess)
	if err != nil {
		return err
	}

	slog.Warn("tokens generated", slog.String("access_token", accessToken), slog.String("refresh_token", refreshToken))

	http.SetCookie(w, GenerateCookie("accessToken", expirationTimeAccess, false, accessToken))
	http.SetCookie(w, GenerateCookie("refreshToken", expirationTimeRefresh, true, refreshToken))

	return nil
}
func ParseToken(tokenString string) (*Claims, error) {
	// Парсинг токена с использованием функции обратного вызова для получения секретного ключа
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка используемого метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Возвращаем секретный ключ
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	// Проверка валидности токена и получение claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
func NewAccessToken(tokenString string, threshold time.Duration, w http.ResponseWriter) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// Проверка оставшегося времени жизни токена
	remainingTime := time.Until(claims.ExpiresAt.Time)
	if remainingTime > threshold {
		return tokenString, nil
	}

	// Генерация нового токена с обновленным временем истечения
	newExpire := time.Now().Add(1 * time.Hour) // Задайте желаемое время жизни нового токена
	newToken, err := GenerateJwtToken(claims.ID, claims.Role, newExpire)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, GenerateCookie("accessToken", newExpire, false, newToken))
	return newToken, nil
}

func GenerateCookie(name string, expire time.Time, httpOnly bool, value string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expire,
		Path:     "/",
		HttpOnly: httpOnly,
	}

	return cookie
}

func DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, GenerateCookie("refreshToken", time.Unix(0, 0), true, ""))
}
