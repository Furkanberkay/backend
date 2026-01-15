package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Furkanberkay/ticket-booking-project-v1/config"
	"github.com/Furkanberkay/ticket-booking-project-v1/db"
	"github.com/Furkanberkay/ticket-booking-project-v1/handlers"
	"github.com/Furkanberkay/ticket-booking-project-v1/repositories"
	"github.com/Furkanberkay/ticket-booking-project-v1/services"
	"github.com/gofiber/fiber/v2"
	"github.com/lmittmann/tint"
)

func main() {
	logLevel := slog.LevelInfo
	if os.Getenv("APP_ENV") == "dev" || os.Getenv("ENV") == "dev" {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	envConfig := config.NewEnvConfig()
	db := db.Init(envConfig, db.DBMigrator)

	eventRepository := repositories.NewEventRepository(db)
	ticketRepository := repositories.NewTicketRepository(db)

	eventService := services.NewEventService(eventRepository)
	ticketService := services.NewTicketService(ticketRepository, logger)

	app := fiber.New()
	api := app.Group("/api")

	events := api.Group("/events")
	handlers.NewEventHandler(events, eventService)

	tickets := api.Group("/tickets")
	handlers.NewTicketHandler(tickets, ticketService)

	app.Listen(fmt.Sprintf(":") + envConfig.ServerPort)
}
