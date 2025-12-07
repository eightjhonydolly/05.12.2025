package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eightjhonydolly/05.12.2025/internal/app/handlers/check_links_handler"
	"github.com/eightjhonydolly/05.12.2025/internal/app/handlers/generate_report_handler"
	"github.com/eightjhonydolly/05.12.2025/internal/domain/links/repository"
	"github.com/eightjhonydolly/05.12.2025/internal/domain/links/service"
	"github.com/eightjhonydolly/05.12.2025/internal/infra/config"
	"github.com/eightjhonydolly/05.12.2025/internal/infra/http/middlewares"
)

type App struct {
	config *config.Config
	server http.Server
}

func NewApp(configPath string) (*App, error) {
	configImpl, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{
		config: configImpl,
	}

	app.server.Handler = bootstrapHandler()

	return app, nil
}

func (app *App) ListenAndServe() error {
	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)

	log.Printf("Starting server on %s", address)
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("Failed to listen on %s: %v", address, err)
		return err
	}

	go app.gracefulShutdown()

	log.Printf("Server listening on %s", address)
	return app.server.Serve(l)
}

func (app *App) gracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Received shutdown signal, starting graceful shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		log.Printf("Error during graceful shutdown: %v", err)
	} else {
		log.Println("Server shutdown completed")
	}
}

func bootstrapHandler() http.Handler {
	linkRepository := repository.NewInMemoryLinkRepository()
	linkService := service.NewLinkService(linkRepository)

	mx := http.NewServeMux()
	mx.Handle("POST /api/check-links", check_links_handler.NewCheckLinksHandler(linkService))
	mx.Handle("POST /api/generate-report", generate_report_handler.NewGenerateReportHandler(linkService))

	middleware := middlewares.NewTimerMiddleware(mx)

	return middleware
}
