package main

import (
	"time"

	log "github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	pb "github.com/syedomair/micro-api/batch-tasks-svc/proto"
	common "github.com/syedomair/micro-api/common"
)

type Repository interface {
	Create(apiName string, clientId string) (*pb.BatchTask, error)
	Get(batchTaskId string, clientId string) (*pb.BatchTask, error)
}

type BatchTaskRepository struct {
	db     *gorm.DB
	logger log.Logger
}

func (repo *BatchTaskRepository) Create(apiName string, clientId string) (*pb.BatchTask, error) {
	start := time.Now()
	repo.logger.Log("METHOD", "Create", "SPOT", "method start", "time_start", start)
	batchTaskId := uuid.NewV4().String()
	batchTask := &pb.BatchTask{
		Id:          batchTaskId,
		ClientId:    clientId,
		ApiName:     apiName,
		Status:      common.BATCH_TASK_PENDING,
		CreatedAt:   time.Now().Format(time.RFC3339),
		CompletedAt: "1111-11-11T10:00:00-05:00"}

	if err := repo.db.Create(batchTask).Error; err != nil {
		return nil, err
	}
	repo.logger.Log("METHOD", "Create", "SPOT", "method end", "time_spent", time.Since(start))
	return batchTask, nil
}
func (repo *BatchTaskRepository) Get(batchTaskId string, clientId string) (*pb.BatchTask, error) {
	start := time.Now()
	repo.logger.Log("METHOD", "Get", "SPOT", "method start", "time_start", start)
	batchTask := pb.BatchTask{}
	if err := repo.db.Where("client_id = ?", clientId).Where("id = ?", batchTaskId).Find(&batchTask).Error; err != nil {
		return nil, err
	}
	repo.logger.Log("METHOD", "Get", "SPOT", "method end", "time_spent", time.Since(start))
	return &batchTask, nil
}
