package main

import (
	"github.com/dafaath/iot-server/internal/handlers"
	"github.com/dafaath/iot-server/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	app            *fiber.App
	authMiddleware *middlewares.AuthenticationMiddleware
}

func NewRouter(app *fiber.App, authMiddleware *middlewares.AuthenticationMiddleware) (Router, error) {
	return Router{
		app:            app,
		authMiddleware: authMiddleware,
	}, nil
}

func (r *Router) CreateHealthCheckRoute() {
	r.app.Get("/", func(c *fiber.Ctx) error {
		accept := c.Accepts("application/json", "text/html")

		switch accept {
		case "text/html":
			return c.Render("index", fiber.Map{}, "layouts/main")
		default:
			return c.SendString("Server OK")
		}
	})
}

func (r *Router) CreateUserRoute(handler *handlers.UserHandler) {
	userRouter := r.app.Group("/user")
	userRouter.Post("/signup", handler.Register)
	userRouter.Get("/signup", handler.RegisterPage)
	userRouter.Post("/login", handler.Login)
	userRouter.Get("/login", handler.LoginPage)
	userRouter.Post("/forget-password", handler.ForgotPassword)
	userRouter.Get("/forget-password", handler.ForgotPasswordPage)
	userRouter.Get("/activation", handler.Activation)
	userRouter.Get("/", r.authMiddleware.ValidateAdmin, handler.GetAll)
	userRouter.Get("/:id", r.authMiddleware.ValidateAdmin, handler.GetOne)
	userRouter.Put("/:id", r.authMiddleware.ValidateUserSameAsUrlIdOrAdmin, handler.Update)
	userRouter.Delete("/:id", r.authMiddleware.ValidateUserSameAsUrlIdOrAdmin, handler.Delete)
}

func (r *Router) CreateHardwareRoute(handler *handlers.HardwareHandler) {
	hardwareRouter := r.app.Group("/hardware")
	hardwareRouter.Get("/create", r.authMiddleware.ValidateUser, handler.CreateForm)
	hardwareRouter.Post("/", r.authMiddleware.ValidateUser, handler.Create)
	hardwareRouter.Get("/", r.authMiddleware.ValidateUser, handler.GetAll)
	hardwareRouter.Get("/:id/edit", r.authMiddleware.ValidateUser, handler.UpdateForm)
	hardwareRouter.Get("/:id", r.authMiddleware.ValidateUser, handler.GetById)
	hardwareRouter.Put("/:id", r.authMiddleware.ValidateUser, handler.Update)
	hardwareRouter.Delete("/:id", r.authMiddleware.ValidateUser, handler.Delete)
}

func (r *Router) CreateNodeRoute(handler *handlers.NodeHandler) {
	nodeRouter := r.app.Group("/node")
	nodeRouter.Get("/create", r.authMiddleware.ValidateUser, handler.CreateForm)
	nodeRouter.Post("/", r.authMiddleware.ValidateUser, handler.Create)
	nodeRouter.Get("/", r.authMiddleware.ValidateUser, handler.GetAll)
	nodeRouter.Get("/:id/edit", r.authMiddleware.ValidateUser, handler.UpdateForm)
	nodeRouter.Get("/:id", r.authMiddleware.ValidateUser, handler.GetById)
	nodeRouter.Put("/:id", r.authMiddleware.ValidateUser, handler.Update)
	nodeRouter.Delete("/:id", r.authMiddleware.ValidateUser, handler.Delete)
}

func (r *Router) CreateSensorRoute(handler *handlers.SensorHandler) {
	sensorRouter := r.app.Group("/sensor")
	sensorRouter.Get("/create", r.authMiddleware.ValidateUser, handler.CreateForm)
	sensorRouter.Post("/", r.authMiddleware.ValidateUser, handler.Create)
	sensorRouter.Get("/", r.authMiddleware.ValidateUser, handler.GetAll)
	sensorRouter.Get("/:id/edit", r.authMiddleware.ValidateUser, handler.UpdateForm)
	sensorRouter.Get("/:id", r.authMiddleware.ValidateUser, handler.GetById)
	sensorRouter.Put("/:id", r.authMiddleware.ValidateUser, handler.Update)
	sensorRouter.Delete("/:id", r.authMiddleware.ValidateUser, handler.Delete)
}

func (r *Router) CreateChannelRoute(handler *handlers.ChannelHandler) {
	channelRouter := r.app.Group("/channel")
	channelRouter.Get("/create", r.authMiddleware.ValidateUser, handler.CreateForm)
	channelRouter.Post("/", r.authMiddleware.ValidateUser, handler.Create)
}
