package controller

import (
	"graded-challenge-2-client/dto"
	"graded-challenge-2-client/helper"
	"graded-challenge-2-client/pb"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerController struct {
	client pb.ServerServiceClient
}

func NewServerController(client pb.ServerServiceClient) *ServerController {
	return &ServerController{client}
}

func (sc *ServerController) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(helper.ErrBadRequest.ErrorFormat(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	data := &pb.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	}

	ctx, cancel, err := helper.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	resGrpc, err := sc.client.Register(ctx, data)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unauthenticated:
				return echo.NewHTTPError(helper.ErrUnauthorized.ErrorFormat(e.Message()))
			case codes.AlreadyExists:
				return echo.NewHTTPError(helper.ErrBadRequest.ErrorFormat(e.Message()))
			case codes.Internal:
				return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(e.Message()))
			}
		}

		return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(err.Error()))
	}

	res := dto.WebSuccessResponse{
		Status:  http.StatusCreated,
		Message: "User created successfully",
		Data:    resGrpc,
	}

	return c.JSON(http.StatusCreated, res)
}

func (sc *ServerController) Login(c echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(helper.ErrBadRequest.ErrorFormat(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	data := &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	ctx, cancel, err := helper.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	resGrpc, err := sc.client.Login(ctx, data)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound, codes.Unauthenticated:
				return echo.NewHTTPError(helper.ErrUnauthorized.ErrorFormat(e.Message()))
			case codes.Internal:
				return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(e.Message()))
			}
		}
		return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(err.Error()))
	}

	res := dto.WebSuccessResponse{
		Status:  http.StatusOK,
		Message: "Successfully login",
		Data:    resGrpc,
	}

	return c.JSON(http.StatusOK, res)
}
