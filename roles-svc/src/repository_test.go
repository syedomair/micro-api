package main

import (
	"testing"
	"time"

	"github.com/syedomair/micro-api/common"
	pb "github.com/syedomair/micro-api/roles-svc/proto"
	testdata "github.com/syedomair/micro-api/testdata"
)

func TestRoleDB(t *testing.T) {

	db, _ := common.CreateDBConnection()
	repo := &RoleRepository{db, common.GetLogger()}
	defer repo.db.Close()

	start := time.Now()
	repo.logger.Log("METHOD", "TestRoleDB", "SPOT", "method start", "time_start", start)

	//Create
	role := &pb.Role{Title: testdata.RoleTitle1, RoleType: testdata.RoleType}
	roleId, err := repo.Create(role, testdata.ClientId)
	var expected error = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}
	repo.logger.Log("METHOD", "TestRoleDB", "roleId", roleId)

	//Get
	roleResponse, err := repo.Get(roleId, testdata.ClientId)
	expected = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}
	repo.logger.Log("METHOD", "TestRoleDB", "roleResponse", roleResponse)

	//Update
	role.Title = testdata.RoleTitle2
	err = repo.Update(role, testdata.ClientId)
	expected = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}

	//Get
	roleResponse, err = repo.Get(roleId, testdata.ClientId)
	expected = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}
	expectedString := testdata.RoleTitle2
	if expectedString != roleResponse.Title {
		t.Errorf("\n...expected = %v\n...obtained = %v", expectedString, roleResponse.Title)
	}
	repo.logger.Log("METHOD", "TestRoleDB", "roleResponse", roleResponse)

	//GetAll
	roles, _, err := repo.GetAll("1", "0", "created_at", "desc", testdata.ClientId)
	expectedInt := 1
	if expectedInt != len(roles) {
		t.Errorf("\n...expected = %v\n...obtained = %v", expectedInt, len(roles))
	}
	expected = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}

	//cleaning up after test
	if err = repo.db.Delete(&role).Error; err != nil {
		repo.logger.Log("METHOD", "TestRoleDB", "Error in deleting", err)
	}

	repo.logger.Log("METHOD", "TestRoleDB", "SPOT", "method end", "time_spent", time.Since(start))
}
