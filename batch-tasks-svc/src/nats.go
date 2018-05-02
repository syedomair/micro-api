package main

import (
	log "github.com/go-kit/kit/log"
	"github.com/gogo/protobuf/proto"
	nats "github.com/nats-io/go-nats"
	pb "github.com/syedomair/micro-api/batch-tasks-svc/proto"
)

type Nats interface {
	PublishBatchGetAllUserEvent(string, string, string, string, string, string, string, string) error
}

type NatsWrapper struct {
	nats   *nats.Conn
	logger log.Logger
}

func (natsWrap *NatsWrapper) PublishBatchGetAllUserEvent(batchTaskId string, clientId string, limit string, offset string, orderby string, sort string, filter string, search string) error {

	natsWrap.logger.Log("ACTION", "PublishBatchGetAllUserEvent", "SPOT", "method start")
	batchTaskMessage := pb.BatchTaskMessage{
		BatchTaskId: batchTaskId,
		ClientId:    clientId,
		Limit:       limit,
		Offset:      offset,
		Orderby:     orderby,
		Sort:        sort,
		Filter:      filter,
		Search:      search}
	data, _ := proto.Marshal(&batchTaskMessage)
	err := natsWrap.nats.Publish("Batch.GetAllUser", data)
	if err != nil {
		natsWrap.logger.Log("Error during publishing: ", err)
		return err
	}
	natsWrap.nats.Flush()

	natsWrap.logger.Log("ACTION", "PublishBatchGetAllUserEvent", "SPOT", "method end")
	return nil
}
