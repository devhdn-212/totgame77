package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on real environment variables")
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	model := os.Getenv("ANTHROPIC_MODEL")
	if model == "" {
		model = "claude-sonnet-5"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	ocr := newOcrHandler(apiKey, model)

	app := fiber.New()

	api := app.Group("/api")
	api.Post("/ocr/scan", ocr.Scan)

	// Serve the built Svelte SPA (run `yarn build` first) and fall back to
	// index.html so client-side routes resolve on refresh.
	app.Static("/", "./dist")
	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile("./dist/index.html")
	})

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
	_ = app.Shutdown()
}
