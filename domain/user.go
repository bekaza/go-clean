package domain

import (
	"context"
	"time"
)

//go:generate go run github.com/golang/mock/mockgen --build_flags=--mod=vendor -destination=./mock/user.go -source=./user.go UserService

// CtxUserKey ...
const CtxUserKey = "user"

// User ...
type User struct {
	ID        uint       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Username  string     `gorm:"UNIQUE;TYPE:varchar(10)" json:"username"`
	Password  string     `gorm:"TYPE:varchar(200)" json:"password"`
	IsActive  bool       `gorm:"DEFAULT:true" json:"is_active"`
	CreatedAt *time.Time `gorm:"DEFAULT:now()" json:"created_at"`
	UpdatedAt *time.Time `gorm:"DEFAULT:now()" json:"updated_at"`
	DeletedAt *time.Time `gorm:"DEFAULT:NULL" json:"deleted_at"`
}

// TableName ...
func (User) TableName() string {
	return "users"
}

// UserService ...
type UserService interface {
	Register(ctx context.Context, username string, password string) error
	Login(ctx context.Context, username string, password string) (string, string, error)
	ParseToken(ctx context.Context, accessToken string) (*User, error)
}

// UserRepository ...
type UserRepository interface {
	CreateUser(ctx context.Context, username string, password string) error
	GetByUsername(ctx context.Context, username string) (*User, error)
}
