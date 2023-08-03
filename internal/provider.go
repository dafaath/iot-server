package internal

import (
	"github.com/dafaath/iot-server/v2/configs"
	"github.com/dafaath/iot-server/v2/internal/database"
	"github.com/dafaath/iot-server/v2/internal/dependencies"
	"github.com/dafaath/iot-server/v2/internal/handlers"
	"github.com/dafaath/iot-server/v2/internal/middlewares"
	"github.com/dafaath/iot-server/v2/internal/repositories"
	"github.com/dafaath/iot-server/v2/internal/server/http"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

var (
	dependenciesSet = wire.NewSet(
		database.GetConnection,
		configs.GetConfig,
		validator.New,
		dependencies.NewValidator,
		dependencies.NewMailDialer,
	)
	middlewaresSet = wire.NewSet(
		middlewares.NewAuthenticationMiddleware,
	)
	repositorySet = wire.NewSet(
		repositories.NewUserRepository,
		repositories.NewHardwareRepository,
		repositories.NewNodeRepository,
		repositories.NewChannelRepository,
	)
	handlerSet = wire.NewSet(
		handlers.NewUserHandler,
		handlers.NewHardwareHandler,
		handlers.NewNodeHandler,
		handlers.NewChannelHandler,
	)

	allSet = wire.NewSet(
		dependenciesSet,
		repositorySet,
		middlewaresSet,
		handlerSet,
		http.NewHTTPServer,
	)
)
