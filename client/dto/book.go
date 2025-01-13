package dto

type AddBookRequest struct {
	Title         string `json:"title" validate:"required"`
	Author        string `json:"author" validate:"required"`
	PublishedDate string `json:"published_date" validate:"required"`
	Status        string `json:"status" validate:"required"`
}

type UpdateBookRequest struct {
	Title         string `json:"title" validate:"required"`
	Author        string `json:"author" validate:"required"`
	PublishedDate string `json:"published_date" validate:"required"`
	Status        string `json:"status" validate:"required"`
}
