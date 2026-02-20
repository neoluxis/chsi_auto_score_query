package repo

import (
"chsi-auto-score-query/internal/logger"
"chsi-auto-score-query/internal/model"
"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		logger.Error("Failed to create user: %v", err)
		return err
	}
	return nil
}

func (r *UserRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logger.Error("Failed to find user by id: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logger.Error("Failed to find user by email: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindPending() ([]model.User, error) {
	var users []model.User
	if err := r.db.Where("score IS NULL OR score = ''").Find(&users).Error; err != nil {
		logger.Error("Failed to find pending users: %v", err)
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) Update(user *model.User) error {
	if err := r.db.Save(user).Error; err != nil {
		logger.Error("Failed to update user: %v", err)
		return err
	}
	return nil
}

func (r *UserRepo) Delete(id uint) error {
	if err := r.db.Delete(&model.User{}, id).Error; err != nil {
		logger.Error("Failed to delete user: %v", err)
		return err
	}
	return nil
}
