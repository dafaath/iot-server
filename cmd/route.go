package main

import (
	"github.com/dafaath/iot-server/v2/internal/handlers"
	"github.com/dafaath/iot-server/v2/internal/middlewares"
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
		return c.SendString("Server OK")
	})
}

func (r *Router) CreateUserRoute(handler *handlers.UserHandler) {
	userRouter := r.app.Group("/user")
	userRouter.Post("/signup", handler.Register)
	userRouter.Post("/login", handler.Login)
	userRouter.Post("/forget-password", handler.ForgotPassword)
	userRouter.Get("/activation", handler.Activation)
	userRouter.Get("/", r.authMiddleware.ValidateAdmin, handler.GetAll)
	userRouter.Get("/:id", r.authMiddleware.ValidateAdmin, handler.GetOne)
	userRouter.Put("/:id", r.authMiddleware.ValidateUserSameAsUrlIdOrAdmin, handler.Update)
	userRouter.Delete("/:id", r.authMiddleware.ValidateUserSameAsUrlIdOrAdmin, handler.Delete)
}

func (r *Router) CreateHardwareRoute(handler *handlers.HardwareHandler) {
	hardwareRouter := r.app.Group("/hardware")
	hardwareRouter.Post("/", r.authMiddleware.ValidateAdmin, handler.Create)
	hardwareRouter.Get("/", r.authMiddleware.ValidateUser, handler.GetAll)
	hardwareRouter.Get("/:id", r.authMiddleware.ValidateUser, handler.GetById)
	hardwareRouter.Put("/:id", r.authMiddleware.ValidateAdmin, handler.Update)
	hardwareRouter.Delete("/:id", r.authMiddleware.ValidateAdmin, handler.Delete)
}

func (r *Router) CreateNodeRoute(handler *handlers.NodeHandler) {
	nodeRouter := r.app.Group("/node")
	nodeRouter.Post("/", r.authMiddleware.ValidateUser, handler.Create)
	nodeRouter.Get("/", r.authMiddleware.ValidateUser, handler.GetAll)
	nodeRouter.Get("/:id", r.authMiddleware.ValidateUser, handler.GetById)
	nodeRouter.Put("/:id", r.authMiddleware.ValidateUser, handler.Update)
	nodeRouter.Delete("/:id", r.authMiddleware.ValidateUser, handler.Delete)
}
func (r *Router) CreateChannelRoute(handler *handlers.ChannelHandler) {
	channelRouter := r.app.Group("/channel")
	channelRouter.Post("/", r.authMiddleware.ValidateUser, handler.Create)
}
