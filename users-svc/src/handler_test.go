package main

import (
	"context"
	"testing"
	"time"

	"github.com/syedomair/micro-api/common"
	testdata "github.com/syedomair/micro-api/testdata"
	pb "github.com/syedomair/micro-api/users-svc/proto"
)

func TestGetAll(t *testing.T) {

	env := Env{repo: &mockDB{}, nats: &mockNATS{}, logger: common.GetLogger()}
	start := time.Now()
	env.logger.Log("METHOD", "TestGetAll", "SPOT", "method start", "time_start", start)
	ctx := context.WithValue(context.Background(), "client_id", testdata.ClientId)

	//ALL Good
	req := &pb.RequestQuery{Limit: "3", Offset: "0", Orderby: "title", Sort: "desc"}
	response, _ := env.GetAll(ctx, req)

	expected := testdata.SUCCESS
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//Invalid Offset
	req = &pb.RequestQuery{Limit: "3", Offset: "A", Orderby: "title", Sort: "desc"}
	response, _ = env.GetAll(ctx, req)

	expected = testdata.FAILURE
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//Invalid Limit
	req = &pb.RequestQuery{Limit: "A", Offset: "0", Orderby: "title", Sort: "desc"}
	response, _ = env.GetAll(ctx, req)

	expected = testdata.FAILURE
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//Invalid orderby
	req = &pb.RequestQuery{Limit: "3", Offset: "0", Orderby: "3", Sort: "desc"}
	response, _ = env.GetAll(ctx, req)

	expected = testdata.FAILURE
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//Invalid sort
	req = &pb.RequestQuery{Limit: "3", Offset: "0", Orderby: "title", Sort: "123"}
	response, _ = env.GetAll(ctx, req)

	expected = testdata.FAILURE
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}
	env.logger.Log("METHOD", "TestGetAll", "SPOT", "method end", "time_spent", time.Since(start))
}

func TestGetUser(t *testing.T) {

	env := Env{repo: &mockDB{}, nats: &mockNATS{}, logger: common.GetLogger()}
	start := time.Now()
	env.logger.Log("METHOD", "TestGetUser", "SPOT", "method start", "time_start", start)
	ctx := context.WithValue(context.Background(), "client_id", testdata.ClientId)

	//All Good
	req := &pb.User{Id: testdata.UserId}
	response, _ := env.GetUser(ctx, req)

	expected := testdata.SUCCESS
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//Invalid UserID
	req = &pb.User{Id: testdata.InValidId}
	response, _ = env.GetUser(ctx, req)

	expected = testdata.FAILURE
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}
	env.logger.Log("METHOD", "TestGetUser", "SPOT", "method end", "time_spent", time.Since(start))
}

func TestUpdateUser(t *testing.T) {

	env := Env{repo: &mockDB{}, nats: &mockNATS{}, logger: common.GetLogger()}
	start := time.Now()
	env.logger.Log("METHOD", "TestUpdateUser", "SPOT", "method start", "time_start", start)
	ctx := context.WithValue(context.Background(), "client_id", testdata.ClientId)

	//All Good
	user := &pb.User{Id: testdata.UserId, FirstName: testdata.ValidFirstName, LastName: testdata.ValidLastName}
	response, _ := env.UpdateUser(ctx, user)

	expected := testdata.SUCCESS
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//Invalid user_id
	user = &pb.User{Id: testdata.InValidId, FirstName: testdata.ValidFirstName}
	response, _ = env.UpdateUser(ctx, user)

	expected = testdata.FAILURE
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//ALL Invalid name
	user = &pb.User{Id: testdata.UserId, FirstName: testdata.ValidFirstName}
	response, _ = env.UpdateUser(ctx, user)

	expected = testdata.SUCCESS
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}
	env.logger.Log("METHOD", "TestUpdateRole", "SPOT", "method end", "time_spent", time.Since(start))
}

func TestDeleteUser(t *testing.T) {

	env := Env{repo: &mockDB{}, nats: &mockNATS{}, logger: common.GetLogger()}
	start := time.Now()
	env.logger.Log("METHOD", "TestDeleteUser", "SPOT", "method start", "time_start", start)
	ctx := context.WithValue(context.Background(), "client_id", testdata.ClientId)

	//All Good
	user := &pb.User{Id: testdata.UserId}
	response, _ := env.DeleteUser(ctx, user)

	expected := testdata.SUCCESS
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//Invalid user_id
	user = &pb.User{Id: testdata.InValidId}
	response, _ = env.DeleteUser(ctx, user)

	expected = testdata.FAILURE
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}
	env.logger.Log("METHOD", "TestDeleteUser", "SPOT", "method end", "time_spent", time.Since(start))
}

type mockNATS struct {
}

func (mnats *mockNATS) PublishUserDeleteEvent(userId string, clientId string) error {
	return nil
}

type mockDB struct {
	users []*pb.UserShorten
}

func (mdb *mockDB) Create(role *pb.User, clientId string) (string, error) {
	return testdata.UserId, nil
}
func (mdb *mockDB) GetAll(limit string, offset string, orderby string, sort string, clientId string) ([]*pb.UserShorten, string, error) {
	return mdb.users, "5", nil
}
func (mdb *mockDB) Get(roleId string, clientId string) (*pb.UserShorten, error) {
	user := &pb.UserShorten{}
	return user, nil
}
func (mdb *mockDB) Update(user *pb.User, clientId string) error {
	return nil
}
func (mdb *mockDB) Delete(user *pb.User, clientId string) error {
	return nil
}
