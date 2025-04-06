package services

import (
	"github.com/hsiaocz/web3-code/book/types"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	CreateUser(user *types.Users) (*types.Users, error)
	GetUserByID(userID string) (*types.Users, error)
	UpdateUser(user *types.Users) (*types.Users, error)
	DeleteUser(userID string) error
	GetAllUsers() ([]types.Users, error)
	GetUserByUsername(username string) (*types.Users, error)
	GetUserByEmail(email string) (*types.Users, error)
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}
func (s *UserService) CreateUser(user *types.Users) (*types.Users, error) {
	err := s.db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserService) GetUserByID(userID string) (*types.Users, error) {
	var user types.Users
	err := s.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (s *UserService) UpdateUser(user *types.Users) (*types.Users, error) {
	err := s.db.Save(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserService) DeleteUser(userID string) error {
	var user types.Users
	err := s.db.Where("user_id = ?", userID).Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}
func (s *UserService) GetAllUsers() ([]types.Users, error) {
	var users []types.Users
	err := s.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (s *UserService) GetUserByUsername(username string) (*types.Users, error) {
	var user types.Users
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (s *UserService) GetUserByEmail(email string) (*types.Users, error) {
	var user types.Users
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
