package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/go-kit/kit/log"
	"github.com/gogo/protobuf/proto"
	nats "github.com/nats-io/go-nats"
	pb "github.com/syedomair/micro-api/batch-tasks-svc/proto"
	pbp "github.com/syedomair/micro-api/public-svc/proto"
	pbu "github.com/syedomair/micro-api/users-svc/proto"
)

type Env struct {
	repo   Repository
	logger log.Logger
}

func (env *Env) HandleBatchGetAllClientUsers(m *nats.Msg) {

	start := time.Now()
	env.logger.Log("METHOD", "HandleBatchGetAllClientUsers", "SPOT", "method start", "time_start", start)
	batchTaskMsg := pb.BatchTaskMessage{}
	err := proto.Unmarshal(m.Data, &batchTaskMsg)
	if err != nil {
		env.logger.Log("METHOD", "HandleBatchGetAllClientUsers", "SPOT:", "proto.Unmarshal", "error:", err.Error())
	}

	bucketName := "batch-all-users-list-bucket"
	keyName := batchTaskMsg.ClientId + "/" + batchTaskMsg.BatchTaskId + ".json"

	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)

	users, count, err := env.repo.GetAllUser(batchTaskMsg.Limit, batchTaskMsg.Offset, batchTaskMsg.Orderby, batchTaskMsg.Sort, batchTaskMsg.ClientId)

	responseDataJson, err := json.Marshal(&pbu.UserList{Count: count, Offset: batchTaskMsg.Offset, Limit: batchTaskMsg.Limit, List: users})
	if err != nil {
		env.logger.Log("METHOD", "HandleBatchGetAllClientUsers", "SPOT:", "proto.Unmarshal", "error:", err.Error())
	}

	var reader io.Reader
	reader = strings.NewReader(string(responseDataJson))

	upParams := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &keyName,
		Body:   reader,
	}

	_, err = uploader.Upload(upParams)
	if err != nil {
		env.logger.Log("METHOD", "HandleBatchGetAllClientUsers", "SPOT:", "uploader.Upload", "error:", err.Error())
	}

	err = env.repo.MarkBatchTaskComplete(batchTaskMsg.BatchTaskId)
	if err != nil {
		env.logger.Log("METHOD", "HandleBatchGetAllClientUsers", "SPOT:", "proto.Unmarshal", "error:", err.Error())
	}

	env.logger.Log("METHOD", "HandleBatchGetAllClientUsers", "SPOT", "method end", "time_spent", time.Since(start))
	return
}

func (env *Env) HandleUserRegister(m *nats.Msg) {
	start := time.Now()
	env.logger.Log("METHOD", "HandleUserRegister", "SPOT", "method start", "time_start", start)
	env.logger.Log("METHOD", "HandleUserRegister", "Received on:", m.Subject, "data:", string(m.Data))
	userMsg := pbp.UserMessage{}
	err := proto.Unmarshal(m.Data, &userMsg)
	if err != nil {
		env.logger.Log("METHOD", "HandleUserRegister", "SPOT:", "proto.Unmarshal", "error:", err.Error())
	}
	data := url.Values{}
	data.Set("username", userMsg.UserId)
	request, err := http.NewRequest("POST", "http://kong-admin:8001/consumers", strings.NewReader(data.Encode()))
	httpErr := HttpKongCall(request, env.logger)
	if httpErr != nil {
		env.logger.Log("METHOD", "HandleUserRegister", "ERROR", err.Error())
	}
	env.logger.Log("METHOD", "HandleUserRegister", "SPOT", "method end", "time_spent", time.Since(start))
}

func (env *Env) HandleUserLogin(m *nats.Msg) {
	start := time.Now()
	env.logger.Log("METHOD", "HandleUserLogin", "SPOT", "method start", "time_start", start)
	env.logger.Log("METHOD", "HandleUserLogin", "Received on:", m.Subject, "data:", string(m.Data))
	userMsg := pbp.UserTokenMessage{}
	err := proto.Unmarshal(m.Data, &userMsg)
	if err != nil {
		env.logger.Log("METHOD", "HandleUserRegister", "SPOT:", "proto.Unmarshal", "error:", err.Error())
	}

	data := url.Values{}
	data.Set("key", userMsg.Token)
	request, err := http.NewRequest("POST", "http://kong-admin:8001/consumers/"+userMsg.UserId+"/key-auth", strings.NewReader(data.Encode()))
	httpErr := HttpKongCall(request, env.logger)
	if httpErr != nil {
		env.logger.Log("METHOD", "HandleUserLogin", "ERROR", err.Error())
	}
	env.logger.Log("METHOD", "HandleUserLogin", "SPOT", "method end", "time_spent", time.Since(start))
}

func (env *Env) HandleUserDelete(m *nats.Msg) {
	start := time.Now()
	env.logger.Log("METHOD", "HandleUserDelete", "SPOT", "method start", "time_start", start)
	env.logger.Log("METHOD", "HandleUserDelete", "Received on:", m.Subject, "data:", string(m.Data))
	userMsg := pbp.UserTokenMessage{}
	err := proto.Unmarshal(m.Data, &userMsg)
	if err != nil {
		env.logger.Log("METHOD", "HandleUserRegister", "SPOT:", "proto.Unmarshal", "error:", err.Error())
	}

	data := url.Values{}
	request, err := http.NewRequest("DELETE", "http://kong-admin:8001/consumers/"+userMsg.UserId, strings.NewReader(data.Encode()))
	httpErr := HttpKongCall(request, env.logger)
	if httpErr != nil {
		env.logger.Log("METHOD", "HandleUserDelete", "ERROR", err.Error())
	}
	env.logger.Log("METHOD", "HandleUserDelete", "SPOT", "method end", "time_spent", time.Since(start))
}

func HttpKongCall(request *http.Request, logger log.Logger) error {
	start := time.Now()
	logger.Log("METHOD", "HttpCall", "SPOT", "method start", "time_start", start)

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.Log("METHOD", "HttpCall", "error:", err.Error())
		return err
	}
	defer response.Body.Close()
	logger.Log("METHOD", "HttpCall", "SPOT:", "response.status", "status:", string(response.Status))

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Log("METHOD", "HttpCall", "error:", err.Error())
		return err
	}
	var bodyInterface map[string]interface{}
	json.Unmarshal(body, &bodyInterface)
	//logger.Log("METHOD", "HttpCall", "SPOT", "response", "body:", string(bodyInterface))
	logger.Log("METHOD", "HttpCall", "SPOT", "method end", "time_spent", time.Since(start))
	return nil
}
