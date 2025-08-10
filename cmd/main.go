package main

// @title        ERP Demo API
// @version      1.0
// @description  This is an ERP API demo.
// @host         api.rkp-media.com
// @BasePath     /service/api

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/Be2Bag/erp-demo/config"
	_ "github.com/Be2Bag/erp-demo/docs"
	"github.com/Be2Bag/erp-demo/handler"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/pkg/db"
	"github.com/Be2Bag/erp-demo/pkg/storage"
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

	supabaseStorage, err := storage.NewSupabaseStorage(cfg.Supabase)
	if err != nil {
		log.Fatal("supabase storage:", err)
	}

	cloudflareStorage, err := storage.NewCloudflareStorage(cfg.Cloudflare)
	if err != nil {
		log.Fatal("cloudflare storage:", err)
	}

	userRepo := repository.NewUserRepository(database)
	authRepo := repository.NewAuthRepository(database)
	dropDownRepo := repository.NewDropDownRepository(database)
	adminRepo := repository.NewAdminRepository(database)
	upLoadRepo := repository.NewUpLoadRepository(database)
	kpiRepo := repository.NewKPIRepository(database)

	userSvc := service.NewUserService(*cfg, userRepo, dropDownRepo, supabaseStorage, cloudflareStorage)
	upLoadSvc := service.NewUpLoadService(*cfg, authRepo, upLoadRepo, supabaseStorage, userRepo, cloudflareStorage)
	adminSvc := service.NewAdminService(*cfg, adminRepo, authRepo, userRepo)
	dropDownSvc := service.NewDropDownService(*cfg, dropDownRepo)
	kpiSvc := service.NewKPIService(*cfg, kpiRepo, userRepo)

	authSvc := service.NewAuthService(*cfg, authRepo, userRepo)

	userHdl := handler.NewUserHandler(userSvc, upLoadSvc)
	upLoadHdl := handler.NewUpLoadHandler(upLoadSvc)
	adminHdl := handler.NewAdminHandler(adminSvc)
	dropDownHdl := handler.NewDropDownHandler(dropDownSvc)
	authHdl := handler.NewAuthHandler(authSvc)
	kpiHdl := handler.NewKPIHandler(kpiSvc)

	app := fiber.New()

	// Enable CORS for your frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,https://erp-demo-frontend.vercel.app,https://dev-erp-demo-frontend.vercel.app,https://www.rkp-media.com",
		AllowCredentials: true,
	}))

	app.Use(middleware.TimeoutMiddleware(30 * time.Second))

	apiGroup := app.Group("/service/api")
	userHdl.UserRoutes(apiGroup)
	authHdl.AuthRoutes(apiGroup)
	dropDownHdl.DropDownRoutes(apiGroup)
	adminHdl.AdminRoutes(apiGroup)
	upLoadHdl.UpLoadRoutes(apiGroup)
	kpiHdl.KPIRoutes(apiGroup)

	app.Use("/swagger", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": cfg.Swagger.Key,
		},
	}))
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.Fatal(app.Listen(":3000"))
}
