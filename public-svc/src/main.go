package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	common "github.com/syedomair/micro-api/common"
	pb "github.com/syedomair/micro-api/public-svc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	errors := make(chan error)
	httpPort := "8180"
	grpcPort := "50080"
	fmt.Println("HTTP PORT", httpPort)
	fmt.Println("GRPC PORT", grpcPort)

	go func() { errors <- startGRPC(grpcPort) }()
	go func() { errors <- startHTTP(httpPort, grpcPort) }()
	for err := range errors {
		log.Fatal(err)
		return
	}
}
func startGRPC(port string) error {
	db, err := common.CreateDBConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	} else {
		fmt.Println("Connected to DB")
	}

	logger := common.GetLogger()

	repo := &PublicRepository{db, logger}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	natsCon, _ := common.CreateNATSConnection()
	nats := &NatsWrapper{natsCon, logger}

	s := grpc.NewServer()
	pb.RegisterPublicServiceServer(s, &Env{repo, nats, logger})

	return s.Serve(lis)
}
func startHTTP(httpPort, grpcPort string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := pb.RegisterPublicServiceHandlerFromEndpoint(ctx, gwmux, "127.0.0.1:"+grpcPort, opts); err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/v1/", gwmux)
	http.ListenAndServe(":"+httpPort, mux)
	return nil
}
