package controller

import (
	"fmt"
	"graded-challenge-2-client/dto"
	"graded-challenge-2-client/helper"
	"graded-challenge-2-client/pb"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (sc *ServerController) GetBookByID(c echo.Context) error {
	bookID := c.Param("id")

	ctx, cancel, err := helper.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	book, err := sc.client.GetBookById(ctx, &pb.BookID{Id: bookID})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return echo.NewHTTPError(helper.ErrUnauthorized.ErrorFormat(e.Message()))
			case codes.Internal:
				return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(e.Message()))
			}
		}
		return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(err.Error()))
	}

	res := dto.WebSuccessResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Successfully retrieve book with ID: %s", bookID),
		Data:    book,
	}

	return c.JSON(http.StatusOK, res)
}

func (sc *ServerController) AddBook(c echo.Context) error {
	var req dto.AddBookRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(helper.ErrBadRequest.ErrorFormat(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	claims, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat("internal server error"))
	}

	userID := claims["user_id"].(string)

	data := &pb.AddBookRequest{
		Title:         req.Title,
		Author:        req.Author,
		PublishedDate: req.PublishedDate,
		Status:        req.Status,
		UserId:        userID,
	}

	ctx, cancel, err := helper.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	book, err := sc.client.AddBook(ctx, data)
	if err != nil {
		return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(err.Error()))
	}

	res := dto.WebSuccessResponse{
		Status:  http.StatusCreated,
		Message: "Successfully add new book",
		Data:    book,
	}

	return c.JSON(http.StatusCreated, res)
}

func (sc *ServerController) UpdateBook(c echo.Context) error {
	bookID := c.Param("id")

	ctx, cancel, err := helper.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	book, err := sc.client.GetBookById(ctx, &pb.BookID{Id: bookID})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return echo.NewHTTPError(helper.ErrUnauthorized.ErrorFormat(e.Message()))
			case codes.Internal:
				return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(e.Message()))
			}
		}
		return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(err.Error()))
	}

	claims, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat("internal server error"))
	}

	userID := claims["user_id"].(string)

	if book.UserId != userID {
		return echo.NewHTTPError(helper.ErrUnauthorized.ErrorFormat("Not allowed"))
	}

	var req dto.UpdateBookRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(helper.ErrBadRequest.ErrorFormat(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	data := &pb.UpdateBookRequest{
		Id:            book.Id,
		Title:         req.Title,
		Author:        req.Author,
		PublishedDate: req.PublishedDate,
		Status:        req.Status,
		UserId:        book.UserId,
	}

	ctx, cancel, err = helper.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	updatedBook, err := sc.client.UpdateBook(ctx, data)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return echo.NewHTTPError(helper.ErrUnauthorized.ErrorFormat(e.Message()))
			case codes.Internal:
				return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(e.Message()))
			}
		}
		return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(err.Error()))
	}

	res := dto.WebSuccessResponse{
		Status:  http.StatusOK,
		Message: "Successfully updated book",
		Data:    updatedBook,
	}

	return c.JSON(http.StatusOK, res)
}

func (sc *ServerController) DeleteBook(c echo.Context) error {
	bookID := c.Param("id")

	ctx, cancel, err := helper.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	book, err := sc.client.GetBookById(ctx, &pb.BookID{Id: bookID})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return echo.NewHTTPError(helper.ErrUnauthorized.ErrorFormat(e.Message()))
			case codes.Internal:
				return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(e.Message()))
			}
		}
		return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(err.Error()))
	}

	if book.Status == "Borrowed" {
		return echo.NewHTTPError(helper.ErrUnprocessable.ErrorFormat("book is still borrowed"))
	}

	claims, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat("internal server error"))
	}

	userID := claims["user_id"].(string)

	if book.UserId != userID {
		return echo.NewHTTPError(helper.ErrUnauthorized.ErrorFormat("Not allowed"))
	}

	ctx, cancel, err = helper.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	_, err = sc.client.DeleteBook(ctx, &pb.BookID{Id: bookID})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return echo.NewHTTPError(helper.ErrUnauthorized.ErrorFormat(e.Message()))
			case codes.Internal:
				return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(e.Message()))
			}
		}
		return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(err.Error()))
	}

	res := dto.WebSuccessResponse{
		Status:  http.StatusOK,
		Message: "Successfully deleted book",
	}

	return c.JSON(http.StatusOK, res)
}

func (sc *ServerController) BorrowBook(c echo.Context) error {
	bookID := c.Param("id")

	claims, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat("internal server error"))
	}

	userID := claims["user_id"].(string)

	data := &pb.BorrowBookRequest{
		BookId:       bookID,
		UserId:       userID,
		BorrowDate:   time.Now().Format("2006-01-02"),
		ReturnedDate: time.Time{}.Format("2006-01-02"),
	}

	ctx, cancel, err := helper.NewServiceContext()
	if err != nil {
		return err
	}
	defer cancel()

	_, err = sc.client.BorrowBook(ctx, data)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return echo.NewHTTPError(helper.ErrUnauthorized.ErrorFormat(e.Message()))
			case codes.Unavailable:
				return echo.NewHTTPError(helper.ErrUnprocessable.ErrorFormat(e.Message()))
			case codes.Internal:
				return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(e.Message()))
			}
		}
		return echo.NewHTTPError(helper.ErrInternalServer.ErrorFormat(err.Error()))
	}

	res := dto.WebSuccessResponse{
		Status:  http.StatusCreated,
		Message: fmt.Sprintf("Successfully borrowed book with ID: %s", bookID),
	}

	return c.JSON(http.StatusCreated, res)
}
