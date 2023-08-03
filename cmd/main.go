package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dafaath/iot-server/v2/configs"
	"github.com/dafaath/iot-server/v2/internal"
	"github.com/dafaath/iot-server/v2/internal/database"
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

	config := configs.GetConfig()
	app, err := internal.InitializeApp()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server is running on", config.Server.Env, "environment:")
	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)))
}
