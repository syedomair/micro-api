package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	nats "github.com/nats-io/go-nats"
	common "github.com/syedomair/micro-api/common"
)

func main() {
	serverList := os.Getenv("NATS_SERVER")
	rootCACertFile := os.Getenv("NATS_CACERT")
	clientCertFile := os.Getenv("NATS_CERT")
	clientKeyFile := os.Getenv("NATS_KEY")
	// Connect options
	rootCA := nats.RootCAs(rootCACertFile)
	clientCert := nats.ClientCert(clientCertFile, clientKeyFile)
	alwaysReconnect := nats.MaxReconnects(-1)

	logger := common.GetLogger()
	db, _ := common.CreateDBConnection()

	repo := &BatchRepository{db, logger}
	env := &Env{repo, logger}

	var nc *nats.Conn
	var err1 error
	for {
		fmt.Println("for loop")
		nc, err1 = nats.Connect(serverList, rootCA, clientCert, alwaysReconnect)
		if err1 != nil {
			fmt.Printf("Error while connecting to NATS, backing off for a sec... (error: %s)", err1)
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Println("Connected NATS")
		break
	}

	nc.Subscribe("User.UserRegister", func(m *nats.Msg) {
		env.HandleUserRegister(m)
	})
	nc.Subscribe("User.UserLogin", func(m *nats.Msg) {
		env.HandleUserLogin(m)
	})
	nc.Subscribe("User.UserDelete", func(m *nats.Msg) {
		env.HandleUserDelete(m)
	})
	nc.Subscribe("Batch.GetAllUser", func(m *nats.Msg) {
		env.HandleBatchGetAllClientUsers(m)
	})
	runtime.Goexit()
}
