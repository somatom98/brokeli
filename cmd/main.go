package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/somatom98/brokeli/internal/setup"
)

func main() {
	app, err := setup.Setup()
	if err != nil {
		log.Fatalf("Setup: %v", err)
	}

	errCh := app.Start()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case err := <-errCh:
		// Server crashed or failed to bind
		if err != nil {
			log.Fatalf("server error: %v", err)
		}
		log.Println("server stopped cleanly")

	case sig := <-sigCh:
		log.Printf("caught %s, stopping server…", sig)

		// Give outstanding requests up to 5 s to finish.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := app.Stop(ctx); err != nil {
			log.Fatalf("server stop failed: %v", err)
		}
		log.Println("server stop complete")
	}
}
