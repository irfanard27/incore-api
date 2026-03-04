package route

import (
	"irfanard27/incore-api/internal/delivery/http"
	"irfanard27/incore-api/internal/delivery/middleware"
	"irfanard27/incore-api/internal/infra/jwt"
	"irfanard27/incore-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Route struct {
	jwtMiddleware *middleware.JWTMiddleware

	authHandler      *http.AuthHandler
	inventoryHandler *http.InventoryHandler
	stockinHandler   *http.StockInHandler
}

type RouteCfg struct {
	JWTService       jwt.JWTService
	AuthUsecase      usecase.AuthUsecase
	InventoryUsecase usecase.InventoryUsecase
	StockInUsecase   usecase.StockInUsecase
}

func NewRoute(cfg RouteCfg) *Route {
	authHandler := http.NewAuthHandler(cfg.AuthUsecase)
	inventoryHandler := http.NewInventoryHandler(cfg.InventoryUsecase)
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWTService)
	stockinHandler := http.NewStockInHandler(cfg.StockInUsecase)

	return &Route{
		authHandler:      authHandler,
		inventoryHandler: inventoryHandler,
		jwtMiddleware:    jwtMiddleware,
		stockinHandler:   stockinHandler,
	}
}

func (r *Route) Setup(router *gin.Engine) {
	api := router.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/logout", r.authHandler.Logout)
	}

	// private routes

	inventtory := api.Group("/inventories")
	inventtory.Use(r.jwtMiddleware.RequireAuth())
	{
		inventtory.POST("", r.inventoryHandler.CreateInventory)
		inventtory.GET("", r.inventoryHandler.GetAllInventories)
		inventtory.GET("/:id", r.inventoryHandler.GetInventoryById)
		inventtory.PUT("/:id", r.inventoryHandler.UpdateInventory)
		inventtory.DELETE("/:id", r.inventoryHandler.DeleteInventory)
		inventtory.GET("/search", r.inventoryHandler.SearchInventory)
	}

	stockin := api.Group("/stocks-in")
	stockin.Use(r.jwtMiddleware.RequireAuth())
	{
		stockin.POST("", r.stockinHandler.CreateStockIn)
	}

}
