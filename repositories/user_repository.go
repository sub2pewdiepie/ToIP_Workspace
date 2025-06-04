package repositories

import (
	"space/models"
	"space/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByUsernameOrEmail(username, email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Create(user *models.User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) GetByUsernameOrEmail(username, email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ? OR email = ?", username, email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// repositories/user_repository.go
func (r *userRepo) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Logger.WithFields(logrus.Fields{
				"username": username,
			}).Debug("User not found")
			return nil, err
		}
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
		}).Error("Failed to find user by username")
		return nil, err
	}
	return &user, nil
}
