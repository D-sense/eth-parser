package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trustwallet/api/server"
	"trustwallet/business/logger"
	"trustwallet/business/parser"
	"trustwallet/business/storage"
)

// config is used to represent runtime configuration.
type config struct {
	ethereumGatewayURL string
}

// cfg provides parsed runtime configuration as a convenient global variable.
var cfg config

func init() {
	cfg.ethereumGatewayURL = os.Getenv("ETHEREUM_GATEWAY_URL")
	if cfg.ethereumGatewayURL == "" {
		// I created this project on Infura (https://app.infura.io/dashboard) to help speed up testing the service
		cfg.ethereumGatewayURL = "https://mainnet.infura.io/v3/3b7ef887e2b244b9b0bd9b2a0c36cdf1"
	}
}

func main() {
	// Construct the application logger.
	log, err := logger.New("API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err = run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	// =========================================================================
	// Start API Service

	log.Infow("startup", "status", "initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Initialize Ethereum Parser
	ethereumParser := parser.NewEthereumParser(storage.NewMemoryStorage(), cfg.ethereumGatewayURL, 5, log)

	// Construct the mux for the API calls.
	apiMux := server.APIMux(server.APIMuxConfig{
		Ctx:      context.Background(),
		Shutdown: shutdown,
		Log:      log,
		Parser:   ethereumParser,
	})

	// Construct a server to service the requests against the mux.
	api := http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      apiMux,
		ReadTimeout:  time.Second * 5, //cfg.Web.ReadTimeout,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 120,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for server requests.
	go func() {
		log.Infow("startup", "status", "server router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
