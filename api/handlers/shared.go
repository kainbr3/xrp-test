package handlers

import (
	"crypto-braza-tokens-api/api/handlers/types"
	l "crypto-braza-tokens-api/utils/logger"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func UnauthorizedWrapper(ctx *fiber.Ctx, message string) error {
	l.Logger.Error("handler: unauthorized acces", zap.String("error", message))

	return ctx.Status(fiber.StatusUnauthorized).JSON(types.ErrorMessage{Message: message})
}

func BadRequestWrapper(ctx *fiber.Ctx, resource string, err error) error {
	msg := fmt.Sprintf("handler: error saving %s", resource)

	l.Logger.Error(msg, zap.Error(err))

	formattedError := fmt.Sprintf("%s with error: %v", msg, err)

	return ctx.Status(http.StatusBadRequest).JSON(types.ErrorMessage{Message: formattedError})
}

func NotFoundWrapper(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusNotFound).JSON("{}")
}

func ConflictWrapper(ctx *fiber.Ctx, resource string) error {
	msg := fmt.Sprintf("handler: resource %s already exists", resource)

	l.Logger.Error(msg)

	return ctx.Status(http.StatusConflict).JSON(types.ErrorMessage{Message: msg})
}

func InternalErrorWrapper(ctx *fiber.Ctx, resource string, err error) error {
	msg := fmt.Sprintf("handler: error saving %s", resource)

	l.Logger.Error(msg, zap.Error(err))

	formattedError := fmt.Sprintf("%s with error: %v", msg, err)

	return ctx.Status(http.StatusInternalServerError).JSON(types.ErrorMessage{Message: formattedError})
}

func MessageResultWrapper(ctx *fiber.Ctx, result string) error {
	return ctx.Status(http.StatusOK).JSON(types.Result{Result: result})
}

func ObjectResultWrapper(ctx *fiber.Ctx, result any) error {
	return ctx.Status(http.StatusOK).JSON(result)
}

func ValidatePathParam(ctx *fiber.Ctx, param string) error {
	if ctx.Params(param) == "" || ctx.Params(param) == "undefined" {
		l.Logger.Error(fmt.Sprintf("handler: error validating %s", param))
		return fmt.Errorf("%s is required", param)
	}

	return nil
}

// DefaultPath root path validation to redirect to default route path (swagger)
func DefaultPath(ctx *fiber.Ctx) error {
	ctx.Redirect("/api/docs")

	return nil
}

// AuthorizerHandler - middleware to validate the Authorization header
func AuthorizerHandler(ctx *fiber.Ctx) error {
	if ctx.Get("x-webhook-secret") != "" {
		l.Logger.Info("auth: webhook", zap.String("path", ctx.Get("x-webhook-secret")))
	}

	if ctx.Get("Authorization") != "" {
		l.Logger.Info("auth: authorizer", zap.String("Authorization", ctx.Get("Authorization")))
	}

	// authHeader := ctx.Get("Authorization")
	// if authHeader == "" {
	// 	return UnauthorizedWrapper(ctx, "Missing Authorization header")
	// }

	// Check if the Authorization header starts with "Bearer "
	// if !strings.HasPrefix(authHeader, "Bearer ") {
	// 	UnauthorizedWrapper(ctx, "Invalid Authorization header format. Must start with 'Bearer '")
	// }

	// Extract the token
	// token := strings.TrimPrefix(authHeader, "Bearer ")
	// if token == "" {
	// 	UnauthorizedWrapper(ctx, "Missing token")
	// }

	// Optionally, you can add more checks for the token here (e.g., validate the token)

	// Continue to the next middleware/handler
	return ctx.Next()
}
