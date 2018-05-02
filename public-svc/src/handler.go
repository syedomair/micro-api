package main

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/go-kit/kit/log"
	"google.golang.org/grpc/metadata"

	"golang.org/x/net/context"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	common "github.com/syedomair/micro-api/common"
	pb "github.com/syedomair/micro-api/public-svc/proto"
)

type Env struct {
	repo   Repository
	nats   Nats
	logger log.Logger
}

func (env *Env) Register(ctx context.Context, req *pb.User) (*pb.Response, error) {
	env.logger.Log("METHOD", "Register", "SPOT", "method start")
	start := time.Now()
	meta, _ := metadata.FromIncomingContext(ctx)
	apiKey, err := common.GetAPIKey(meta["authorization"][0])
	if err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2001", err.Error())}, nil
	}

	client, err := env.repo.GetClientFromApiKey(apiKey)
	if err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2002", err.Error())}, nil
	}
	_, err = common.ValidateJWTToken(meta["authorization"][0], client.Secret)
	if err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2003", err.Error())}, nil
	}
	clientId := client.Id

	err = validateParameters(req)
	if err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2004", err.Error())}, nil
	}

	err = env.repo.IsEmailUnique(req, clientId)
	if err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2014", err.Error())}, nil
	}

	userId, err := env.repo.Create(req, clientId)
	if err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2005", err.Error())}, nil
	}
	responseUserId := map[string]string{"user_id": userId}

	env.logger.Log("METHOD", "Register", "SPOT", "before NATS event", "time_spent", time.Since(start))
	//NATS Event Publish
	go func() {
		natsError := env.nats.PublishRegisterEvent(userId, clientId)
		if natsError != nil {
			env.logger.Log("Error during publishing: ", natsError)
		}
	}()
	env.logger.Log("METHOD", "Register", "SPOT", "after NATS event", "time_spent", time.Since(start))
	env.logger.Log("METHOD", "Register", "SPOT", "method end")
	rtn := &pb.Response{Result: common.SUCCESS, Data: responseUserId, Error: nil}
	return rtn, err
}

func (env *Env) Authenticate(ctx context.Context, req *pb.LoginRequest) (*pb.Response, error) {

	env.logger.Log("METHOD", "Authenticate", "SPOT", "method start")
	meta, _ := metadata.FromIncomingContext(ctx)
	apiKey, err := common.GetAPIKey(meta["authorization"][0])
	if err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2001", err.Error())}, nil
	}

	client, err := env.repo.GetClientFromApiKey(apiKey)
	if err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2002", err.Error())}, nil
	}
	_, err = common.ValidateJWTToken(meta["authorization"][0], client.Secret)
	if err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2003", err.Error())}, nil
	}
	clientId := client.Id

	user, err := env.repo.Authenticate(req, clientId)
	if err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2006", err.Error())}, nil
	}

	signedJwtToken := createUserToken(user)

	//NATS Event Publish
	go func() {
		err = env.nats.PublishAuthEvent(user.Id, signedJwtToken)
		if err != nil {
			env.logger.Log("Error during publishing: ", err)
		}
	}()

	tokenStr := map[string]string{"token": signedJwtToken}
	env.logger.Log("METHOD", "Authenticate", "SPOT", "method end")
	return &pb.Response{Result: common.SUCCESS, Data: tokenStr, Error: nil}, nil
}

func validateParameters(user *pb.User) error {
	if err := validation.Validate(
		user.Email,
		validation.Required.Error("email is required"),
		is.Email.Error("email must be a valid email address")); err != nil {
		return err
	}
	if err := validation.Validate(
		user.FirstName, validation.Required,
		validation.Length(1, 32)); err != nil {
		return err
	}
	if err := validation.Validate(
		user.LastName,
		validation.Required,
		validation.Length(1, 32)); err != nil {
		return err
	}
	if err := validation.Validate(
		user.Password,
		validation.Required,
		validation.Length(6, 32)); err != nil {
		return err
	}
	return nil
}

func createUserToken(user *pb.User) string {
	type Claims struct {
		CurrentUserId string `json:"current_user_id"`
		ClientId      string `json:"client_id"`
		CreatedAt     int64  `json:"created_at"`
		jwt.StandardClaims
	}

	claims := Claims{
		user.Id,
		user.ClientId,
		time.Now().UnixNano(),
		jwt.StandardClaims{
			Issuer: "MEEM",
		},
	}
	signingKey := []byte(common.SIGNING_KEY)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJwtToken, _ := token.SignedString(signingKey)
	return signedJwtToken
}
