package main

import (
	"fmt"

	"github.com/Furkanberkay/ticket-booking-project-v1/config"
	"github.com/Furkanberkay/ticket-booking-project-v1/db"
	"github.com/Furkanberkay/ticket-booking-project-v1/handlers"
	"github.com/Furkanberkay/ticket-booking-project-v1/repositories"
	"github.com/Furkanberkay/ticket-booking-project-v1/services"
	"github.com/gofiber/fiber/v2"
)

func main() {

	envConfig := config.NewEnvConfig()
	db := db.Init(envConfig, db.DBMigrator)

	eventRepository := repositories.NewEventRepository(db)

	eventService := services.NewEventService(eventRepository)

	app := fiber.New()
	api := app.Group("/api")
	events := api.Group("/events")
	handlers.NewEventHandler(events, eventService)

	app.Listen(fmt.Sprintf(":") + envConfig.ServerPort)
}
