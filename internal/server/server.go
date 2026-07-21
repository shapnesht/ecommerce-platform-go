package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/shapnesht/ecommerce-platform-go/internal/config"
	"github.com/shapnesht/ecommerce-platform-go/internal/services"
	"gorm.io/gorm"
)

type Server struct {
	config         *config.Config
	db             *gorm.DB
	logger         *zerolog.Logger
	authService    *services.AuthService
	userService    *services.UserService
	productService *services.ProductService
	uploadService  *services.UploadService
	cartService    *services.CartService
	orderService   *services.OrderService
}

func New(config *config.Config, db *gorm.DB, logger *zerolog.Logger, authService *services.AuthService, userService *services.UserService, productService *services.ProductService, uploadService *services.UploadService, cartService *services.CartService, orderService *services.OrderService) *Server {
	return &Server{
		config:         config,
		db:             db,
		logger:         logger,
		authService:    authService,
		userService:    userService,
		productService: productService,
		uploadService:  uploadService,
		cartService:    cartService,
		orderService:   orderService,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	router.GET("/health", s.healthCheck)

	router.Static("/uploads", "./uploads")

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{ //nolint:gocritic // Need this for readability
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("/refresh", s.refreshToken)
			auth.POST("/logout", s.logout)
		}

		protected := router.Group("/")
		protected.Use(s.authMiddleware())
		{
			users := protected.Group("/users")
			{
				userRoutes := users
				userRoutes.GET("/profile", s.getProfile)
				userRoutes.PUT("/profile", s.updateProfile)
			}

			categories := protected.Group("/categories")
			{
				categoryRoutes := categories
				categoryRoutes.POST("/", s.adminMiddleware(), s.createCategory)
				categoryRoutes.PUT("/:id", s.adminMiddleware(), s.updateCategory)
				categoryRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)
			}

			products := protected.Group("/products")
			{
				productRoutes := products
				productRoutes.POST("/", s.adminMiddleware(), s.createProduct)
				productRoutes.PUT("/", s.adminMiddleware(), s.updateProduct)
				productRoutes.DELETE("/", s.adminMiddleware(), s.deleteProduct)
				productRoutes.POST("/:id/images", s.adminMiddleware(), s.uploadImage)
			}

		}

		api.GET("/products", s.getProducts)
		api.GET("/categories", s.getCategories)
		api.GET("/products/:id", s.getProduct)
	}

	return router
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
