package main

import (
	"log"
	"net/http"
	"os"

	"irfanard27/incore-api/db"
	"irfanard27/incore-api/internal/delivery/route"
	"irfanard27/incore-api/internal/infra/jwt"
	"irfanard27/incore-api/internal/infra/pgsql"
	"irfanard27/incore-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setup database connection
	database := db.SetupPGSQL()

	log.Println("Database connected successfully!")

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Migrations completed successfully!")
	defer database.Close()

	// Setup repositories
	userRepo := pgsql.NewUserRepository(database)
	inventoryRepo := pgsql.NewInventoryRepository(database)
	stockInRepo := pgsql.NewStockInRepository(database)
	stockInItemRepo := pgsql.NewStockInItemRepository(database)

	// Setup JWT service with secret key from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // Default for development
		log.Println("Warning: Using default JWT secret key. Set JWT_SECRET in production!")
	}
	jwtService := jwt.NewJWTService(jwtSecret)

	// Setup routes
	routeHandler := route.NewRoute(route.RouteCfg{
		AuthUsecase:      usecase.NewAuthUsecase(userRepo, jwtService),
		InventoryUsecase: usecase.NewInventoryUsecase(inventoryRepo),
		JWTService:       jwtService,
		StockInUsecase:   usecase.NewStockInUsecase(stockInRepo, stockInItemRepo),
	})

	// Setup Gin router
	router := gin.Default()

	// Setup routes
	routeHandler.Setup(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
