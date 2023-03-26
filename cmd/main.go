package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dafaath/iot-server/v2/configs"
	"github.com/dafaath/iot-server/v2/internal/database"
	"github.com/dafaath/iot-server/v2/internal/dependencies"
	"github.com/dafaath/iot-server/v2/internal/handlers"
	"github.com/dafaath/iot-server/v2/internal/helper"
	"github.com/dafaath/iot-server/v2/internal/middlewares"
	"github.com/dafaath/iot-server/v2/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/handlebars"
)

var createDatabaseMode bool

func init() {
	flag.BoolVar(&createDatabaseMode, "create-db", false, "If set to true, this will drop the current table, create the table and create initial user. Then exit program")
}

// Declare all dependencies and run server
func main() {
	// Parse flag
	flag.Parse()

	if createDatabaseMode {
		database.DropTable()
		database.CreateTableAndMockData()
		os.Exit(0)
	}

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

	// BEGIN Other dependencies declaration
	config := configs.GetConfig()
	validate := validator.New()
	db, err := database.GetConnection()
	helper.PanicIfError(err)
	myValidator := dependencies.NewValidator(validate)
	dialer, err := dependencies.NewMailDialer(config)
	helper.PanicIfError(err)
	// END

	// BEGIN Middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	authenticationMiddleware := middlewares.NewAuthenticationMiddleware(&myValidator)
	// END

	// BEGIN Repositories declaration
	userRepository, err := repositories.NewUserRepository(dialer)
	helper.PanicIfError(err)
	hardwareRepository, err := repositories.NewHardwareRepository()
	helper.PanicIfError(err)
	nodeRepository, err := repositories.NewNodeRepository()
	helper.PanicIfError(err)
	helper.PanicIfError(err)
	channelRepository, err := repositories.NewChannelRepository()
	helper.PanicIfError(err)
	// END

	// BEGIN Handlers declaration
	userHandler, err := handlers.NewUserHandler(db, &userRepository, &myValidator)
	helper.PanicIfError(err)
	hardwareHandler, err := handlers.NewHardwareHandler(db, &hardwareRepository, &nodeRepository, &myValidator)
	helper.PanicIfError(err)
	nodeHandler, err := handlers.NewNodeHandler(db, &nodeRepository, &channelRepository, &hardwareRepository, &myValidator)
	helper.PanicIfError(err)
	helper.PanicIfError(err)
	channelHandler, err := handlers.NewChannelHandler(db, &channelRepository, &nodeRepository, &myValidator)
	helper.PanicIfError(err)
	// END

	// BEGIN Routes declaration
	router, err := NewRouter(app, &authenticationMiddleware)
	helper.PanicIfError(err)
	router.CreateHealthCheckRoute()
	router.CreateUserRoute(&userHandler)
	router.CreateHardwareRoute(&hardwareHandler)
	router.CreateNodeRoute(&nodeHandler)
	router.CreateChannelRoute(&channelHandler)
	// END

	// Middleware
	app.Use(cache.New())
	app.Use(etag.New())

	// Initialize default config

	log.Println("Server is running on", config.Server.Env, "environment:")
	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)))
}
