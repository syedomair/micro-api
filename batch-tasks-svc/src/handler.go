package main

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/go-kit/kit/log"
	"golang.org/x/net/context"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	pb "github.com/syedomair/micro-api/batch-tasks-svc/proto"
	"github.com/syedomair/micro-api/common"
)

type Env struct {
	repo   Repository
	nats   Nats
	logger log.Logger
}

func (env *Env) GetAllUser(ctx context.Context, req *pb.RequestQuery) (*pb.ResponseBatchTask, error) {

	start := time.Now()
	env.logger.Log("METHOD", "GetAllUser", "SPOT", "method start", "time_start", start)
	clientId, _ := ctx.Value("client_id").(string)

	limit, offset, orderby, sort, err := common.ValidateQueryString(req.Limit, "10", req.Offset, "0", req.Orderby, "created_at", req.Sort, "desc")
	if err != nil {
		return &pb.ResponseBatchTask{Result: common.FAILURE, Error: common.CommonError(err.Error()), Data: nil}, nil
	}
	batchTask, err := env.repo.Create("batch-network-users-list", clientId)
	if err != nil {
		return &pb.ResponseBatchTask{Result: common.FAILURE, Data: nil, Error: common.DatabaseError()}, nil
	}

	env.logger.Log("METHOD", "GetAllUser", "SPOT", "before NATS event", "time_spent", time.Since(start))
	//NATS Event Publish
	go func() {
		natsError := env.nats.PublishBatchGetAllUserEvent(batchTask.Id, clientId, limit, offset, orderby, sort, req.Filter, req.Search)
		if natsError != nil {
			env.logger.Log("Error during publishing: ", natsError)
		}
	}()
	env.logger.Log("METHOD", "GetAllUser", "SPOT", "after NATS event", "time_spent", time.Since(start))

	env.logger.Log("METHOD", "GetAllUser", "SPOT", "method end", "time_spent", time.Since(start))
	return &pb.ResponseBatchTask{Result: common.SUCCESS, Error: nil, Data: batchTask}, nil
}

func (env *Env) GetAllUserStatus(ctx context.Context, req *pb.BatchTask) (*pb.ResponseBatchTask, error) {

	start := time.Now()
	env.logger.Log("METHOD", "GetAllUserStatus", "SPOT", "method start", "time_start", start)
	clientId, _ := ctx.Value("client_id").(string)

	if err := validateBatchTaskId(req); err != nil {
		return &pb.ResponseBatchTask{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2004", err.Error())}, nil
	}

	batchTask, err := env.repo.Get(req.Id, clientId)
	if err != nil {
		return &pb.ResponseBatchTask{Result: common.FAILURE, Data: nil, Error: common.CommonError(err.Error())}, nil
	}
	env.logger.Log("METHOD", "GetAllUserStatus", "SPOT", "method end", "time_spent", time.Since(start))
	return &pb.ResponseBatchTask{Result: common.SUCCESS, Data: batchTask, Error: nil}, nil
}

func (env *Env) GetAllUserOutput(ctx context.Context, req *pb.BatchTask) (*pb.ResponseList, error) {

	start := time.Now()
	env.logger.Log("METHOD", "GetAllUserOutput", "SPOT", "method start", "time_start", start)
	clientId, _ := ctx.Value("client_id").(string)

	if err := validateBatchTaskId(req); err != nil {
		return &pb.ResponseList{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2004", err.Error())}, nil
	}

	sess := session.Must(session.NewSession())
	bucketName := "batch-all-users-list-bucket"
	keyName := clientId + "/" + req.Id + ".json"

	downloader := s3manager.NewDownloader(sess)
	buff := &aws.WriteAtBuffer{}
	_, _ = downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(keyName),
		})
	userList := &pb.UserList{}
	if err := json.Unmarshal(buff.Bytes(), userList); err != nil {
		return &pb.ResponseList{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2004", err.Error())}, nil
	}

	env.logger.Log("METHOD", "GetAllUserOutput", "SPOT", "method end", "time_spent", time.Since(start))
	return &pb.ResponseList{Result: common.SUCCESS, Error: nil, Data: userList}, nil
}
func validateBatchTaskId(batchTask *pb.BatchTask) error {
	if err := validation.Validate(
		batchTask.Id,
		validation.Required.Error("batch_task_id is a required field"),
		is.UUIDv4.Error("invalid batch_task_id.")); err != nil {
		return err
	}
	return nil
}
