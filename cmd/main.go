package main

// @title        ERP Demo API
// @version      1.0
// @description  This is an ERP API demo.
// @host         https://erp-demo-9ux8.onrender.com
// @BasePath     /service/api

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/Be2Bag/erp-demo/config"
	_ "github.com/Be2Bag/erp-demo/docs"
	"github.com/Be2Bag/erp-demo/handler"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/pkg/db"
	"github.com/Be2Bag/erp-demo/repository"
	"github.com/Be2Bag/erp-demo/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("config:", err)
	}
	client, err := db.Connect(cfg.Mongo.URI)
	if err != nil {
		log.Fatal("mongo connect:", err)
	}
	database := db.GetDB(client, cfg.Mongo.Database)

	userRepo := repository.NewUserRepository(database)
	userSvc := service.NewUserService(*cfg, userRepo)
	userHdl := handler.NewUserHandler(userSvc)

	app := fiber.New()

	app.Use(middleware.TimeoutMiddleware(5 * time.Second))

	apiGroup := app.Group("/service/api")
	userHdl.RegisterRoutes(apiGroup)

	app.Use("/swagger", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "123456",
		},
	}))
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.Fatal(app.Listen(":3000"))
}
