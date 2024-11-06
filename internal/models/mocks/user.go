package mocks

import (
	"time"

	"github.com/tneuqole/snippetbox/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "user@example.com" && password == "password" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) UpdatePassword(id int, password string) error {
	return nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return &models.User{
			ID:      1,
			Name:    "John Smith",
			Email:   "user@example.com",
			Created: time.Now(),
		}, nil
	default:
		return nil, models.ErrNoRecord
	}
}
