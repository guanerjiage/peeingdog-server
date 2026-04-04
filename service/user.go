package service

import (
	"context"
	"fmt"
	"time"

	"peeingdog-server/sql/queries/generated"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserService struct {
	queries *generated.Queries
}

func NewUserService(q *generated.Queries) *UserService {
	return &UserService{queries: q}
}

// GetAllUsers retrieves all users from the database
func (s *UserService) GetAllUsers(ctx context.Context) ([]User, error) {
	users, err := s.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	result := make([]User, len(users))
	for i, u := range users {
		result[i] = User{
			ID:        int(u.ID),
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Time,
		}
	}
	return result, nil
}

// GetUserByID retrieves a single user by ID
func (s *UserService) GetUserByID(ctx context.Context, id int) (*User, error) {
	user, err := s.queries.GetUserByID(ctx, int32(id))
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &User{
		ID:        int(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}

// CreateUser creates a new user in the database
func (s *UserService) CreateUser(ctx context.Context, name, email string) (*User, error) {
	if name == "" || email == "" {
		return nil, fmt.Errorf("name and email are required")
	}

	user, err := s.queries.CreateUser(ctx, generated.CreateUserParams{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &User{
		ID:        int(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}
