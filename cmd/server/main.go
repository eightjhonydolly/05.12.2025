package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eightjhonydolly/05.12.2025/internal/app"
)

func main() {
	fmt.Println("Link checker service starting...")

	app, err := app.NewApp(os.Getenv("CONFIG_ENV_VAR"))
	if err != nil {
		log.Fatal("Failed to create app:", err)
	}

	if err := app.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed:", err)
	}

	fmt.Println("Server stopped gracefully")
}
