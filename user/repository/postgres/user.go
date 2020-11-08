package postgrel

import (
	"context"

	"github.com/bekaza/go-clean/domain"
	"gorm.io/gorm"
)

type postgrelUserRepo struct {
	DB *gorm.DB
}

// NewPostgrelUserRepository ...
func NewPostgrelUserRepository(db *gorm.DB) domain.UserRepository {
	db.AutoMigrate(&domain.User{})
	return &postgrelUserRepo{db}
}

func (mur *postgrelUserRepo) CreateUser(ctx context.Context, username string, password string) error {
	user := domain.User{
		Username: username,
		Password: password,
	}
	return mur.DB.Create(&user).Error
}

func (mur *postgrelUserRepo) GetByUsername(ctx context.Context, username string) (user *domain.User, err error) {
	err = mur.DB.First(&user).Error
	return user, err
}
