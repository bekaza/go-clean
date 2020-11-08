package usecase

import (
	"testing"
)

func TestAuthFlow(t *testing.T) {
	// repo := new(mock_domain.MockUserRepository)

	// uc := NewUserService(repo, "salt", []byte("secret"), 86400, time.Duration(2)*time.Second)

	// var (
	// 	username = "user"
	// 	password = "pass"

	// 	ctx = context.Background()

	// 	user = &domain.User{
	// 		Username: username,
	// 		Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
	// 	}
	// )

	// // Sign Up
	// repo.On("CreateUser", user).Return(nil)
	// err := uc.Login(ctx, username, password)
	// assert.NoError(t, err)

	// // Sign In (Get Auth Token)
	// repo.On("GetUser", user.Username, user.Password).Return(user, nil)
	// token, err := uc.SignIn(ctx, username, password)
	// assert.NoError(t, err)
	// assert.NotEmpty(t, token)

	// // Verify token
	// parsedUser, err := uc.ParseToken(ctx, token)
	// assert.NoError(t, err)
	// assert.Equal(t, user, parsedUser)
}
