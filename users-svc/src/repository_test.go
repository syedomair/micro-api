package main

import (
	"testing"
	"time"

	"github.com/syedomair/micro-api/common"
	testdata "github.com/syedomair/micro-api/testdata"
	pb "github.com/syedomair/micro-api/users-svc/proto"
)

func TestUserDB(t *testing.T) {

	db, _ := common.CreateDBConnection()
	repo := &UserRepository{db, common.GetLogger()}
	defer repo.db.Close()

	start := time.Now()
	repo.logger.Log("METHOD", "TestUserDB", "SPOT", "method start", "time_start", start)

	//Create
	user := &pb.User{FirstName: testdata.ValidFirstName,
		LastName: testdata.ValidLastName,
		Email:    testdata.ValidEmail,
		Password: testdata.ValidPassword,
		IsAdmin:  testdata.IsAdmin}
	userId, err := repo.Create(user, testdata.ClientId)
	var expected error = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}
	repo.logger.Log("METHOD", "TestUserDB", "userId", userId)

	//Get
	userResponse, err := repo.Get(userId, testdata.ClientId)
	expected = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}
	repo.logger.Log("METHOD", "TestUserDB", "userResponse", userResponse)

	//Update
	user.FirstName = testdata.ValidFirstName + "-changed"
	err = repo.Update(user, testdata.ClientId)
	expected = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}

	//Get
	userResponse, err = repo.Get(userId, testdata.ClientId)
	expected = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}
	expectedString := testdata.ValidFirstName + "-changed"
	if expectedString != userResponse.FirstName {
		t.Errorf("\n...expected = %v\n...obtained = %v", expectedString, userResponse.FirstName)
	}
	repo.logger.Log("METHOD", "TestUserDB", "userResponse", userResponse)

	//GetAll
	users, _, err := repo.GetAll("1", "0", "created_at", "desc", testdata.ClientId)
	expectedInt := 1
	if expectedInt != len(users) {
		t.Errorf("\n...expected = %v\n...obtained = %v", expectedInt, len(users))
	}
	expected = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}

	//cleaning up after test
	if err = repo.db.Delete(&user).Error; err != nil {
		repo.logger.Log("METHOD", "TestUserDB", "Error in deleting", err)
	}
	repo.logger.Log("METHOD", "TestUserDB", "SPOT", "method end", "time_spent", time.Since(start))
}
