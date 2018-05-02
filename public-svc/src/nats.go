package main

import (
	log "github.com/go-kit/kit/log"
	"github.com/gogo/protobuf/proto"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	nats "github.com/nats-io/go-nats"
	pb "github.com/syedomair/micro-api/public-svc/proto"
)

type Nats interface {
	PublishRegisterEvent(string, string) error
	PublishAuthEvent(string, string) error
}

type NatsWrapper struct {
	nats   *nats.Conn
	logger log.Logger
}

func (natsWrap *NatsWrapper) PublishRegisterEvent(userId string, clientId string) error {

	natsWrap.logger.Log("ACTION", "PublishRegisterEvent", "SPOT", "method start")
	userMessage := pb.UserMessage{UserId: userId, ClientId: clientId}
	data, _ := proto.Marshal(&userMessage)
	err := natsWrap.nats.Publish("User.UserRegister", data)
	if err != nil {
		natsWrap.logger.Log("Error during publishing: ", err)
		return err
	}
	natsWrap.nats.Flush()

	natsWrap.logger.Log("ACTION", "PublishRegisterEvent", "SPOT", "method end")
	return nil
}

func (natsWrap *NatsWrapper) PublishAuthEvent(userId string, signedJwtToken string) error {

	natsWrap.logger.Log("ACTION", "PublishAuthEvent", "SPOT", "method start")
	userMessage := pb.UserTokenMessage{UserId: userId, Token: signedJwtToken}
	data, _ := proto.Marshal(&userMessage)
	err := natsWrap.nats.Publish("User.UserLogin", data)
	if err != nil {
		natsWrap.logger.Log("Error during publishing: ", err)
		return err
	}
	natsWrap.nats.Flush()
	natsWrap.logger.Log("ACTION", "PublishAuthEvent", "SPOT", "method end")

	return nil
}
