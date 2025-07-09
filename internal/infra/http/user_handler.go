package http

import (
	"net/http"

	"github.com/JoaoVitorFerreiro/golang-start/internal/application/dto"
	"github.com/JoaoVitorFerreiro/golang-start/internal/application/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

// POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
            "details": err.Error(),
        })
        return
    }

    user, err := h.userService.CreateUser(req)
    if err != nil {
        if err.Error() == "email already exists" {
            c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, user)
}

// GET /users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
    id := c.Param("id")
    
    user, err := h.userService.GetUserByID(id)
    if err != nil {
        if err.Error() == "user not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}

// GET /users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
    users, err := h.userService.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "users": users,
        "total": len(users),
    })
}

// PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
    id := c.Param("id")
    
    var req dto.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
            "details": err.Error(),
        })
        return
    }

    user, err := h.userService.UpdateUser(id, req)
    if err != nil {
        if err.Error() == "user not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        if err.Error() == "email already exists" {
            c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}

// DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
    id := c.Param("id")
    
    err := h.userService.DeleteUser(id)
    if err != nil {
        if err.Error() == "user not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusNoContent, nil)
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
    userRoutes := router.Group("/users")
    {
        userRoutes.POST("", h.CreateUser)
        userRoutes.GET("/:id", h.GetUserByID)
        userRoutes.GET("", h.GetAllUsers)
        userRoutes.PUT("/:id", h.UpdateUser)
        userRoutes.DELETE("/:id", h.DeleteUser)
    }
}