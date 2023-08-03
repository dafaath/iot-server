package http

import (
	"encoding/json"

	"github.com/dafaath/iot-server/v2/internal/handlers"
	"github.com/dafaath/iot-server/v2/internal/helper"
	"github.com/dafaath/iot-server/v2/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/handlebars"
)

type Server struct {
	app             *fiber.App
	authMiddleware  *middlewares.AuthenticationMiddleware
	userHandler     *handlers.UserHandler
	hardwareHandler *handlers.HardwareHandler
	nodeHandler     *handlers.NodeHandler
	channelHandler  *handlers.ChannelHandler
}

func NewHTTPServer(
	authMiddleware *middlewares.AuthenticationMiddleware,
	userHandler *handlers.UserHandler, hardwareHandler *handlers.HardwareHandler,
	nodeHandler *handlers.NodeHandler, channelHandler *handlers.ChannelHandler,
) (*fiber.App, error) {
	engine := handlebars.New("./internal/views", ".hbs")

	app := fiber.New(
		fiber.Config{
			// Override default error handler
			Views:        engine,
			ErrorHandler: helper.FiberErrorHandler,
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		},
	)
	app.Static("/static", "./internal/public")
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	server := Server{
		app:             app,
		authMiddleware:  authMiddleware,
		userHandler:     userHandler,
		hardwareHandler: hardwareHandler,
		nodeHandler:     nodeHandler,
		channelHandler:  channelHandler,
	}
	registerRoute(server)

	return app, nil
}
