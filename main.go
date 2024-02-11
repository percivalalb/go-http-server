package main

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const (
	shutdownTimeout = 10 * time.Second

	requestTimeout = 15 * time.Second

	port = 8080
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	log.Info("Starting")

	t1 := time.Now()

	log.Info("Listening for interrupt signals")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Info("Setting up mux")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", exampleHandler(log))

	server := http.Server{
		Handler: mux,
		// Protect against attacks. By default there is no timeout.
		ReadHeaderTimeout: requestTimeout,
		ReadTimeout:       requestTimeout,
		WriteTimeout:      requestTimeout,
	}

	// Listen on all interfaces on the given port.
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	log.Info("Listening on", slog.String("addr", listener.Addr().String()))

	serveErr := make(chan error, 1)

	// Serve requests in a seperate go-routine and then watch ctx & errors below.
	go func() {
		defer close(serveErr)

		log.Info("Serving", slog.Duration("time", time.Since(t1)))
		serveErr <- server.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		t1 := time.Now()

		log.Info("Initiating shutdown", slog.Duration("timeout", shutdownTimeout))

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// Wait for all HTTP requests to finish
		if err := server.Shutdown(ctx); err != nil {
			panic(err)
		}

		log.Info("Shutdown", slog.Duration("time", time.Since(t1)))
	case err := <-serveErr:
		if err != nil {
			panic(err)
		}
	}

	// Serve immediately returns ErrServerClosed when Shutdown is called.
	// Check there hasn't been an an error between ctx finishing and Shutdown
	// being called.
	if err := <-serveErr; err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func exampleHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("serving request", slog.String("source", r.RemoteAddr))
		io.WriteString(w, "Hello, world!\n")
	}
}
