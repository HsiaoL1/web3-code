package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hsiaocz/web3-code/book/services"
)

type UserHandler struct {
	UserService services.UserServiceInterface
}

func (h *UserHandler) RegisterRoutes(c *fiber.Ctx) {

}
