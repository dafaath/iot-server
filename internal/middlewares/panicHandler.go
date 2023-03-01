package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// New creates a new middleware handler
func NewPanicHandlerMiddleware() fiber.Handler {
	// Return new handler
	return func(c *fiber.Ctx) (err error) {
		// Catch panics
		defer func() {
			if r := recover(); r != nil {
				// if cfg.EnableStackTrace {
				// 	cfg.StackTraceHandler(c, r)
				// }

				var ok bool
				if err, ok = r.(error); !ok {
					// Set error that will call the global error handler
					err = fmt.Errorf("%v", r)
				}
			}
		}()

		// Return err if exist, else move to next handler
		if err != nil {
			return err
		} else {
			return c.Next()
		}
	}
}
