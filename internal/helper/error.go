package helper

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func HandleStackTrace(e interface{}) {
	_, _ = os.Stderr.WriteString(fmt.Sprintf("panic: %v\n%s\n", e, debug.Stack()))
}

func IsErrorNotFound(err error) bool {
	if err == nil {
		return false
	}

	var e *fiber.Error
	if errors.As(err, &e) && e.Code == 404 {
		return true
	}
	return false
}

func ChangeErrorIfErrorIsNotFound(err error, newError error) error {
	var e *fiber.Error
	if errors.As(err, &e) && e.Code == 404 {
		return newError
	}
	return err
}

func FiberErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	} else {
		log.Printf("[UNHANDLED ERROR] %v", err)
		HandleStackTrace(e)
	}
	// Set Content-Type: text/plain; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	accept := c.Accepts("application/json", "text/html")

	switch accept {
	case "text/html":
		message := err.Error()
		showLogin := false
		if code == 401 || code == 403 {
			showLogin = true
		}
		return c.Render("error", fiber.Map{
			"code":      code,
			"status":    http.StatusText(code),
			"message":   message,
			"showLogin": showLogin,
		}, "layouts/main")
	default:
		// Return status code with error message
		return c.Status(code).SendString(err.Error())
	}

}
