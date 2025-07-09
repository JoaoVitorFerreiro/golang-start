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
)

func main() {
    // Carregar configurações
    cfg := config.Load()
    
    // Configurar Gin
    gin.SetMode(cfg.GinMode)
    
    // Configurar repositório baseado no ambiente
    var userRepo repository.UserRepository
    
    if cfg.Env == "production" || cfg.Env == "staging" {
        // Usar PostgreSQL
        pool, err := setupDatabase(cfg)
        if err != nil {
            log.Fatal("Failed to connect to database:", err)
        }
        defer pool.Close()
        
        userRepo = repository.NewUserRepository(repository.Postgres, pool)
        log.Println("Using PostgreSQL repository")
    } else {
        // Usar memória em desenvolvimento
        userRepo = repository.NewUserRepository(repository.InMemory, nil)
        log.Println("Using in-memory repository")
    }
    
    // Configurar serviços
    userService := service.NewUserService(userRepo)
    
    // Configurar handlers
    userHandler := http.NewUserHandler(userService)
    
    // Configurar router
    router := setupRouter()
    userHandler.RegisterRoutes(router)
    
    // Iniciar servidor
    log.Printf("Server starting on port %s (env: %s)", cfg.Port, cfg.Env)
    if err := router.Run(":" + cfg.Port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

func setupDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(cfg.DatabaseURL)
    if err != nil {
        return nil, err
    }
    
    // Configurações de performance do .env
    config.MaxConns = cfg.DBMaxConns
    config.MinConns = cfg.DBMinConns
    config.MaxConnLifetime = cfg.DBMaxConnLifetime
    config.MaxConnIdleTime = cfg.DBMaxConnIdleTime
    
    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return nil, err
    }
    
    // Testar conexão
    if err := pool.Ping(context.Background()); err != nil {
        return nil, err
    }
    
    return pool, nil
}

func setupRouter() *gin.Engine {
    router := gin.Default()
    
    // Middleware básico
    router.Use(gin.Recovery())
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "service": "user-api",
        })
    })
    
    return router
}