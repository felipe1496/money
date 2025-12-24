package categories

import (
	"rango-backend/utils"
	"time"
)

// ==============================================================================
// 1. HTTP MODELS
//    Models that represents request or response objects
// ==============================================================================

type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

type CreateCategoryResponse struct {
	Data CreateCategoryResponseData `json:"data"`
}

type CreateCategoryResponseData struct {
	Category Category `json:"category"`
}

type ListCategoriesResponse struct {
	Data  ListCategoriesResponseData `json:"data"`
	Query utils.QueryMeta            `json:"query"`
}

type ListCategoriesResponseData struct {
	Categories []Category `json:"categories"`
}

type ListCategoryAmountPerPeriodResponse struct {
	Data ListCategoryAmountPerPeriodResponseData `json:"data"`
}

type ListCategoryAmountPerPeriodResponseData struct {
	Categories []CategoryAmountPerPeriod `json:"categories"`
}

// ==============================================================================
// 2. DTO MODELS
//    Models that represents data transfer objects between api layers
// ==============================================================================

type CreateCategoryDTO struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

// ==============================================================================
// 3. DATABASE
//    Models that represents database objects
// ==============================================================================

type Category struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryAmountPerPeriod struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Period      string  `json:"period"`
	TotalAmount float64 `json:"total_amount"`
}
