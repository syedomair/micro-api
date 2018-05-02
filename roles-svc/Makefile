
build:
	protoc -I/usr/local/include -I. \
	  -I$(GOPATH)/src \
	  -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
          --go_out=google/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
	  proto/roles.proto 

reverse:
	protoc -I/usr/local/include -I. \
          -I$(GOPATH)/src \
          -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
          --grpc-gateway_out=logtostderr=true:. \
	  proto/roles.proto

docker:
	docker build -t ${AWS_ECR}/${REPOSITORY_NAME} .

race:
	go test ./... -race

test:
	go test ./... -v
