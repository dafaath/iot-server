package http

import (
	"github.com/gofiber/fiber/v2"
)

func registerRoute(server Server) error {
	server.CreateHealthCheckRoute()
	server.CreateHardwareRoute()
	server.CreateChannelRoute()
	server.CreateNodeRoute()
	server.CreateUserRoute()
	return nil
}

func (r *Server) CreateHealthCheckRoute() {
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

func (r *Server) CreateUserRoute() {
	userRouter := r.app.Group("/user")
	userRouter.Post("/signup", r.userHandler.Register)
	userRouter.Get("/signup", r.userHandler.RegisterPage)
	userRouter.Post("/login", r.userHandler.Login)
	userRouter.Get("/login", r.userHandler.LoginPage)
	userRouter.Post("/forget-password", r.userHandler.ForgotPassword)
	userRouter.Get("/forget-password", r.userHandler.ForgotPasswordPage)
	userRouter.Get("/activation", r.userHandler.Activation)
	userRouter.Get("/", r.authMiddleware.ValidateAdmin, r.userHandler.GetAll)
	userRouter.Get("/:id", r.authMiddleware.ValidateAdmin, r.userHandler.GetOne)
	userRouter.Put("/:id", r.authMiddleware.ValidateUserSameAsUrlIdOrAdmin, r.userHandler.Update)
	userRouter.Delete("/:id", r.authMiddleware.ValidateUserSameAsUrlIdOrAdmin, r.userHandler.Delete)
}

func (r *Server) CreateHardwareRoute() {
	hardwareRouter := r.app.Group("/hardware")
	hardwareRouter.Get("/create", r.authMiddleware.ValidateAdmin, r.hardwareHandler.CreateForm)
	hardwareRouter.Post("/", r.authMiddleware.ValidateAdmin, r.hardwareHandler.Create)
	hardwareRouter.Get("/", r.authMiddleware.ValidateUser, r.hardwareHandler.GetAll)
	hardwareRouter.Get("/:id/edit", r.authMiddleware.ValidateAdmin, r.hardwareHandler.UpdateForm)
	hardwareRouter.Get("/:id", r.authMiddleware.ValidateUser, r.hardwareHandler.GetById)
	hardwareRouter.Put("/:id", r.authMiddleware.ValidateAdmin, r.hardwareHandler.Update)
	hardwareRouter.Delete("/:id", r.authMiddleware.ValidateAdmin, r.hardwareHandler.Delete)
}

func (r *Server) CreateNodeRoute() {
	nodeRouter := r.app.Group("/node")
	nodeRouter.Get("/create", r.authMiddleware.ValidateUser, r.nodeHandler.CreateForm)
	nodeRouter.Post("/", r.authMiddleware.ValidateUser, r.nodeHandler.Create)
	nodeRouter.Get("/", r.authMiddleware.ValidateUser, r.nodeHandler.GetAll)
	nodeRouter.Get("/:id/edit", r.authMiddleware.ValidateUser, r.nodeHandler.UpdateForm)
	nodeRouter.Get("/:id", r.authMiddleware.ValidateUser, r.nodeHandler.GetById)
	nodeRouter.Put("/:id", r.authMiddleware.ValidateUser, r.nodeHandler.Update)
	nodeRouter.Delete("/:id", r.authMiddleware.ValidateUser, r.nodeHandler.Delete)
}
func (r *Server) CreateChannelRoute() {
	channelRouter := r.app.Group("/channel")
	channelRouter.Get("/create", r.authMiddleware.ValidateUser, r.channelHandler.CreateForm)
	channelRouter.Post("/", r.authMiddleware.ValidateUser, r.channelHandler.Create)
}
