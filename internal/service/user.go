package service

import (
	"errors"

	"go_web/internal/model"
	"go_web/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(name, email, password string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(id uint, name string, status int) (*model.User, error)
	DeleteUser(id uint) error
	ListUsers(page, pageSize int) ([]*model.User, int64, error)
	// 权限检查
	HasPermission(userID uint, resource, action string) (bool, error)
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

	// 使用 bcrypt 对密码进行哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
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

// HasPermission 检查用户是否拥有指定资源与操作的权限
func (s *userService) HasPermission(userID uint, resource, action string) (bool, error) {
	return s.userRepo.HasPermission(userID, resource, action)
}
