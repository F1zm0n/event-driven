package transport

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/F1zm0n/event-driven/auth/internal/dto"
	"github.com/F1zm0n/event-driven/auth/internal/usecases"
	event_pb_authv1 "github.com/F1zm0n/event-driven/auth/protos/gen/auth/v1"
)

type Customer struct {
	event_pb_authv1.UnimplementedAuthServiceServer
	service usecases.CustomerUsecases
}

func (c *Customer) Login(
	ctx context.Context,
	req *event_pb_authv1.LoginRequest,
) (*event_pb_authv1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}
	customerID, err := dto.ParseUUID(req.CustomerId)
	if err != nil {
		return nil, err
	}
	customerDto := dto.CustomerDto{
		CustomerID: customerID,
		Email:      req.Email,
		Password:   req.Password,
	}
	token, err := c.service.LoginJWT(ctx, customerDto)
	if err != nil {
		return nil, err
	}

	return &event_pb_authv1.LoginResponse{
		Token: token,
	}, err
}

func (c *Customer) CustomerByID(
	ctx context.Context,
	req *event_pb_authv1.UserByIDRequest,
) (*event_pb_authv1.UserByIDResponse, error) {
	if err := validateGetByID(req); err != nil {
		return nil, err
	}
	customerID, err := dto.ParseUUID(req.CustomerId)
	if err != nil {
		return nil, err
	}
	customer, err := c.service.GetByID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	return &event_pb_authv1.UserByIDResponse{
		CustomerId: req.CustomerId,
		Email:      customer.Email,
	}, nil
}

func (c *Customer) CustomerByEmail(
	ctx context.Context,
	req *event_pb_authv1.UserByEmailRequest,
) (*event_pb_authv1.UserByEmailResponse, error) {
	if err := validateGetByEmail(req); err != nil {
		return nil, err
	}
	customer, err := c.service.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &event_pb_authv1.UserByEmailResponse{
		CustomerId: customer.CustomerID.String(),
		Email:      customer.Email,
	}, nil
}

func RegisterServer(gRPC *grpc.Server, service usecases.CustomerUsecases) {
	event_pb_authv1.RegisterAuthServiceServer(gRPC, &Customer{service: service})
}

func validateLogin(req *event_pb_authv1.LoginRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email field shouldn't be empty")
	}
	if req.Password == "" {
		return fmt.Errorf("password field shouldn't be empty")
	}
	if req.CustomerId == "" {
		return fmt.Errorf("customer_id field shouldn't be empty")
	}
	return nil
}

func validateGetByID(req *event_pb_authv1.UserByIDRequest) error {
	if req.CustomerId == "" {
		return fmt.Errorf("customer_id field shouldn't be empty")
	}
	return nil
}

func validateGetByEmail(req *event_pb_authv1.UserByEmailRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email field shouldn't be empty")
	}
	return nil
}
