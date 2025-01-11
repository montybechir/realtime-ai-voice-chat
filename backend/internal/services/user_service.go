package services

import (
	"context"
	"errors"
	"strconv"
)

var (
	ErrNotFound     = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input")

	mockUsers = []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
		{ID: 3, Name: "Bob Wilson", Email: "bob@example.com"},
		{ID: 4, Name: "Alice Brown", Email: "alice@example.com"},
	}
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserService struct{}

func NewUserService(dbUrl string) *UserService {
	return &UserService{}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return mockUsers, nil
	}
}

func (s *UserService) CreateUser(ctx context.Context, name, email string) (User, error) {
	if name == "" || email == "" {
		return User{}, ErrInvalidInput
	}

	newUser := User{
		ID:    len(mockUsers) + 1,
		Name:  name,
		Email: email,
	}

	return newUser, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (User, error) {
	if id == "" {
		return User{}, ErrInvalidInput
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		return User{}, ErrInvalidInput
	}

	for _, user := range mockUsers {
		if user.ID == userID {
			return user, nil
		}
	}

	return User{}, ErrNotFound
}
