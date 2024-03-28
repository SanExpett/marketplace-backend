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
	ErrWrongCredentials = myerrors.NewError("Некорректный логин (должен быть длиной от 1 до 25 " +
		"символов) или пароль (должен быть не менее 6 символов, содержать цифры, " +
		"строчные и заглавные буквы и специальные символы)")
	ErrDecodeUser = myerrors.NewError("Некорректный json пользователя")
)

func ValidateUserWithoutID(r io.Reader) (*models.UserWithoutID, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	decoder := json.NewDecoder(r)

	userWithoutID := new(models.UserWithoutID)
	if err := decoder.Decode(userWithoutID); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodeUser)
	}

	userWithoutID.Trim()

	_, err = govalidator.ValidateStruct(userWithoutID)
	if err != nil {
		return nil, ErrWrongCredentials
	}

	return userWithoutID, nil
}

func ValidateUserCredentials(login string, password string) (*models.UserWithoutID, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	userWithoutID := new(models.UserWithoutID)

	userWithoutID.Login = login
	userWithoutID.Password = password
	userWithoutID.Trim()
	logger.Infoln(userWithoutID)

	_, err = govalidator.ValidateStruct(userWithoutID)
	if err != nil && (govalidator.ErrorByField(err, "login") != "" ||
		govalidator.ErrorByField(err, "password") != "") {
		logger.Errorln(err)

		return nil, ErrWrongCredentials
	}

	return userWithoutID, nil
}
