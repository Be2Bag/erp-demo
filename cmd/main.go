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
	"github.com/Be2Bag/erp-demo/cron"
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
	departmentRepo := repositories.NewDepartmentRepository(database)
	positionRepo := repositories.NewPositionRepository(database)
	kpiEvaluationRepo := repositories.NewKPIEvaluationRepository(database)
	categoryRepo := repositories.NewCategoryRepository(database)
	signTypeRepo := repositories.NewSignTypeRepository(database)
	bankAccountsRepo := repositories.NewBankAccountsRepository(database)
	transactionCategoryRepo := repositories.NewTransactionCategoryRepository(database)
	inComeRepo := repositories.NewInComeRepository(database)
	expenseRepo := repositories.NewExpenseRepository(database)
	payableRepo := repositories.NewPayableRepository(database)
	receivableRepo := repositories.NewReceivableRepository(database)
	receiptRepo := repositories.NewReceiptRepository(database)

	userSvc := services.NewUserService(*cfg, userRepo, dropDownRepo, cloudflareStorage, taskRepo)
	upLoadSvc := services.NewUpLoadService(*cfg, authRepo, upLoadRepo, userRepo, cloudflareStorage)
	adminSvc := services.NewAdminService(*cfg, adminRepo, authRepo, userRepo)
	dropDownSvc := services.NewDropDownService(*cfg, dropDownRepo)
	kpiSvc := services.NewKPIService(*cfg, kpiRepo, userRepo)
	taskSvc := services.NewTaskService(*cfg, taskRepo, userRepo, workFlowRepo, departmentRepo, kpiEvaluationRepo, kpiRepo, signJobRepo)
	authSvc := services.NewAuthService(*cfg, authRepo, userRepo)
	workFlowSvc := services.NewWorkflowService(*cfg, workFlowRepo)
	signJobSvc := services.NewSignJobService(*cfg, signJobRepo, dropDownRepo, taskRepo, inComeRepo, receivableRepo)
	projectSvc := services.NewProjectService(*cfg, projectRepo, userRepo, signJobRepo, taskRepo)
	departmentSvc := services.NewDepartmentService(*cfg, departmentRepo, userRepo)
	positionSvc := services.NewPositionService(*cfg, positionRepo, departmentRepo, userRepo)
	kpiEvaluationSvc := services.NewKPIEvaluationService(*cfg, kpiRepo, userRepo, kpiEvaluationRepo, taskRepo, departmentRepo, projectRepo, signJobRepo)
	categorySvc := services.NewCategoryService(*cfg, categoryRepo)
	signTypeSvc := services.NewSignTypeService(*cfg, signTypeRepo, userRepo)
	bankAccountsSvc := services.NewBankAccountService(*cfg, bankAccountsRepo)
	transactionCategorySvc := services.NewTransactionCategoryService(*cfg, transactionCategoryRepo)
	inComeSvc := services.NewInComeService(*cfg, inComeRepo, transactionCategoryRepo)
	expenseSvc := services.NewExpenseService(*cfg, expenseRepo, transactionCategoryRepo)
	payableSvc := services.NewPayablesService(*cfg, payableRepo, bankAccountsRepo)
	receivableSvc := services.NewReceivableService(*cfg, receivableRepo, bankAccountsRepo, signJobRepo, inComeRepo)
	receiptSvc := services.NewReceiptService(*cfg, receiptRepo, bankAccountsRepo)

	// เริ่มต้น Cronjob สำหรับตรวจสอบสถานะ Payable และ Receivable
	statusChecker := cron.NewStatusChecker(payableRepo, receivableRepo)
	if err := statusChecker.Start(); err != nil {
		log.Printf("เริ่ม cronjob ไม่สำเร็จ: %v", err)
	} else {
		log.Println("🚀🚀🚀 Cronjob เริ่มทำงานแล้ว - ระบบจะตรวจสอบสถานะ Payable/Receivable ทุกวันเวลา 00:00 น.")
	}

	userHdl := handlers.NewUserHandler(userSvc, upLoadSvc, authCookieMiddleware)
	upLoadHdl := handlers.NewUpLoadHandler(upLoadSvc, authCookieMiddleware)
	adminHdl := handlers.NewAdminHandler(adminSvc, authCookieMiddleware)
	dropDownHdl := handlers.NewDropDownHandler(dropDownSvc, authCookieMiddleware)
	authHdl := handlers.NewAuthHandler(authSvc)
	kpiHdl := handlers.NewKPIHandler(kpiSvc, authCookieMiddleware)
	taskHdl := handlers.NewTaskHandler(taskSvc, authCookieMiddleware)
	workFlowHdl := handlers.NewWorkFlowHandler(workFlowSvc, authCookieMiddleware)
	signJobHdl := handlers.NewSignJobHandler(signJobSvc, authCookieMiddleware)
	projectHdl := handlers.NewProjectHandler(projectSvc, authCookieMiddleware)
	departmentHdl := handlers.NewDepartmentHandler(departmentSvc, authCookieMiddleware)
	positionHdl := handlers.NewPositionHandler(positionSvc, authCookieMiddleware)
	kpiEvaluationHdl := handlers.NewKPIEvaluationHandler(kpiEvaluationSvc, authCookieMiddleware)
	categoryHdl := handlers.NewCategoryHandler(categorySvc, authCookieMiddleware)
	signTypeHdl := handlers.NewSignTypeHandler(signTypeSvc, authCookieMiddleware)
	bankAccountsHdl := handlers.NewBankAccountsHandler(bankAccountsSvc, authCookieMiddleware)
	transactionCategoryHdl := handlers.NewTransactionCategoryHandler(transactionCategorySvc, authCookieMiddleware)
	inComeHdl := handlers.NewInComeHandler(inComeSvc, authCookieMiddleware)
	expenseHdl := handlers.NewExpenseHandler(expenseSvc, authCookieMiddleware)
	payableHdl := handlers.NewPayableHandler(payableSvc, authCookieMiddleware)
	receivableHdl := handlers.NewReceivableHandler(receivableSvc, authCookieMiddleware)
	receiptHdl := handlers.NewReceiptHandler(receiptSvc, authCookieMiddleware)
	cronHdl := handlers.NewCronHandler(statusChecker, authCookieMiddleware)

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
	departmentHdl.DepartmentRoutes(apiGroup)
	positionHdl.PositionRoutes(apiGroup)
	kpiEvaluationHdl.KPIEvaluationRoutes(apiGroup)
	categoryHdl.CategoryRoutes(apiGroup)
	signTypeHdl.SignTypeRoutes(apiGroup)
	bankAccountsHdl.BankAccountsRoutes(apiGroup)
	transactionCategoryHdl.TransactionCategoryRoutes(apiGroup)
	inComeHdl.InComeRoutes(apiGroup)
	expenseHdl.ExpenseRoutes(apiGroup)
	payableHdl.PayableRoutes(apiGroup)
	receivableHdl.ReceivableRoutes(apiGroup)
	receiptHdl.ReceiptRoutes(apiGroup)
	cronHdl.CronRoutes(apiGroup)

	app.Use("/swagger", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": cfg.Swagger.Key,
		},
	}))
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.Fatal(app.Listen(":3000"))
}
