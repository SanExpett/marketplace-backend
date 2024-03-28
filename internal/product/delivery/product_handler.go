package delivery

import (
	"context"
	"errors"
	"fmt"
	"github.com/SanExpett/marketplace-backend/internal/product/usecases"
	"github.com/SanExpett/marketplace-backend/internal/server/delivery"
	"github.com/SanExpett/marketplace-backend/pkg/models"
	myerrors "github.com/SanExpett/marketplace-backend/pkg/my_errors"
	"github.com/SanExpett/marketplace-backend/pkg/my_logger"
	"github.com/SanExpett/marketplace-backend/pkg/utils"
	"go.uber.org/zap"
	"io"
	"math"
	"net/http"
)

var _ IProductService = (*usecases.ProductService)(nil)

type IProductService interface {
	AddProduct(ctx context.Context, r io.Reader, userID uint64) (*models.Product, error)
	GetProduct(ctx context.Context, productID uint64, userID uint64) (*models.ProductWithIsMy, error)
	GetProductsList(ctx context.Context, limit uint64, offset uint64, sortType uint64, minPrice uint64,
		maxPrice uint64, userID uint64) ([]*models.ProductWithIsMy, error)
}

type ProductHandler struct {
	service IProductService
	logger  *zap.SugaredLogger
}

func NewProductHandler(productService IProductService) (*ProductHandler, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &ProductHandler{
		service: productService,
		logger:  logger,
	}, nil
}

// AddProductHandler godoc
//
//	@Summary    add product
//	@Description  add product by data
//	@Description Error.status can be:
//	@Description StatusErrBadRequest      = 400
//	@Description  StatusErrInternalServer  = 500
//	@Tags product
//
//	@Accept      json
//	@Produce    json
//	@Param      product  body models.PreProduct true  "product data for adding"
//	@Success    200  {object} ProductResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/add [post]
func (p *ProductHandler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	product, err := p.service.AddProduct(ctx, r.Body, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, NewProductResponse(delivery.StatusResponseSuccessful, product))
	p.logger.Infof("in AddProductHandler: add product: %+v", product)
}

// GetProductHandler godoc
//
//	@Summary    get product
//	@Description  get product by id
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      id  query uint64 true  "product id"
//	@Success    200  {object} ProductResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/get [get]
func (p *ProductHandler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		if errors.Is(err, delivery.ErrCookieNotPresented) {
			userID = 0
		} else {
			delivery.HandleErr(w, p.logger, err)

			return
		}
	}

	productID, err := utils.ParseUint64FromRequest(r, "id")
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	product, err := p.service.GetProduct(ctx, productID, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, NewProductWithIsMyResponse(delivery.StatusResponseSuccessful, product))
	p.logger.Infof("in GetProductHandler: get product: %+v", product)
}

// GetProductsListHandler godoc
//
//	@Summary    get Products list
//	@Description  get Products by count and last_id return old Products
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      limit  query uint64 true  "limit Products"
//	@Param      offset  query uint64 true  "offset of Products"
//	@Param      min_price  query uint64 true  "min price of product"
//	@Param      max_price  query uint64 true  "max price of product"
//	@Param      sort_type query uint64 true  "type of sort(nil - by date desc, 1 - by price asc, 2 - by price desc, 3 - by date asc, 4 - by date desc)"
//	@Success    200  {object} ProductListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /Product/get_list [get]
func (p *ProductHandler) GetProductListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		if errors.Is(err, delivery.ErrCookieNotPresented) {
			userID = 0
		} else {
			delivery.HandleErr(w, p.logger, err)

			return
		}
	}

	limit, err := utils.ParseUint64FromRequest(r, "limit")
	if err != nil {
		limit = 10
	}

	offset, err := utils.ParseUint64FromRequest(r, "offset")
	if err != nil {
		offset = 0
	}

	sortType, err := utils.ParseUint64FromRequest(r, "sort_type")
	if err != nil {
		sortType = 0
	}

	minPrice, err := utils.ParseUint64FromRequest(r, "min_price")
	if err != nil {
		minPrice = 0
	}

	maxPrice, err := utils.ParseUint64FromRequest(r, "max_price")
	if err != nil {
		maxPrice = math.MaxUint64
	}

	products, err := p.service.GetProductsList(ctx, limit, offset, sortType, minPrice, maxPrice, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, NewProductListResponse(delivery.StatusResponseSuccessful, products))
	p.logger.Infof("in GetProductListHandler: get Product list: %+v", products)
}
