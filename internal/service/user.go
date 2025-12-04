package service

import (
	"errors"

	"go_test/internal/model"
	"go_test/internal/repository"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(name, email, password string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(id uint, name string, status int) (*model.User, error)
	DeleteUser(id uint) error
	ListUsers(page, pageSize int) ([]*model.User, int64, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(name, email, password string) (*model.User, error) {
	// 检查邮箱是否已存在
	existingUser, err := s.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("邮箱已存在")
	}
	// 如果查询出错但不是"记录不存在"的错误，应该返回错误
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: password, // 实际应用中应该加密
		Status:   1,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *userService) UpdateUser(id uint, name string, status int) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 只有当name不为空时才更新
	if name != "" {
		user.Name = name
	}
	// status: -1表示不更新，0或1表示更新
	if status >= 0 {
		user.Status = status
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *userService) ListUsers(page, pageSize int) ([]*model.User, int64, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(offset, pageSize)
}

