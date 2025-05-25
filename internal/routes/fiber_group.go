package authsystemroutes

import (
	"github.com/gofiber/fiber/v2"
	authsystem "github.com/viniciusfonseca/auth-system/internal/shared"
)

func AddAuthGroup(a *authsystem.AuthSystem) fiber.Router {

	authGroup := a.FiberApp.Group("/auth")

	authGroup.Post("/sessions", a.GetFiberHandler(CreateSession))
	authGroup.Post("/users", a.GetFiberHandler(CreateUser))

	return authGroup
}
