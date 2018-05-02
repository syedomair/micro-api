package main

import (
	"strconv"
	"testing"
	"time"

	"github.com/syedomair/micro-api/common"
	pb "github.com/syedomair/micro-api/public-svc/proto"
	testdata "github.com/syedomair/micro-api/testdata"
)

func TestPublicDB(t *testing.T) {

	db, _ := common.CreateDBConnection()
	repo := &PublicRepository{db, common.GetLogger()}
	defer repo.db.Close()

	start := time.Now()
	repo.logger.Log("METHOD", "TestPublicDB", "SPOT", "method start", "time_start", start)

	currentTime := time.Now()
	uniqueEmail = "email_" + strconv.FormatInt(currentTime.UnixNano(), 10) + "@gmail.com"
	user := &pb.User{
		FirstName: testdata.ValidFirstName,
		LastName:  testdata.ValidLastName,
		Email:     uniqueEmail,
		Password:  testdata.ValidPassword,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339)}

	userId, err := repo.Create(user, testdata.ClientId)
	var expected error = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}
	repo.logger.Log("METHOD", "TestPublicDB", "userId", userId)

	currentTime = time.Now()
	uniqueEmail1 := "email_" + strconv.FormatInt(currentTime.UnixNano(), 10) + "@gmail.com"

	userUniqueEmail := &pb.User{
		FirstName: testdata.ValidFirstName,
		LastName:  testdata.ValidLastName,
		Email:     uniqueEmail1,
		Password:  testdata.ValidPassword,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339)}
	err = repo.IsEmailUnique(userUniqueEmail, testdata.ClientId)
	expected = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}

	err = repo.IsEmailUnique(user, testdata.ClientId)
	expectedStr := "Email already exist for this client."
	if expectedStr != err.Error() {
		t.Errorf("\n...expected = %v\n...obtained = %v", expectedStr, err)
	}

	loginReq := &pb.LoginRequest{Email: uniqueEmail, Password: testdata.ValidPassword}
	userResponse, err := repo.Authenticate(loginReq, testdata.ClientId)
	expected = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}
	expectedString := userId
	if expectedString != userResponse.Id {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, userResponse.Id)
	}

	if err = repo.db.Delete(&user).Error; err != nil {
		repo.logger.Log("METHOD", "TestPublicDB", "Error in deleting", err)
	}
	repo.logger.Log("METHOD", "TestPublicDB", "SPOT", "method end", "time_spent", time.Since(start))
}
