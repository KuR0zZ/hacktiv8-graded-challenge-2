package repository

import (
	"graded-challenge-2-server/model"

	"gorm.io/gorm"
)

type ServerRepository interface {
	CreateUser(user *model.User) error
	GetUser(username string) (*model.User, error)
	GetBookByID(bookID string) (*model.Book, error)
	AddBook(book *model.Book) error
	UpdateBook(book *model.Book) error
	DeleteBook(bookID string) error
}

type ServerRepositoryImpl struct {
	db *gorm.DB
}

func NewServerRepositoryImpl(db *gorm.DB) ServerRepository {
	return &ServerRepositoryImpl{db}
}

func (sr *ServerRepositoryImpl) CreateUser(user *model.User) error {
	return sr.db.Create(user).Error
}

func (sr *ServerRepositoryImpl) GetUser(username string) (*model.User, error) {
	var user model.User
	err := sr.db.Where("username = ?", username).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (sr *ServerRepositoryImpl) GetBookByID(bookID string) (*model.Book, error) {
	var book model.Book
	err := sr.db.Where("id = ?", bookID).Take(&book).Error
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (sr *ServerRepositoryImpl) AddBook(book *model.Book) error {
	return sr.db.Create(book).Error
}

func (sr *ServerRepositoryImpl) UpdateBook(book *model.Book) error {
	return sr.db.Save(book).Error
}

func (sr *ServerRepositoryImpl) DeleteBook(bookID string) error {
	return sr.db.Delete(&model.Book{}, bookID).Error
}
