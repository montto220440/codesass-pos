package dto

type CatRequest struct {
	Name string `json:"name" binding:"required"`
}

type CategoryResponse struct {
	ID uint `json:"id"`
	Name string `json:"name"`
}