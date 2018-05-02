package main

import (
	log "github.com/go-kit/kit/log"
	"github.com/gogo/protobuf/proto"
	nats "github.com/nats-io/go-nats"
	pb "github.com/syedomair/micro-api/users-svc/proto"
)

type Nats interface {
	PublishUserDeleteEvent(string, string) error
}

type NatsWrapper struct {
	nats   *nats.Conn
	logger log.Logger
}

func (natsWrap *NatsWrapper) PublishUserDeleteEvent(userId string, clientId string) error {

	natsWrap.logger.Log("ACTION", "PublishDeleteEvent", "SPOT", "method start")
	userMessage := pb.UserMessage{UserId: userId, ClientId: clientId}
	data, _ := proto.Marshal(&userMessage)
	err := natsWrap.nats.Publish("User.UserDelete", data)
	if err != nil {
		natsWrap.logger.Log("Error during publishing: ", err)
		return err
	}
	natsWrap.nats.Flush()

	natsWrap.logger.Log("ACTION", "PublishDeleteEvent", "SPOT", "method end")
	return nil
}
