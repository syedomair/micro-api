package main

import (
	"testing"
	"time"

	"github.com/syedomair/api_micro/common"
	testdata "github.com/syedomair/micro-api/testdata"
)

func TestBatchTaskDB(t *testing.T) {
	db, _ := common.CreateDBConnection()
	repo := &BatchTaskRepository{db, common.GetLogger()}
	defer repo.db.Close()

	start := time.Now()
	repo.logger.Log("METHOD", "TestBatchTaskDB", "SPOT", "method start", "time_start", start)

	//Create
	batchTask, err := repo.Create("test_batch_api_name", testdata.ClientId)
	var expected error = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}
	repo.logger.Log("METHOD", "TestBatchTaskDB", "batch_task_id", batchTask.Id)

	//Get
	batchTask, err = repo.Get(batchTask.Id, testdata.ClientId)
	expected = nil
	if expected != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, err)
	}
	repo.logger.Log("METHOD", "TestBatchTaskDB", "batchTaskResponse.status", batchTask.Status)

	repo.logger.Log("METHOD", "TestBatchTaskDB", "SPOT", "method end", "time_spent", time.Since(start))
}
