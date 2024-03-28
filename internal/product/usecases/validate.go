package usecases

import (
	"encoding/json"
	"fmt"
	"github.com/SanExpett/marketplace-backend/pkg/models"
	myerrors "github.com/SanExpett/marketplace-backend/pkg/my_errors"
	"github.com/SanExpett/marketplace-backend/pkg/my_logger"
	"github.com/asaskevich/govalidator"
	"io"
)

var (
	ErrDecodePreProduct = myerrors.NewError("Некорректный json объявления")
)

func validatePreProduct(r io.Reader, userID uint64) (*models.PreProduct, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(r)
	preProduct := &models.PreProduct{}
	if err := decoder.Decode(preProduct); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodePreProduct)
	}

	preProduct.Trim()

	preProduct.SalerID = userID

	_, err = govalidator.ValidateStruct(preProduct)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return preProduct, nil
}

func ValidatePreProduct(r io.Reader, userID uint64) (*models.PreProduct, error) {
	preProduct, err := validatePreProduct(r, userID)
	if err != nil {
		return nil, myerrors.NewError(err.Error())
	}

	return preProduct, nil
}
