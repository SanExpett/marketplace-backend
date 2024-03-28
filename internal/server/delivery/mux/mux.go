package mux

import (
	"context"
	"github.com/SanExpett/marketplace-backend/pkg/middleware"
	"net/http"

	productdelivery "github.com/SanExpett/marketplace-backend/internal/product/delivery"
	userdelivery "github.com/SanExpett/marketplace-backend/internal/user/delivery"

	"go.uber.org/zap"
)

type ConfigMux struct {
	addrOrigin string
	schema     string
	portServer string
}

func NewConfigMux(addrOrigin string, schema string, portServer string) *ConfigMux {
	return &ConfigMux{
		addrOrigin: addrOrigin,
		schema:     schema,
		portServer: portServer,
	}
}

func NewMux(ctx context.Context, configMux *ConfigMux, userService userdelivery.IUserService,
	productService productdelivery.IProductService, logger *zap.SugaredLogger,
) (http.Handler, error) {
	router := http.NewServeMux()

	userHandler, err := userdelivery.NewUserHandler(userService)
	if err != nil {
		return nil, err
	}

	productHandler, err := productdelivery.NewProductHandler(productService)
	if err != nil {
		return nil, err
	}

	router.Handle("/api/v1/signup", middleware.Context(ctx,
		middleware.SetupCORS(userHandler.SignUpHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/signin", middleware.Context(ctx,
		middleware.SetupCORS(userHandler.SignInHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/logout", middleware.Context(ctx, http.HandlerFunc(userHandler.LogOutHandler)))

	router.Handle("/api/v1/product/add", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.AddProductHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/product/get", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.GetProductHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/product/get_list", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.GetProductListHandler, configMux.addrOrigin, configMux.schema)))

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(router, logger))

	return mux, nil
}
