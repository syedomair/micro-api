package main

import (
	"context"
	"testing"
	"time"

	pb "github.com/syedomair/micro-api/batch-tasks-svc/proto"
	"github.com/syedomair/micro-api/common"
	testdata "github.com/syedomair/micro-api/testdata"
)

func TestGetAllUser(t *testing.T) {

	env := Env{repo: &mockDB{}, nats: &mockNATS{}, logger: common.GetLogger()}
	start := time.Now()
	env.logger.Log("METHOD", "TestGetAllUser", "SPOT", "method start", "time_start", start)
	ctx := context.WithValue(context.Background(), "client_id", testdata.ClientId)

	//ALL Good
	req := &pb.RequestQuery{Limit: "3", Offset: "0", Orderby: "title", Sort: "desc"}
	response, _ := env.GetAllUser(ctx, req)

	expected := testdata.SUCCESS
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	env.logger.Log("METHOD", "TestGetAllUserTest", "SPOT", "method end", "time_spent", time.Since(start))
}

func TestGetAllUserStatus(t *testing.T) {

	env := Env{repo: &mockDB{}, nats: &mockNATS{}, logger: common.GetLogger()}
	start := time.Now()
	env.logger.Log("METHOD", "TestGetAllUserStatus", "SPOT", "method start", "time_start", start)
	ctx := context.WithValue(context.Background(), "client_id", testdata.ClientId)

	//All Good
	req := &pb.BatchTask{Id: testdata.UserId}
	response, _ := env.GetAllUserStatus(ctx, req)

	expected := testdata.SUCCESS
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	env.logger.Log("METHOD", "TestGetAllUserStatus", "SPOT", "method end", "time_spent", time.Since(start))
}

type mockNATS struct {
}

type mockDB struct {
}

func (mdb *mockDB) Create(apiName string, clientId string) (*pb.BatchTask, error) {
	return nil, nil
}
func (mdb *mockDB) Get(batchTaskId string, clientId string) (*pb.BatchTask, error) {
	return nil, nil
}
