package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"mylittleprice/internal/app"
	"mylittleprice/internal/config"
	"mylittleprice/internal/container"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	c, err := container.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer c.Close()

	fiberApp := fiber.New(fiber.Config{
		AppName:      "MyLittlePrice API",
		ServerHeader: "Fiber",
		ErrorHandler: customErrorHandler,
	})

	fiberApp.Use(recover.New())
	fiberApp.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "15:04:05",
	}))

	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.CORSOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	app.SetupRoutes(fiberApp, c)

	port := cfg.Port
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🔒 Environment: %s", cfg.Env)
	log.Printf("🌍 Allowed origins: %v", cfg.CORSOrigins)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("🛑 Shutting down server...")

		if err := fiberApp.Shutdown(); err != nil {
			log.Printf("❌ Server shutdown error: %v", err)
		}

		log.Println("✅ Server stopped gracefully")
	}()

	if err := fiberApp.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

func customErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return ctx.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": message,
		"code":    code,
	})
}
