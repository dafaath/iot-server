// go:build wireinject
//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

func InitializeApp() (*fiber.App, error) {
	wire.Build(allSet)
	return &fiber.App{}, nil
}
