package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/SanExpett/marketplace-backend/internal/server/repository"
	"github.com/SanExpett/marketplace-backend/pkg/models"
	myerrors "github.com/SanExpett/marketplace-backend/pkg/my_errors"
	"github.com/SanExpett/marketplace-backend/pkg/my_logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"math"
	"time"
)

var (
	ErrProductNotFound = myerrors.NewError("Этот товар не найден")

	NameSeqProduct = pgx.Identifier{"public", "product_id_seq"} //nolint:gochecknoglobals
)

type ProductStorage struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

const (
	byPriceASC  = 1
	byPriceDESC = 2
	byDateASC   = 3
	byDateDESC  = 4
)

func NewProductStorage(pool *pgxpool.Pool) (*ProductStorage, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &ProductStorage{
		pool:   pool,
		logger: logger,
	}, nil
}

func (p *ProductStorage) insertProduct(ctx context.Context, tx pgx.Tx, preProduct *models.PreProduct) error {
	SQLInsertProduct := `INSERT INTO public."product"(saler_id,
		title, description, price, image_url) VALUES(
		$1, $2, $3, $4, $5)`
	_, err := tx.Exec(ctx, SQLInsertProduct, preProduct.SalerID,
		preProduct.Title, preProduct.Description, preProduct.Price, preProduct.ImageUrl)

	if err != nil {
		p.logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) AddProduct(ctx context.Context, preProduct *models.PreProduct) (*models.Product, error) {
	product := &models.Product{Title: preProduct.Title, Description: preProduct.Description,
		Price: preProduct.Price, SalerID: preProduct.SalerID, ImageUrl: preProduct.ImageUrl}

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.insertProduct(ctx, tx, preProduct)
		if err != nil {
			return err
		}

		lastProductID, err := repository.GetLastValSeq(ctx, tx, NameSeqProduct)
		if err != nil {
			return err
		}

		createdAt, err := p.selectCreatedAtByProductID(ctx, tx, lastProductID)
		if err != nil {
			return err
		}

		product.CreatedAt = createdAt

		return err
	})
	if err != nil {
		p.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return product, nil
}

func (p *ProductStorage) selectProductByID(ctx context.Context, tx pgx.Tx, productID uint64, userID uint64,
) (*models.ProductWithIsMy, error) {
	SQLSelectProduct := `SELECT saler_id, image_url, title,
       description, price, created_at FROM public."product" WHERE id=$1`
	product := &models.ProductWithIsMy{ID: productID} //nolint:exhaustruct

	productRow := tx.QueryRow(ctx, SQLSelectProduct, productID)
	if err := productRow.Scan(&product.SalerID, &product.ImageUrl,
		&product.Title, &product.Description, &product.Price, &product.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(myerrors.ErrTemplate, ErrProductNotFound)
		}

		p.logger.Errorf("error with productId=%d: %+v", productID, err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if product.SalerID == userID {
		product.IsMy = true
	} else {
		product.IsMy = false
	}

	return product, nil
}

func (p *ProductStorage) selectCreatedAtByProductID(ctx context.Context, tx pgx.Tx, productID uint64,
) (time.Time, error) {
	SQLSelectCreatedAtByProductID := `SELECT created_at FROM public."product" WHERE id=$1`

	var createdAt time.Time

	createdAtRow := tx.QueryRow(ctx, SQLSelectCreatedAtByProductID, productID)
	if err := createdAtRow.Scan(&createdAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return time.Time{}, fmt.Errorf(myerrors.ErrTemplate, ErrProductNotFound)
		}

		p.logger.Errorf("error with productId=%d: %+v", productID, err)

		return time.Time{}, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return createdAt, nil
}

func (p *ProductStorage) GetProduct(ctx context.Context, productID uint64, userID uint64) (*models.ProductWithIsMy, error) {
	var product *models.ProductWithIsMy

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		productInner, err := p.selectProductByID(ctx, tx, productID, userID)
		if err != nil {
			return err
		}

		product = productInner

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return product, nil
}

func (p *ProductStorage) selectProductsWithWhereOrderLimitOffset(ctx context.Context, tx pgx.Tx,
	limit uint64, offset uint64, whereClause any, orderByClause []string, userID uint64,
) ([]*models.ProductWithIsMy, error) {
	query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("id, saler_id, title," +
		"description, price, created_at, image_url").From(`public."product"`).
		Where(whereClause).OrderBy(orderByClause...).Limit(limit).Offset(offset)

	SQLQuery, args, err := query.ToSql()
	if err != nil {
		p.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsProducts, err := tx.Query(ctx, SQLQuery, args...)
	if err != nil {
		p.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curProduct := new(models.ProductWithIsMy)

	var slProduct []*models.ProductWithIsMy

	_, err = pgx.ForEachRow(rowsProducts, []any{
		&curProduct.ID, &curProduct.SalerID, &curProduct.Title, &curProduct.Description,
		&curProduct.Price, &curProduct.CreatedAt, &curProduct.ImageUrl,
	}, func() error {
		slProduct = append(slProduct, &models.ProductWithIsMy{ //nolint:exhaustruct
			ID:          curProduct.ID,
			SalerID:     curProduct.SalerID,
			Title:       curProduct.Title,
			Description: curProduct.Description,
			Price:       curProduct.Price,
			CreatedAt:   curProduct.CreatedAt,
			ImageUrl:    curProduct.ImageUrl,
		})

		return nil
	})
	if err != nil {
		p.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, product := range slProduct {
		if product.SalerID == userID {
			product.IsMy = true
		} else {
			product.IsMy = false
		}
	}

	return slProduct, nil
}

func (p *ProductStorage) GetProductsList(ctx context.Context,
	limit uint64, offset uint64, sortType uint64, minPrice uint64, maxPrice uint64, userID uint64,
) ([]*models.ProductWithIsMy, error) {
	var slProduct []*models.ProductWithIsMy

	var orderByClause []string

	switch sortType {
	case byPriceASC:
		orderByClause = []string{"price ASC"}
	case byPriceDESC:
		orderByClause = []string{"price DESC"}
	case byDateASC:
		orderByClause = []string{"created_at ASC"}
	case byDateDESC:
		orderByClause = []string{"created_at DESC"}
	default:
		orderByClause = []string{"created_at DESC"}
	}

	whereClause := ""
	if !(minPrice == 0 && maxPrice == math.MaxUint64) && (minPrice <= maxPrice) {
		whereClause = fmt.Sprintf("price >= %d AND price <= %d", minPrice, maxPrice)
	}

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		var err error
		slProduct, err = p.selectProductsWithWhereOrderLimitOffset(ctx,
			tx, limit, offset, whereClause, orderByClause, userID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		p.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slProduct, nil
}
