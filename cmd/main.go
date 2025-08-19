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
	handlers "github.com/Be2Bag/erp-demo/handlers"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/pkg/db"
	"github.com/Be2Bag/erp-demo/pkg/storage"
	repositories "github.com/Be2Bag/erp-demo/repositories"
	services "github.com/Be2Bag/erp-demo/services"
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

	authCookieMiddleware := middleware.NewMiddleware(cfg.JWT)

	userRepo := repositories.NewUserRepository(database)
	authRepo := repositories.NewAuthRepository(database)
	dropDownRepo := repositories.NewDropDownRepository(database)
	adminRepo := repositories.NewAdminRepository(database)
	upLoadRepo := repositories.NewUpLoadRepository(database)
	kpiRepo := repositories.NewKPIRepository(database)
	taskRepo := repositories.NewTaskRepository(database)
	workFlowRepo := repositories.NewWorkFlowRepository(database)
	signJobRepo := repositories.NewSignJobRepository(database)
	projectRepo := repositories.NewProjectRepository(database)

	userSvc := services.NewUserService(*cfg, userRepo, dropDownRepo, supabaseStorage, cloudflareStorage)
	upLoadSvc := services.NewUpLoadService(*cfg, authRepo, upLoadRepo, supabaseStorage, userRepo, cloudflareStorage)
	adminSvc := services.NewAdminService(*cfg, adminRepo, authRepo, userRepo)
	dropDownSvc := services.NewDropDownService(*cfg, dropDownRepo)
	kpiSvc := services.NewKPIService(*cfg, kpiRepo, userRepo)
	taskSvc := services.NewTaskService(*cfg, taskRepo, userRepo)
	authSvc := services.NewAuthService(*cfg, authRepo, userRepo)
	workFlowSvc := services.NewWorkflowService(*cfg, workFlowRepo)
	signJobSvc := services.NewSignJobService(*cfg, signJobRepo)
	projectSvc := services.NewProjectService(*cfg, projectRepo)

	userHdl := handlers.NewUserHandler(userSvc, upLoadSvc, authCookieMiddleware)
	upLoadHdl := handlers.NewUpLoadHandler(upLoadSvc, authCookieMiddleware)
	adminHdl := handlers.NewAdminHandler(adminSvc)
	dropDownHdl := handlers.NewDropDownHandler(dropDownSvc)
	authHdl := handlers.NewAuthHandler(authSvc)
	kpiHdl := handlers.NewKPIHandler(kpiSvc, authCookieMiddleware)
	taskHdl := handlers.NewTaskHandler(taskSvc, authCookieMiddleware)
	workFlowHdl := handlers.NewWorkFlowHandler(workFlowSvc, authCookieMiddleware)
	signJobHdl := handlers.NewSignJobHandler(signJobSvc, authCookieMiddleware)
	projectHdl := handlers.NewProjectHandler(projectSvc, authCookieMiddleware)

	app := fiber.New()

	// Enable CORS for your frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,https://erp-demo-frontend.vercel.app,https://dev-erp-demo-frontend.vercel.app,https://www.rkp-media.com",
		AllowCredentials: true,
	}))

	app.Use(authCookieMiddleware.TimeoutMiddleware(10 * time.Second))

	apiGroup := app.Group("/service/api")
	userHdl.UserRoutes(apiGroup)
	authHdl.AuthRoutes(apiGroup)
	dropDownHdl.DropDownRoutes(apiGroup)
	adminHdl.AdminRoutes(apiGroup)
	upLoadHdl.UpLoadRoutes(apiGroup)
	kpiHdl.KPIRoutes(apiGroup)
	taskHdl.TaskRoutes(apiGroup)
	workFlowHdl.WorkFlowRoutes(apiGroup)
	signJobHdl.SignJobRoutes(apiGroup)
	projectHdl.ProjectRoutes(apiGroup)

	app.Use("/swagger", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": cfg.Swagger.Key,
		},
	}))
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.Fatal(app.Listen(":3000"))
}
