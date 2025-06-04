package service

import (
	"testing"

	"github.com/strazhnikovt/TestShop/internal/domain"
	"github.com/strazhnikovt/TestShop/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Register(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("success", func(t *testing.T) {
		user := &domain.User{
			FirstName: "John",
			LastName:  "Doe",
			Login:     "john@example.com",
			Password:  "password123",
			Age:       25,
		}

		mockRepo.On("GetByLogin", "john@example.com").Return((*domain.User)(nil), nil)
		mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Run(func(args mock.Arguments) {
			u := args.Get(0).(*domain.User)
			assert.Equal(t, "John Doe", u.FullName)
			assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("password123")))
		}).Return(nil)

		err := service.Register(user)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("duplicate email", func(t *testing.T) {
		user := &domain.User{
			Login: "existing@example.com",
		}

		existingUser := &domain.User{}
		mockRepo.On("GetByLogin", "existing@example.com").Return(existingUser, nil)

		err := service.Register(user)
		assert.Error(t, err)
		assert.Equal(t, "user already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Login(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("success", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
		user := &domain.User{
			ID:       1,
			Login:    "user@example.com",
			Password: string(hashedPassword),
		}

		mockRepo.On("GetByLogin", "user@example.com").Return(user, nil)

		authUser, err := service.Login("user@example.com", "correctpassword")
		assert.NoError(t, err)
		assert.Equal(t, 1, authUser.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
		user := &domain.User{
			Login:    "user@example.com",
			Password: string(hashedPassword),
		}

		mockRepo.On("GetByLogin", "user@example.com").Return(user, nil)

		_, err := service.Login("user@example.com", "wrongpassword")
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("GetByLogin", "nonexistent@example.com").Return((*domain.User)(nil), nil)

		_, err := service.Login("nonexistent@example.com", "password")
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
