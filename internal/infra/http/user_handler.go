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

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
    userGroup := router.Group("/users")
    {
        userGroup.POST("", h.CreateUser)
        userGroup.GET("", h.GetAllUsers)
        userGroup.GET("/:id", h.GetUserByID)
        userGroup.PUT("/:id", h.UpdateUser)
        userGroup.DELETE("/:id", h.DeleteUser)
    }
}

// CreateUser godoc
// @Summary      Criar usuário
// @Description  Cria um novo usuário no sistema
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.CreateUserRequest  true  "Dados do usuário"
// @Success      201   {object}  dto.UserResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      409   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "invalid request",
            Message: err.Error(),
        })
        return
    }
    
    user, err := h.userService.CreateUser(req)
    if err != nil {
        if err.Error() == "email already exists" {
            c.JSON(http.StatusConflict, dto.ErrorResponse{
                Error:   "email already exists",
                Message: "Um usuário com este email já existe",
            })
            return
        }
        
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Error:   "internal server error",
            Message: "Erro interno do servidor",
        })
        return
    }
    
    c.JSON(http.StatusCreated, user)
}

// GetAllUsers godoc
// @Summary      Listar usuários
// @Description  Retorna todos os usuários cadastrados
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   dto.UserResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
    users, err := h.userService.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Error:   "internal server error",
            Message: "Erro ao buscar usuários",
        })
        return
    }
    
    c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary      Buscar usuário por ID
// @Description  Retorna um usuário específico pelo ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID do usuário"
// @Success      200  {object}  dto.UserResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
    id := c.Param("id")
    
    user, err := h.userService.GetUserByID(id)
    if err != nil {
        if err.Error() == "user not found" {
            c.JSON(http.StatusNotFound, dto.ErrorResponse{
                Error:   "user not found",
                Message: "Usuário não encontrado",
            })
            return
        }
        
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Error:   "internal server error",
            Message: "Erro interno do servidor",
        })
        return
    }
    
    c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary      Atualizar usuário
// @Description  Atualiza os dados de um usuário existente
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string                 true  "ID do usuário"
// @Param        user  body      dto.UpdateUserRequest  true  "Dados para atualização"
// @Success      200   {object}  dto.UserResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      404   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
    id := c.Param("id")
    var req dto.UpdateUserRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "invalid request",
            Message: err.Error(),
        })
        return
    }
    
    user, err := h.userService.UpdateUser(id, req)
    if err != nil {
        if err.Error() == "user not found" {
            c.JSON(http.StatusNotFound, dto.ErrorResponse{
                Error:   "user not found",
                Message: "Usuário não encontrado",
            })
            return
        }
        
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Error:   "internal server error",
            Message: "Erro interno do servidor",
        })
        return
    }
    
    c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary      Deletar usuário
// @Description  Remove um usuário do sistema
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "ID do usuário"
// @Success      204
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
    id := c.Param("id")
    
    err := h.userService.DeleteUser(id)
    if err != nil {
        if err.Error() == "user not found" {
            c.JSON(http.StatusNotFound, dto.ErrorResponse{
                Error:   "user not found",
                Message: "Usuário não encontrado",
            })
            return
        }
        
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Error:   "internal server error",
            Message: "Erro interno do servidor",
        })
        return
    }
    
    c.Status(http.StatusNoContent)
}