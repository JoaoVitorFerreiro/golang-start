package dto

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" example:"João Silva"`
	Email string `json:"email" binding:"required,email" example:"joao@email.com"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" example:"João Santos"`
	Email string `json:"email,omitempty" example:"joao.santos@email.com"`
}

type UserResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string `json:"name" example:"João Silva"`
	Email     string `json:"email" example:"joao@email.com"`
	CreatedAt string `json:"created_at" example:"2024-07-08T10:30:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-07-08T11:45:00Z"`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"user not found"`
	Message string `json:"message,omitempty" example:"O usuário com ID especificado não foi encontrado"`
}