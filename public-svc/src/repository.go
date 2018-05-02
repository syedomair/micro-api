package main

import (
	"errors"
	"time"

	b64 "encoding/base64"

	log "github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	pb "github.com/syedomair/micro-api/public-svc/proto"
)

type Repository interface {
	Create(user *pb.User, clientId string) (string, error)
	IsEmailUnique(user *pb.User, clientId string) error
	Authenticate(req *pb.LoginRequest, clientId string) (*pb.User, error)
	GetClientFromApiKey(apiKey string) (*pb.Client, error)
}

type PublicRepository struct {
	db     *gorm.DB
	logger log.Logger
}

func (repo *PublicRepository) GetClientFromApiKey(apikey string) (*pb.Client, error) {

	repo.logger.Log("METHOD", "GetClientFromApiKey", "SPOT", "method start")
	client := pb.Client{}
	if err := repo.db.Where("api_key = ?", apikey).Find(&client).Error; err != nil {
		return nil, err
	}
	repo.logger.Log("METHOD", "GetClientFromApiKey", "SPOT", "method end")
	return &client, nil
}

func (repo *PublicRepository) IsEmailUnique(user *pb.User, clientId string) error {

	repo.logger.Log("METHOD", "IsEmailUnique", "SPOT", "method start")
	if err := repo.db.Where("client_id = ?", clientId).Where("email = ?", user.Email).Find(&user).Error; err == nil {
		return errors.New("Email already exist for this client.")
	}
	repo.logger.Log("METHOD", "IsEmailUnique", "SPOT", "method end")
	return nil
}
func (repo *PublicRepository) Create(user *pb.User, clientId string) (string, error) {

	repo.logger.Log("METHOD", "Create", "SPOT", "method start")
	userId := uuid.NewV4().String()
	password := b64.StdEncoding.EncodeToString([]byte(user.Password))
	user = &pb.User{
		Id:        userId,
		ClientId:  clientId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  password,
		IsAdmin:   "0",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339)}

	if err := repo.db.Create(user).Error; err != nil {
		return "", err
	}
	repo.logger.Log("METHOD", "Create", "SPOT", "method end")
	return userId, nil
}

func (repo *PublicRepository) Authenticate(req *pb.LoginRequest, clientId string) (*pb.User, error) {

	repo.logger.Log("METHOD", "Authenticate", "SPOT", "method start")
	user := pb.User{}
	password := b64.StdEncoding.EncodeToString([]byte(req.Password))
	if err := repo.db.Where("client_id = ?", clientId).Where("email = ?", req.Email).Where("password = ?", password).Find(&user).Error; err != nil {
		return nil, err
	}
	repo.logger.Log("METHOD", "Authenticate", "SPOT", "method end")
	return &user, nil
}
