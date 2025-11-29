package main

import (
	"log"
	"patungan-server/internal/config"
	"patungan-server/internal/handler"
	"patungan-server/internal/middleware"
	"patungan-server/internal/models"
	"patungan-server/internal/repository"
	"patungan-server/internal/service"
	"patungan-server/pkg/database"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	db, err := database.NewPostgresDB(cfg.Database.DSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	log.Println("Running database migrations...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Bill{},
		&models.BillParticipant{},
		&models.BillItem{},
		&models.ItemParticipant{},
		&models.Payment{},
	); err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	billRepo := repository.NewBillRepository(db)
	participantRepo := repository.NewParticipantRepository(db)
	itemRepo := repository.NewItemRepository(db)
	itemParticipantRepo := repository.NewItemParticipantRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)
	billService := service.NewBillService(billRepo)
	participantService := service.NewParticipantService(participantRepo, billRepo, userRepo, itemRepo)
	itemService := service.NewItemService(itemRepo, billRepo, itemParticipantRepo, participantService)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	billHandler := handler.NewBillHandler(billService)
	participantHandler := handler.NewParticipantHandler(participantService)
	itemHandler := handler.NewItemHandler(itemService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))


	// Routes
	api := app.Group("/api/v1")

	// Health check 
	api.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Patungan API is running",
		})
	})

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Bill routes (protected)
	authMiddleware := middleware.AuthMiddleware(cfg)
	bills := api.Group("/bills", authMiddleware)
	bills.Get("/", billHandler.GetAll)
	bills.Post("/", billHandler.Create)
	bills.Get("/:id", billHandler.GetByID)
	bills.Put("/:id", billHandler.Update)
	bills.Delete("/:id", billHandler.Delete)

	// Participant routes (protected)
	bills.Post("/:billId/participants", participantHandler.AddParticipants)
	bills.Get("/:billId/participants", participantHandler.GetParticipants)
	bills.Delete("/:billId/participants/:participantId", participantHandler.RemoveParticipant)
	bills.Put("/participants/:participantId", participantHandler.UpdateParticipant)
	bills.Post("/:billId/recalculate", participantHandler.RecalculateSplits)

	// Item routes (protected)
	bills.Post("/:billId/items", itemHandler.AddItems)
	bills.Get("/:billId/items", itemHandler.GetItems)
	bills.Put("/items/:itemId", itemHandler.UpdateItem)
	bills.Delete("/items/:itemId", itemHandler.DeleteItem)

	// Start server
	log.Printf("Server starting on port %s (Environment: %s)", cfg.Server.Port, cfg.Server.Env)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}