package usecases

import (
	"context"
	"fmt"
	productrepo "github.com/SanExpett/marketplace-backend/internal/product/repository"
	"github.com/SanExpett/marketplace-backend/pkg/models"
	myerrors "github.com/SanExpett/marketplace-backend/pkg/my_errors"
	"github.com/SanExpett/marketplace-backend/pkg/my_logger"
	"go.uber.org/zap"
	"io"
)

var _ IProductStorage = (*productrepo.ProductStorage)(nil)

type IProductStorage interface {
	AddProduct(ctx context.Context, preProduct *models.PreProduct) (*models.Product, error)
	GetProduct(ctx context.Context, productID uint64, userID uint64) (*models.ProductWithIsMy, error)
	GetProductsList(ctx context.Context, limit uint64, offset uint64, sortType uint64, minPrice uint64,
		maxPrice uint64, userID uint64) ([]*models.ProductWithIsMy, error)
}

type ProductService struct {
	storage IProductStorage
	logger  *zap.SugaredLogger
}

func NewProductService(productStorage IProductStorage) (*ProductService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &ProductService{storage: productStorage, logger: logger}, nil
}

func (p *ProductService) AddProduct(ctx context.Context, r io.Reader, userID uint64) (*models.Product, error) {
	preProduct, err := ValidatePreProduct(r, userID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	product, err := p.storage.AddProduct(ctx, preProduct)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return product, nil
}

func (p *ProductService) GetProduct(ctx context.Context, productID uint64, userID uint64) (*models.ProductWithIsMy, error) {
	product, err := p.storage.GetProduct(ctx, productID, userID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	product.Sanitize()

	return product, nil
}

func (p *ProductService) GetProductsList(ctx context.Context, limit uint64, offset uint64, sortType uint64,
	minPrice uint64, maxPrice uint64, userID uint64,
) ([]*models.ProductWithIsMy, error) {
	products, err := p.storage.GetProductsList(ctx, limit, offset, sortType, minPrice, maxPrice, userID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, product := range products {
		product.Sanitize()
	}

	return products, nil
}
