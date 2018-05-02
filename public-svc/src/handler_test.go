package main

import (
	"context"
	"strconv"
	"testing"
	"time"

	common "github.com/syedomair/micro-api/common"
	pb "github.com/syedomair/micro-api/public-svc/proto"
	testdata "github.com/syedomair/micro-api/testdata"
	"google.golang.org/grpc/metadata"
)

var uniqueEmail = ""

func TestCreate(t *testing.T) {
	env := Env{repo: &mockDB{}, nats: &mockNATS{}, logger: common.GetLogger()}
	md := metadata.New(map[string]string{"authorization": testdata.TestValidPublicToken})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	currentTime := time.Now()
	uniqueEmail = "email_" + strconv.FormatInt(currentTime.UnixNano(), 10) + "@gmail.com"

	user := &pb.User{
		FirstName: testdata.ValidFirstName,
		LastName:  testdata.ValidLastName,
		Email:     uniqueEmail,
		Password:  testdata.ValidPassword}
	response, _ := env.Register(ctx, user)

	//TEST 1 correct authorization
	expected := "success"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	md = metadata.New(map[string]string{"authorization": testdata.TestInValidPublicToken})
	ctx = metadata.NewIncomingContext(context.Background(), md)

	response, _ = env.Register(ctx, user)

	//TEST 2 incorrect authorization
	expected = "failure"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}
}

func TestAuthenticate(t *testing.T) {
	env := Env{repo: &mockDB{}, nats: &mockNATS{}, logger: common.GetLogger()}
	md := metadata.New(map[string]string{"authorization": testdata.TestValidPublicToken})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	user := &pb.LoginRequest{
		Email:    uniqueEmail,
		Password: testdata.ValidPassword}
	response, _ := env.Authenticate(ctx, user)

	//TEST 1 correct authorization
	expected := "success"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	md = metadata.New(map[string]string{"authorization": testdata.TestInValidPublicToken})
	ctx = metadata.NewIncomingContext(context.Background(), md)

	response, _ = env.Authenticate(ctx, user)

	//TEST 2 incorrect authorization
	expected = "failure"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//TEST 3 correct authorization with invalid password
	md = metadata.New(map[string]string{"authorization": testdata.TestInValidPublicToken})
	ctx = metadata.NewIncomingContext(context.Background(), md)

	user = &pb.LoginRequest{
		Email:    uniqueEmail,
		Password: testdata.InValidPassword}
	response, _ = env.Authenticate(ctx, user)

	expected = "failure"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//TEST 4 correct authorization with invalid email
	user = &pb.LoginRequest{
		Email:    testdata.InValidEmail,
		Password: testdata.ValidPassword}
	response, _ = env.Authenticate(ctx, user)

	expected = "failure"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}
}

type mockNATS struct {
}

func (mnats *mockNATS) PublishRegisterEvent(userId string, clientId string) error {
	return nil
}
func (mnats *mockNATS) PublishAuthEvent(userId string, token string) error {
	return nil
}

type mockDB struct {
}

func (mdb *mockDB) IsEmailUnique(user *pb.User, clientId string) error {
	return nil
}
func (mdb *mockDB) Create(user *pb.User, clientId string) (string, error) {
	return testdata.UserId, nil
}
func (mdb *mockDB) Authenticate(user *pb.LoginRequest, clientId string) (*pb.User, error) {
	return &pb.User{
		Id:        testdata.UserId,
		ClientId:  testdata.ClientId,
		FirstName: testdata.ValidFirstName,
		LastName:  testdata.ValidLastName,
		Email:     testdata.ValidEmail,
		Password:  testdata.ValidPassword,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339)}, nil
}
func (mdb *mockDB) GetClientFromApiKey(apiKey string) (*pb.Client, error) {
	return &pb.Client{
		Id:        testdata.ClientId,
		Name:      testdata.ClientName,
		ApiKey:    testdata.APIKey,
		Secret:    testdata.APISecret,
		Status:    testdata.ClientStatus,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339)}, nil

}

/*
func (mdb *mockDB) initMockDb() ([]*pb.User, error) {
	users := make([]*pb.User, 0)
	users = append(users, &pb.User{Id: "04b58e6e-f910-4ff0-83f1-27fbfa85dc2f", ClientId: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", FirstName: "First Name 1", LastName: "Last Name 1", Email: "email1@gmail.com", Password: "123", CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)})
	users = append(users, &pb.User{Id: "04b58e6e-f910-4ff0-83f1-27fbfa85dc2f", ClientId: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", FirstName: "First Name 2", LastName: "Last Name 2", Email: "email2@gmail.com", Password: "123",  CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)})
	users = append(users, &pb.User{Id: "04b58e6e-f910-4ff0-83f1-27fbfa85dc2f", ClientId: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", FirstName: "First Name 3", LastName: "Last Name 3", Email: "email3@gmail.com", Password: "123",  CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)})
	return users, nil
}
*/
