package delivery

import (
	"context"
	"github.com/SanExpett/marketplace-backend/pkg/models"
	"github.com/SanExpett/marketplace-backend/pkg/my_logger"
	"github.com/SanExpett/marketplace-backend/pkg/utils"
	"io"
	"net/http"
	"time"

	"github.com/SanExpett/marketplace-backend/internal/server/delivery"
	userusecases "github.com/SanExpett/marketplace-backend/internal/user/usecases"
	"github.com/SanExpett/marketplace-backend/pkg/jwt"
	"go.uber.org/zap"
)

const (
	timeTokenLife = 24 * time.Hour

	StatusUnauthorized = 401

	ResponseSuccessfulSignUp = "Successful sign up"
	ResponseSuccessfulSignIn = "Successful sign in"
	ResponseSuccessfulLogOut = "Successful log out"

	ErrUnauthorized = "Вы не авторизованны"
)

var _ IUserService = (*userusecases.UserService)(nil)

type IUserService interface {
	AddUser(ctx context.Context, r io.Reader) (*models.User, error)
	GetUser(ctx context.Context, login string, password string) (*models.UserWithoutPassword, error)
}

type UserHandler struct {
	service IUserService
	logger  *zap.SugaredLogger
}

func NewUserHandler(userService IUserService) (*UserHandler, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	return &UserHandler{
		service: userService,
		logger:  logger,
	}, nil
}

// SignUpHandler godoc
//
//	@Summary    signup
//	@Description  signup in app
//
//	@Description Error.status can be:
//	@Description StatusErrBadRequest      = 400
//	@Description  StatusErrInternalServer  = 500
//	@Tags auth
//
//	@Accept      json
//	@Produce    json
//	@Param      preUser  body models.UserWithoutID true  "user data for signup"
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /signup [post]
func (u *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	user, err := u.service.AddUser(ctx, r.Body)
	if err != nil {
		delivery.HandleErr(w, u.logger, err)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	jwtStr, err := jwt.GenerateJwtToken(
		&jwt.UserJwtPayload{UserID: user.ID, Login: user.Login, Expire: expire.Unix()}, jwt.Secret, u.logger)
	if err != nil {
		delivery.SendErrResponse(w, u.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    delivery.CookieAuthName,
		Value:   jwtStr,
		Expires: expire,
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	delivery.SendOkResponse(w, u.logger, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulSignUp))
	u.logger.Infof("in SignUpHandler: added user: %+v", user)
}

// SignInHandler godoc
//
//	@Summary    signin
//	@Description  signin in app
//	@Tags auth
//	@Produce    json
//	@Param      login  query string true  "user login for signin"
//	@Param      password  query string true  "user password for signin"
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /signin [get]
func (u *UserHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	login := utils.ParseStringFromRequest(r, "login")
	password := utils.ParseStringFromRequest(r, "password")

	user, err := u.service.GetUser(ctx, login, password)
	if err != nil {
		delivery.HandleErr(w, u.logger, err)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	jwtStr, err := jwt.GenerateJwtToken(&jwt.UserJwtPayload{
		UserID: user.ID,
		Login:  user.Login,
		Expire: expire.Unix(),
	},
		jwt.Secret,
		u.logger,
	)
	if err != nil {
		delivery.SendErrResponse(w, u.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    delivery.CookieAuthName,
		Value:   jwtStr,
		Expires: expire,
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	delivery.SendOkResponse(w, u.logger, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulSignIn))
	u.logger.Infof("in SignInHandler: signin user: %+v", user)
}

// LogOutHandler godoc
//
//	@Summary    logout
//	@Description  logout in app
//	@Tags auth
//	@Produce    json
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /logout [post]
func (u *UserHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	cookie, err := r.Cookie(delivery.CookieAuthName)
	if err != nil {
		u.logger.Errorln(err)
		delivery.SendErrResponse(w, u.logger, delivery.NewErrResponse(StatusUnauthorized, ErrUnauthorized))

		return
	}

	cookie.Expires = time.Now()

	http.SetCookie(w, cookie)
	delivery.SendOkResponse(w, u.logger, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulLogOut))
	u.logger.Infof("in LogOutHandler: logout user with cookie: %+v", cookie)
}
