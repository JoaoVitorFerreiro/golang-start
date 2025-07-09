package main

import (
	"context"
	"log"

	"github.com/JoaoVitorFerreiro/golang-start/config"
	"github.com/JoaoVitorFerreiro/golang-start/internal/application/service"
	"github.com/JoaoVitorFerreiro/golang-start/internal/infra/http"
	"github.com/JoaoVitorFerreiro/golang-start/internal/infra/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// IMPORTANTE: Import dos docs gerados
	_ "github.com/JoaoVitorFerreiro/golang-start/docs"
)

// @title           User API
// @version         1.0
// @description     API para gerenciamento de usuários usando Domain-Driven Design
// @termsOfService  http://swagger.io/terms/

// @contact.name   João Vitor
// @contact.url    http://www.swagger.io/support
// @contact.email  jvferreiro1@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

func main() {
    cfg := config.Load()
    
    gin.SetMode(cfg.GinMode)
    
    var userRepo repository.UserRepository
    
    if cfg.Env == "production" || cfg.Env == "staging" {
        pool, err := setupDatabase(cfg)
        if err != nil {
            log.Fatal("Failed to connect to database:", err)
        }
        defer pool.Close()
        
        userRepo = repository.NewUserRepository(repository.Postgres, pool)
        log.Println("Using PostgreSQL repository")
    } else {
        userRepo = repository.NewUserRepository(repository.InMemory, nil)
        log.Println("Using in-memory repository")
    }
    
    userService := service.NewUserService(userRepo)    
    userHandler := http.NewUserHandler(userService)
    
    router := setupRouter()
    userHandler.RegisterRoutes(router)
    
    log.Printf("Server starting on port %s (env: %s)", cfg.Port, cfg.Env)
    log.Printf("Swagger docs available at: http://localhost:%s/swagger/index.html", cfg.Port)
    if err := router.Run(":" + cfg.Port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

func setupDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(cfg.DatabaseURL)
    if err != nil {
        return nil, err
    }
    
    config.MaxConns = cfg.DBMaxConns
    config.MinConns = cfg.DBMinConns
    config.MaxConnLifetime = cfg.DBMaxConnLifetime
    config.MaxConnIdleTime = cfg.DBMaxConnIdleTime
    
    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return nil, err
    }
    
    if err := pool.Ping(context.Background()); err != nil {
        return nil, err
    }
    
    return pool, nil
}

func setupRouter() *gin.Engine {
    router := gin.Default()
    
    router.Use(gin.Recovery())
    
    // CORS para Swagger
    router.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })

    // Swagger endpoint - configuração correta
    router.GET("/swagger/*any", ginSwagger.WrapHandler(
        swaggerFiles.Handler,
        ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
        ginSwagger.DeepLinking(true),
        ginSwagger.DocExpansion("none"),
    ))

    router.GET("/health", healthCheck)
    
    return router
}

// healthCheck godoc
// @Summary      Health Check
// @Description  Verifica se a API está funcionando
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func healthCheck(c *gin.Context) {
    c.JSON(200, gin.H{
        "status":  "ok",
        "service": "user-api",
    })
}