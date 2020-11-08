package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bekaza/go-clean/domain"
	"github.com/dgrijalva/jwt-go/v4"
	"golang.org/x/crypto/bcrypt"
)

// AuthClaims ...
type AuthClaims struct {
	jwt.StandardClaims
	User *domain.User `json:"user"`
}

type userServiceImp struct {
	userRepo       domain.UserRepository
	signingKey     []byte        // for hash access token and refresh token
	expireDuration time.Duration // expired token
	contextTimeout time.Duration
}

// NewUserService ...
func NewUserService(
	ur domain.UserRepository,
	signingKey []byte,
	tokenTTLSeconds time.Duration,
	timeout time.Duration) domain.UserService {
	return &userServiceImp{
		ur,
		signingKey,
		time.Second * tokenTTLSeconds,
		timeout,
	}
}

func (usi *userServiceImp) Register(ctx context.Context, username string, password string) error {
	ctx, cancel := context.WithTimeout(ctx, usi.contextTimeout)
	defer cancel()
	hashPassword := hashPassword(password)
	return usi.userRepo.CreateUser(ctx, username, hashPassword)
}

func (usi *userServiceImp) Login(c context.Context, username string, password string) (string, string, error) {
	ctx, cancel := context.WithTimeout(c, usi.contextTimeout)
	defer cancel()

	user, err := usi.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	if !user.IsActive {
		return "", "", errors.New("User not active")
	}

	isValid := checkPassword(user.Password, password)
	if !isValid {
		return "", "", errors.New("Password is incorrect")
	}

	claims := AuthClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(usi.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, _ := token.SignedString(usi.signingKey)
	return accessToken, "", nil
}

func hashPassword(password string) (hashedPassword string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func checkPassword(hashedPassword string, password string) (isPasswordValid bool) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}

func (usi *userServiceImp) ParseToken(ctx context.Context, accessToken string) (user *domain.User, err error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return usi.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims.User, nil
	}
	return user, err
}
