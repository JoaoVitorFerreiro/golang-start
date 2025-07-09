package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(name, email string) (*User, error){
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}

	now := time.Now()

	return &User{
		ID: 			uuid.New().String(),
		Name: name,
		Email: email,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) UpdateName(name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) UpdateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	u.Email = email
	u.UpdatedAt = time.Now()
	return nil
}
