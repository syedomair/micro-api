package main

import (
	log "github.com/go-kit/kit/log"
	nats "github.com/nats-io/go-nats"
)

type Nats interface {
	PublishDeleteEvent(string, string) error
}

type NatsWrapper struct {
	nats   *nats.Conn
	logger log.Logger
}

func (natsWrap *NatsWrapper) PublishDeleteEvent(roleId string, clientId string) error {

	natsWrap.logger.Log("ACTION", "PublishDeleteEvent", "SPOT", "method start")

	natsWrap.logger.Log("ACTION", "PublishDeleteEvent", "SPOT", "method end")
	return nil
}
