package service

import (
	"testing"

	"github.com/JoaoVitorFerreiro/golang-start/internal/application/dto"
	"github.com/JoaoVitorFerreiro/golang-start/internal/infra/repository"
)

func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    repo := repository.NewUserRepository(repository.InMemory, nil)
    service := NewUserService(repo)
    
    req := dto.CreateUserRequest{
        Name:  "João Silva",
        Email: "joao@email.com",
    }
    
    // Act
    user, err := service.CreateUser(req)
    
    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if user.Name != req.Name {
        t.Errorf("Expected name %s, got %s", req.Name, user.Name)
    }
    
    if user.Email != req.Email {
        t.Errorf("Expected email %s, got %s", req.Email, user.Email)
    }
    
    if user.ID == "" {
        t.Error("Expected ID to be generated")
    }
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
    // Arrange
    repo := repository.NewUserRepository(repository.InMemory, nil)
    service := NewUserService(repo)
    
    req := dto.CreateUserRequest{
        Name:  "João Silva",
        Email: "joao@email.com",
    }
    
    // Criar primeiro usuário
    _, err := service.CreateUser(req)
    if err != nil {
        t.Fatalf("Failed to create first user: %v", err)
    }
    
    // Act - tentar criar usuário com mesmo email
    _, err = service.CreateUser(req)
    
    // Assert
    if err == nil {
        t.Error("Expected error for duplicate email")
    }
    
    if err.Error() != "email already exists" {
        t.Errorf("Expected 'email already exists', got %s", err.Error())
    }
}

func TestUserService_GetUserByID(t *testing.T) {
    // Arrange
    repo := repository.NewUserRepository(repository.InMemory, nil)
    service := NewUserService(repo)
    
    // Criar usuário
    createReq := dto.CreateUserRequest{
        Name:  "João Silva",
        Email: "joao@email.com",
    }
    createdUser, _ := service.CreateUser(createReq)
    
    // Act
    user, err := service.GetUserByID(createdUser.ID)
    
    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if user.ID != createdUser.ID {
        t.Errorf("Expected ID %s, got %s", createdUser.ID, user.ID)
    }
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
    // Arrange
    repo := repository.NewUserRepository(repository.InMemory, nil)
    service := NewUserService(repo)
    
    // Act
    _, err := service.GetUserByID("non-existent-id")
    
    // Assert
    if err == nil {
        t.Error("Expected error for non-existent user")
    }
    
    if err.Error() != "user not found" {
        t.Errorf("Expected 'user not found', got %s", err.Error())
    }
}

func TestUserService_UpdateUser(t *testing.T) {
    // Arrange
    repo := repository.NewUserRepository(repository.InMemory, nil)
    service := NewUserService(repo)
    
    // Criar usuário
    createReq := dto.CreateUserRequest{
        Name:  "João Silva",
        Email: "joao@email.com",
    }
    createdUser, _ := service.CreateUser(createReq)
    
    // Act
    updateReq := dto.UpdateUserRequest{
        Name: "João Santos",
    }
    updatedUser, err := service.UpdateUser(createdUser.ID, updateReq)
    
    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if updatedUser.Name != "João Santos" {
        t.Errorf("Expected name 'João Santos', got %s", updatedUser.Name)
    }
    
    if updatedUser.Email != createdUser.Email {
        t.Error("Email should not have changed")
    }
}

func TestUserService_DeleteUser(t *testing.T) {
    // Arrange
    repo := repository.NewUserRepository(repository.InMemory, nil)
    service := NewUserService(repo)
    
    // Criar usuário
    createReq := dto.CreateUserRequest{
        Name:  "João Silva",
        Email: "joao@email.com",
    }
    createdUser, _ := service.CreateUser(createReq)
    
    // Act
    err := service.DeleteUser(createdUser.ID)
    
    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    // Verificar se foi deletado
    _, err = service.GetUserByID(createdUser.ID)
    if err == nil {
        t.Error("Expected error when getting deleted user")
    }
}