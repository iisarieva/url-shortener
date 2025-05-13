package http

import (
	"net/http"

	"github.com/iisarieva/url-shortener/configs"
	"github.com/iisarieva/url-shortener/internal/usecase"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase *usecase.URLUseCase
}

func NewHandler(u *usecase.URLUseCase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.POST("/shorten", h.ShortenURL)
	e.GET("/:short", h.Redirect)
	e.DELETE("/:short", h.Delete)
}

type ShortenRequest struct {
	OriginalURL string `json:"original_url"`
}

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

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("invalid JSON"))
	}

	if req.OriginalURL == "" {
		return c.JSON(http.StatusBadRequest, errorResponse("original_url is required"))
	}

	shortCode, err := h.usecase.CreateShortURL(c.Request().Context(), req.OriginalURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("failed to create short URL"))
	}

	return c.JSON(http.StatusOK, ShortenResponse{
		ShortURL: configs.BaseURL + shortCode,
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
		return c.JSON(http.StatusBadRequest, errorResponse("short code is required"))
	}

	originalURL, err := h.usecase.GetOriginalURL(c.Request().Context(), shortCode)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResponse("URL not found"))
	}

	return c.Redirect(http.StatusMovedPermanently, originalURL)
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
		return c.JSON(http.StatusBadRequest, errorResponse("short code is required"))
	}

	err := h.usecase.DeleteShortURL(c.Request().Context(), shortCode)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResponse("URL not found or already deleted"))
	}

	return c.NoContent(http.StatusNoContent)
}

func errorResponse(msg string) map[string]string {
	return map[string]string{"error": msg}
}
