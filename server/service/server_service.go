package service

import (
	"context"
	"errors"
	"graded-challenge-2-server/model"
	"graded-challenge-2-server/pb"
	"graded-challenge-2-server/repository"
	"math"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type ServerService struct {
	pb.UnimplementedServerServiceServer
	repo repository.ServerRepository
}

func NewServerService(repo repository.ServerRepository) *ServerService {
	return &ServerService{repo: repo}
}

func (ss *ServerService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	existingUser, err := ss.repo.GetUser(req.Username)
	if err == nil && existingUser != nil {
		return nil, status.Error(codes.AlreadyExists, "username already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Password: string(hashPassword),
	}

	err = ss.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	res := &pb.RegisterResponse{
		Id:       user.ID,
		Username: user.Username,
	}

	return res, nil
}

func (ss *ServerService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := ss.repo.GetUser(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "invalid username/password")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid username/password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	res := &pb.LoginResponse{Token: tokenString}

	return res, nil
}

func (ss *ServerService) GetBookById(ctx context.Context, req *pb.BookID) (*pb.GetBookByIDResponse, error) {
	book, err := ss.repo.GetBookByID(req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "book not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &pb.GetBookByIDResponse{
		Id:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		PublishedDate: book.PublishedDate.Format("2006-01-02"),
		Status:        book.Status,
		UserId:        book.UserID,
	}

	return res, nil
}

func (ss *ServerService) AddBook(ctx context.Context, req *pb.AddBookRequest) (*pb.AddBookResponse, error) {
	date, err := time.Parse("2006-01-02", req.PublishedDate)
	if err != nil {
		return nil, err
	}

	book := &model.Book{
		ID:            uuid.New().String(),
		Title:         req.Title,
		Author:        req.Author,
		PublishedDate: date,
		Status:        req.Status,
		UserID:        req.UserId,
	}

	if err := ss.repo.AddBook(book); err != nil {
		return nil, err
	}

	res := &pb.AddBookResponse{
		Id:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		PublishedDate: book.PublishedDate.Format("2006-01-02"),
		Status:        book.Status,
		UserId:        book.UserID,
	}

	return res, nil
}

func (ss *ServerService) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	date, err := time.Parse("2006-01-02", req.PublishedDate)
	if err != nil {
		return nil, err
	}

	book, err := ss.repo.GetBookByID(req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "book not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	book.Title = req.Title
	book.Author = req.Author
	book.PublishedDate = date
	book.Status = req.Status
	book.UserID = req.UserId

	if err := ss.repo.UpdateBook(book); err != nil {
		return nil, err
	}

	res := &pb.UpdateBookResponse{
		Id:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		PublishedDate: book.PublishedDate.Format("2006-01-02"),
		Status:        book.Status,
		UserId:        book.UserID,
	}

	return res, nil
}

func (ss *ServerService) DeleteBook(ctx context.Context, req *pb.BookID) (*emptypb.Empty, error) {
	_, err := ss.repo.GetBookByID(req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "book not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := ss.repo.DeleteBook(req.Id); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (ss *ServerService) BorrowBook(ctx context.Context, req *pb.BorrowBookRequest) (*emptypb.Empty, error) {
	book, err := ss.repo.GetBookByID(req.BookId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "book not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if book.Status == "Borrowed" {
		return nil, status.Error(codes.Unavailable, "book is already borrowed")
	}

	borrowDate, err := time.Parse("2006-01-02", req.BorrowDate)
	if err != nil {
		return nil, err
	}

	returnDate, err := time.Parse("2006-01-02", req.ReturnedDate)
	if err != nil {
		return nil, err
	}

	borrowDetail := &model.BorrowedBook{
		ID:           uuid.New().String(),
		BookID:       req.BookId,
		UserID:       req.UserId,
		BorrowedDate: borrowDate,
		ReturnDate:   returnDate,
	}

	if err := ss.repo.BorrowBook(borrowDetail); err != nil {
		return nil, err
	}

	book.Status = "Borrowed"

	if err := ss.repo.UpdateBook(book); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (ss *ServerService) UpdateBookStatus(ctx context.Context, req *emptypb.Empty) (*pb.UpdateBookStatusResponse, error) {
	books, err := ss.repo.GetAllBook()
	if err != nil {
		return nil, err
	}

	counter := 0
	for _, book := range books {
		if book.Status == "Borrowed" {
			borrowDetail, err := ss.repo.GetBorrowByBookID(book.ID)
			if err != nil {
				return nil, err
			}

			currDate := time.Now()
			dayPass := int(math.Round(currDate.Sub(borrowDetail.BorrowedDate).Hours() / 24))

			zeroValue := time.Time{}

			if borrowDetail.ReturnDate == zeroValue && dayPass > 7 {
				book.Status = "Late"

				err := ss.repo.UpdateBook(&book)
				if err != nil {
					return nil, err
				}

				counter++
			}
		}
	}

	return &pb.UpdateBookStatusResponse{Counter: int32(counter)}, nil
}
