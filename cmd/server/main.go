package main

import (
	"context"
	"github.com/iisarieva/url-shortener/internal/delivery"
	"github.com/labstack/echo/v4"
	goRedis "github.com/redis/go-redis/v9"
	"os"
	"time"

	rds "github.com/iisarieva/url-shortener/internal/repository/redis"
	"github.com/iisarieva/url-shortener/internal/usecase"

	_ "github.com/iisarieva/url-shortener/docs"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	e := echo.New()

	e.Use(loggingMiddleware)

	redisClient := createRedisClient()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal().Err(err).Msg("❌ Не удалось подключиться к Redis")
	}

	urlRepo := rds.NewRedisRepository(redisClient)
	urlUsecase := usecase.NewURLUseCase(urlRepo)
	handler := delivery.NewHandler(urlUsecase)

	handler.RegisterRoutes(e)

	e.GET("/docs/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func createRedisClient() *goRedis.Client {
	addr := os.Getenv("REDIS_HOST")
	if addr == "" {
		addr = "localhost:6379"
	}
	log.Info().Str("addr", addr).Msg("🔌 Подключение к Redis")

	return goRedis.NewClient(&goRedis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

}

func loggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		start := time.Now()

		err := next(c)

		log.Info().
			Str("method", req.Method).
			Str("path", req.URL.Path).
			Int("status", res.Status).
			Dur("duration", time.Since(start)).
			Msg("🌐 HTTP-запрос")

		return err
	}
}
