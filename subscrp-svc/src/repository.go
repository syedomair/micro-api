package main

import (
	"time"

	log "github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	pb "github.com/syedomair/micro-api/batch-tasks-svc/proto"
	common "github.com/syedomair/micro-api/common"
	pbu "github.com/syedomair/micro-api/users-svc/proto"
)

type Repository interface {
	GetAllUser(limit string, offset string, orderby string, sort string, clientId string) ([]*pbu.UserShorten, string, error)
	MarkBatchTaskComplete(string) error
}

type BatchRepository struct {
	db     *gorm.DB
	logger log.Logger
}

func (repo *BatchRepository) GetAllUser(limit string, offset string, orderby string, sort string, clientId string) ([]*pbu.UserShorten, string, error) {
	start := time.Now()
	repo.logger.Log("METHOD", "GetAllUser", "SPOT", "method start", "time_start", start)
	var users []*pbu.UserShorten
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
	repo.logger.Log("METHOD", "GetAllUser", "SPOT", "method end", "time_spent", time.Since(start))
	return users, count, nil
}

func (repo *BatchRepository) MarkBatchTaskComplete(batchTaskId string) error {
	start := time.Now()
	repo.logger.Log("METHOD", "MarkBatchTaskComplete", "SPOT", "method start", "time_start", start)
	batchTask := &pb.BatchTask{
		Id:          batchTaskId,
		CompletedAt: time.Now().Format(time.RFC3339),
		Status:      common.BATCH_TASK_COMPLETE}
	if err := repo.db.Model(batchTask).Update(&batchTask).Error; err != nil {
		return err
	}
	repo.logger.Log("METHOD", "MarkBatchTaskComplete", "SPOT", "method end", "time_spent", time.Since(start))
	return nil
}
