package http

import (
	"net/http"

	"github.com/iisarieva/url-shortener/internal/usecase"
	"github.com/labstack/echo/v4"
)

// Handler — структура HTTP-хендлера
type Handler struct {
	usecase *usecase.URLUseCase
}

// NewHandler — конструктор хендлера
func NewHandler(u *usecase.URLUseCase) *Handler {
	return &Handler{usecase: u}
}

// RegisterRoutes — регистрация маршрутов
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.POST("/shorten", h.ShortenURL)
	e.GET("/:short", h.Redirect)
	e.DELETE("/:short", h.Delete)

}

// Структура запроса
type ShortenRequest struct {
	OriginalURL string `json:"original_url"`
}

// Структура ответа
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

// ShortenURL — обработчик POST /shorten
// @Summary Создание короткой ссылки
// @Description Принимает оригинальный URL и возвращает сокращённый
// @Tags shorten
// @Accept json
// @Produce json
// @Param data body ShortenRequest true "Оригинальный URL"
// @Success 200 {object} ShortenResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shorten [post]
func (h *Handler) ShortenURL(c echo.Context) error {
	var req ShortenRequest

	// Парсим JSON из тела запроса
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "невалидный JSON"})
	}

	// Проверяем, что ссылка не пустая
	if req.OriginalURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "original_url обязателен"})
	}

	// Генерация и сохранение короткой ссылки
	shortCode, err := h.usecase.CreateShortURL(c.Request().Context(), req.OriginalURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ошибка при создании короткой ссылки"})
	}

	// Базовый адрес можно вынести в конфиг
	baseURL := "http://localhost:8080/"
	return c.JSON(http.StatusOK, ShortenResponse{
		ShortURL: baseURL + shortCode,
	})
}

// Redirect — обработчик GET /:short
// @Summary Редирект по короткой ссылке
// @Description Перенаправляет на оригинальный URL
// @Tags redirect
// @Produce plain
// @Param short path string true "Короткий код"
// @Success 301 {string} string "Redirect"
// @Failure 404 {object} map[string]string
// @Router /{short} [get]
func (h *Handler) Redirect(c echo.Context) error {
	shortCode := c.Param("short")
	if shortCode == "" {
		return c.JSON(400, map[string]string{"error": "short code обязателен"})
	}

	originalURL, err := h.usecase.GetOriginalURL(c.Request().Context(), shortCode)
	if err != nil {
		return c.JSON(404, map[string]string{"error": "ссылка не найдена"})
	}

	// 301 Moved Permanently
	return c.Redirect(301, originalURL)
}

// Delete — обработчик DELETE /:short
// @Summary Удаление короткой ссылки
// @Description Удаляет ссылку по short-коду
// @Tags delete
// @Param short path string true "Короткий код"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /{short} [delete]
func (h *Handler) Delete(c echo.Context) error {
	shortCode := c.Param("short")
	if shortCode == "" {
		return c.JSON(400, map[string]string{"error": "short code обязателен"})
	}

	err := h.usecase.DeleteShortURL(c.Request().Context(), shortCode)
	if err != nil {
		return c.JSON(404, map[string]string{"error": "ссылка не найдена или уже удалена"})
	}

	return c.NoContent(204) // Успешное удаление
}
