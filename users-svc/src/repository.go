package main

import (
	"time"

	log "github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	pb "github.com/syedomair/micro-api/users-svc/proto"
)

type Repository interface {
	Create(user *pb.User, clientId string) (string, error)
	GetAll(limit string, offset string, orderby string, sort string, clientId string) ([]*pb.UserShorten, string, error)
	Get(userId string, clientId string) (*pb.UserShorten, error)
	Update(user *pb.User, clientId string) error
	Delete(user *pb.User, clientId string) error
}

type UserRepository struct {
	db     *gorm.DB
	logger log.Logger
}

func (repo *UserRepository) Create(user *pb.User, clientId string) (string, error) {

	start := time.Now()
	repo.logger.Log("METHOD", "Create", "SPOT", "method start", "time_start", start)

	userId := uuid.NewV4().String()
	user = &pb.User{
		Id:        userId,
		ClientId:  clientId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		IsAdmin:   user.IsAdmin,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339)}

	if err := repo.db.Create(user).Error; err != nil {
		return "", err
	}
	repo.logger.Log("METHOD", "Create", "SPOT", "method end", "time_spent", time.Since(start))
	return userId, nil
}
func (repo *UserRepository) GetAll(limit string, offset string, orderby string, sort string, clientId string) ([]*pb.UserShorten, string, error) {

	start := time.Now()
	repo.logger.Log("METHOD", "GetAll", "SPOT", "method start", "time_start", start)
	var users []*pb.UserShorten
	count := "0"
	if err := repo.db.Table("users").
		Select("*").
		Count(&count).
		Limit(limit).
		Offset(offset).
		Order(orderby+" "+sort).
		Where("client_id = ?", clientId).
		Scan(&users).Error; err != nil {
		return nil, "", err
	}
	repo.logger.Log("METHOD", "GetAll", "SPOT", "method end", "time_spent", time.Since(start))
	return users, count, nil
}

func (repo *UserRepository) Get(userId string, clientId string) (*pb.UserShorten, error) {
	start := time.Now()
	repo.logger.Log("METHOD", "Get", "SPOT", "method start", "time_start", start)
	user := pb.User{}
	if err := repo.db.Where("client_id = ?", clientId).Where("id = ?", userId).Find(&user).Error; err != nil {
		return nil, err
	}
	userShorten := pb.UserShorten{
		Id:        user.Id,
		ClientId:  user.ClientId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt}

	repo.logger.Log("METHOD", "Get", "SPOT", "method end", "time_spent", time.Since(start))
	return &userShorten, nil
}

func (repo *UserRepository) Update(user *pb.User, clientId string) error {
	start := time.Now()
	repo.logger.Log("METHOD", "Update", "SPOT", "method start", "time_start", start)
	if err := repo.db.Model(user).Update(&user).Error; err != nil {
		return err
	}
	repo.logger.Log("METHOD", "Update", "SPOT", "method end", "time_spent", time.Since(start))
	return nil
}

func (repo *UserRepository) Delete(user *pb.User, clientId string) error {
	start := time.Now()
	repo.logger.Log("METHOD", "Delete", "SPOT", "method start", "time_start", start)
	userId := user.Id
	if err := repo.db.Where("client_id = ?", clientId).Where("id = ?", userId).Find(&user).Error; err != nil {
		return err
	}
	if err := repo.db.Delete(&user).Error; err != nil {
		return err
	}
	repo.logger.Log("METHOD", "Delete", "SPOT", "method end", "time_spent", time.Since(start))
	return nil
}
