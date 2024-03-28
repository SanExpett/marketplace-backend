package server

import (
	"context"
	productrepo "github.com/SanExpett/marketplace-backend/internal/product/repository"
	productusecases "github.com/SanExpett/marketplace-backend/internal/product/usecases"
	"github.com/SanExpett/marketplace-backend/internal/server/delivery/mux"
	"github.com/SanExpett/marketplace-backend/internal/server/repository"
	userrepo "github.com/SanExpett/marketplace-backend/internal/user/repository"
	userusecases "github.com/SanExpett/marketplace-backend/internal/user/usecases"
	"github.com/SanExpett/marketplace-backend/pkg/config"
	"github.com/SanExpett/marketplace-backend/pkg/my_logger"
	"net/http"
	"strings"
	"time"
)

const (
	basicTimeout = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(config *config.Config) error {
	baseCtx := context.Background()

	pool, err := repository.NewPgxPool(baseCtx, config.URLDataBase)
	if err != nil {
		return err //nolint:wrapcheck
	}

	logger, err := my_logger.New(strings.Split(config.OutputLogPath, " "),
		strings.Split(config.ErrorOutputLogPath, " "))
	if err != nil {
		return err //nolint:wrapcheck
	}

	defer logger.Sync()

	userStorage, err := userrepo.NewUserStorage(pool)
	if err != nil {
		return err
	}

	userService, err := userusecases.NewUserService(userStorage)
	if err != nil {
		return err
	}

	productStorage, err := productrepo.NewProductStorage(pool)
	if err != nil {
		return err
	}

	productService, err := productusecases.NewProductService(productStorage)
	if err != nil {
		return err
	}

	handler, err := mux.NewMux(baseCtx, mux.NewConfigMux(config.AllowOrigin,
		config.Schema, config.PortServer), userService, productService, logger)
	if err != nil {
		return err
	}

	s.httpServer = &http.Server{ //nolint:exhaustruct
		Addr:           ":" + config.PortServer,
		Handler:        handler,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ReadTimeout:    basicTimeout,
		WriteTimeout:   basicTimeout,
	}

	logger.Infof("Start server:%s", config.PortServer)

	return s.httpServer.ListenAndServe() //nolint:wrapcheck
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx) //nolint:wrapcheck
}
