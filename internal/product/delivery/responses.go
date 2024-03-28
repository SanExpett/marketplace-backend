package delivery

import (
	"github.com/SanExpett/marketplace-backend/pkg/models"
)

type ProductResponse struct {
	Status int             `json:"status"`
	Body   *models.Product `json:"body"`
}

func NewProductResponse(status int, body *models.Product) *ProductResponse {
	return &ProductResponse{
		Status: status,
		Body:   body,
	}
}

type ProductWithIsMyResponse struct {
	Status int                     `json:"status"`
	Body   *models.ProductWithIsMy `json:"body"`
}

func NewProductWithIsMyResponse(status int, body *models.ProductWithIsMy) *ProductWithIsMyResponse {
	return &ProductWithIsMyResponse{
		Status: status,
		Body:   body,
	}
}

type ProductListResponse struct {
	Status int                       `json:"status"`
	Body   []*models.ProductWithIsMy `json:"body"`
}

func NewProductListResponse(status int, body []*models.ProductWithIsMy) *ProductListResponse {
	return &ProductListResponse{
		Status: status,
		Body:   body,
	}
}
