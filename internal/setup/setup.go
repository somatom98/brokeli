package setup

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/somatom98/brokeli/internal/features/create_transactions"
)

type App struct {
	httpHandler *http.ServeMux
	httpServer  *http.Server
}

func Setup() (*App, error) {
	httpHandler := HttpHandler()
	dispatcher := Dispatcher()

	create_transactions.
		New(httpHandler, dispatcher).
		Setup()

	return &App{
		httpHandler: httpHandler,
	}, nil
}

func (a *App) Start() <-chan error {
	port := os.Getenv("PORT")

	a.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: a.httpHandler,
	}

	errCh := make(chan error)

	go func() {
		defer close(errCh)

		log.Printf("Starting server on :%s", port)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("ListenAndServe: %w", err)
			return
		}
	}()

	return errCh
}

func (a *App) Stop(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}
