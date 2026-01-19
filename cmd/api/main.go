package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Furkanberkay/ticket-booking-project-v1/config"
	"github.com/Furkanberkay/ticket-booking-project-v1/db"
	"github.com/Furkanberkay/ticket-booking-project-v1/handlers"
	"github.com/Furkanberkay/ticket-booking-project-v1/middleware"
	"github.com/Furkanberkay/ticket-booking-project-v1/repositories"
	"github.com/Furkanberkay/ticket-booking-project-v1/services"
	"github.com/Furkanberkay/ticket-booking-project-v1/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lmittmann/tint"
)

func main() {
	logLevel := slog.LevelInfo
	if os.Getenv("APP_ENV") == "dev" || os.Getenv("ENV") == "dev" {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      logLevel,
		TimeFormat: "15:04:05",
	}))
	slog.SetDefault(logger)

	envConfig := config.NewEnvConfig()
	dbConn := db.Init(envConfig, db.DBMigrator)

	redisClient := db.NewRedisClient()
	defer redisClient.Close()

	jwtWrapper := utils.NewJWTWrapper(envConfig.JWTSecret)

	eventRepository := repositories.NewEventRepository(dbConn)
	ticketRepository := repositories.NewTicketRepository(dbConn)

	userRepository := repositories.NewUserRepository(dbConn, logger)
	tokenRepository := repositories.NewRedisTokenRepository(redisClient)

	eventService := services.NewEventService(eventRepository, logger)
	ticketService := services.NewTicketService(ticketRepository, eventRepository, logger)

	authService := services.NewAuthService(userRepository, tokenRepository, jwtWrapper, logger)

	app := fiber.New(fiber.Config{
		AppName: "Ticket Booking API v1",
	})

	api := app.Group("/api")
	validate := validator.New()

	authGroup := api.Group("/auth")
	handlers.NewAuthHandler(validate, authGroup, authService)

	eventsGroup := api.Group("/events")
	handlers.NewEventHandler(validate, eventsGroup, eventService)

	ticketsGroup := api.Group("/tickets")

	ticketsGroup.Use(middleware.AuthProtected(jwtWrapper))

	handlers.NewTicketHandler(validate, ticketsGroup, ticketService)

	logger.Info("Server starting...", "port", envConfig.ServerPort)

	if err := app.Listen(fmt.Sprintf(":%s", envConfig.ServerPort)); err != nil {
		logger.Error("Server failed to start", "error", err)
	}
}
