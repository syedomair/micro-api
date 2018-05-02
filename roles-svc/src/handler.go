package main

import (
	"time"

	log "github.com/go-kit/kit/log"
	"golang.org/x/net/context"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/syedomair/micro-api/common"
	pb "github.com/syedomair/micro-api/roles-svc/proto"
)

type Env struct {
	repo   Repository
	nats   Nats
	logger log.Logger
}

func (env *Env) CreateRole(ctx context.Context, req *pb.Role) (*pb.Response, error) {

	start := time.Now()
	env.logger.Log("METHOD", "Create", "SPOT", "method start", "time_start", start)
	clientId, _ := ctx.Value("client_id").(string)
	env.logger.Log("METHOD", "Create", "SPOT", "method start", "client_id", clientId)

	if err := validateCreateParameters(req); err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2004", err.Error())}, nil
	}

	roleId, err := env.repo.Create(req, clientId)
	if err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.DatabaseError()}, nil
	}
	responseRoleId := map[string]string{"role_id": roleId}
	env.logger.Log("METHOD", "Create", "SPOT", "method end", "time_spent", time.Since(start))
	return &pb.Response{Result: common.SUCCESS, Data: responseRoleId, Error: nil}, err
}

func (env *Env) GetAll(ctx context.Context, req *pb.RequestQuery) (*pb.ResponseList, error) {

	start := time.Now()
	env.logger.Log("METHOD", "GetAll", "SPOT", "method start", "time_start", start)
	clientId, _ := ctx.Value("client_id").(string)

	limit, offset, orderby, sort, err := common.ValidateQueryString(req.Limit, "3", req.Offset, "0", req.Orderby, "title", req.Sort, "asc")
	if err != nil {
		return &pb.ResponseList{Result: common.FAILURE, Error: common.CommonError(err.Error()), Data: nil}, nil
	}

	roles, count, err := env.repo.GetAll(limit, offset, orderby, sort, clientId)
	if err != nil {
		return &pb.ResponseList{Result: common.FAILURE, Error: common.CommonError(err.Error()), Data: nil}, nil
	}

	roleList := &pb.RoleList{Offset: offset, Limit: limit, Count: count, List: roles}
	env.logger.Log("METHOD", "GetAll", "SPOT", "method end", "time_spent", time.Since(start))
	return &pb.ResponseList{Result: common.SUCCESS, Error: nil, Data: roleList}, nil
}

func (env *Env) GetRole(ctx context.Context, req *pb.Role) (*pb.ResponseRole, error) {

	start := time.Now()
	env.logger.Log("METHOD", "GetRole", "SPOT", "method start", "time_start", start)
	clientId, _ := ctx.Value("client_id").(string)

	if err := validateRoleId(req); err != nil {
		return &pb.ResponseRole{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2004", err.Error())}, nil
	}

	role, err := env.repo.Get(req.Id, clientId)
	if err != nil {
		return &pb.ResponseRole{Result: common.FAILURE, Data: nil, Error: common.CommonError(err.Error())}, nil
	}
	env.logger.Log("METHOD", "GetRole", "SPOT", "method end", "time_spent", time.Since(start))
	return &pb.ResponseRole{Result: common.SUCCESS, Data: role, Error: nil}, nil
}

func (env *Env) UpdateRole(ctx context.Context, req *pb.Role) (*pb.Response, error) {

	start := time.Now()
	env.logger.Log("METHOD", "UpdateRole", "SPOT", "method start", "time_start", start)
	env.logger.Log("METHOD", "UpdateRole", "SPOT", "input request:", "req:", req)
	clientId, _ := ctx.Value("client_id").(string)

	if err := validateRoleId(req); err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2004", err.Error())}, nil
	}

	if err := validateUpdateParameters(req); err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2004", err.Error())}, nil
	}

	if err := env.repo.Update(req, clientId); err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.CommonError(err.Error())}, nil
	}
	responseRoleId := map[string]string{"role_id": req.Id}
	env.logger.Log("METHOD", "UpdateRole", "SPOT", "method end", "time_spent", time.Since(start))
	return &pb.Response{Result: common.SUCCESS, Data: responseRoleId, Error: nil}, nil
}

func (env *Env) DeleteRole(ctx context.Context, req *pb.Role) (*pb.Response, error) {

	start := time.Now()
	env.logger.Log("METHOD", "DeleteRole", "SPOT", "method start", "time_start", start)
	clientId, _ := ctx.Value("client_id").(string)

	if err := validateRoleId(req); err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.ErrorMessage("2004", err.Error())}, nil
	}
	if err := env.repo.Delete(req, clientId); err != nil {
		return &pb.Response{Result: common.FAILURE, Data: nil, Error: common.CommonError(err.Error())}, nil
	}
	responseRoleId := map[string]string{"role_id": req.Id}
	env.logger.Log("METHOD", "DeleteRole", "SPOT", "method end", "time_spent", time.Since(start))
	return &pb.Response{Result: common.SUCCESS, Data: responseRoleId, Error: nil}, nil
}

func validateUpdateParameters(role *pb.Role) error {
	if role.Title != "" {
		if err := validation.Validate(
			role.Title,
			validation.Required.Error("title is a required field."),
			validation.Length(1, 64).Error("title is a rqquired field with the max character of 64")); err != nil {
			return err
		}
	}
	if role.RoleType != "" {
		if err := validation.Validate(
			role.RoleType,
			validation.Required.Error("role_type is a required field."),
			is.Digit.Error("role_type must be a digit between 0 and 9")); err != nil {
			return err
		}
	}
	return nil
}

func validateCreateParameters(role *pb.Role) error {
	if err := validation.Validate(
		role.Title,
		validation.Required.Error("title is a required field."),
		validation.Length(1, 64).Error("title is a rqquired field with the max character of 64")); err != nil {
		return err
	}
	if err := validation.Validate(
		role.RoleType,
		validation.Required.Error("role_type is a required field."),
		is.Digit.Error("role_type must be a digit between 0 and 9")); err != nil {
		return err
	}
	return nil
}
func validateRoleId(role *pb.Role) error {
	if err := validation.Validate(
		role.Id,
		validation.Required.Error("role_id is a required field"),
		is.UUIDv4.Error("invalid role_id.")); err != nil {
		return err
	}
	return nil
}
