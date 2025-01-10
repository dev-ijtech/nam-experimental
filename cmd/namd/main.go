package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/dev-ijtech/nam-experimental/namhttp"
	"github.com/dev-ijtech/nam-experimental/namsql"
	"github.com/dev-ijtech/nam-experimental/southbound"

	_ "github.com/mattn/go-sqlite3"
)

func run(ctx context.Context, stdout io.Writer, stderr io.Writer, getenv func(string) string) error {
	var err error

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := log.New(stdout, "namd: ", log.LstdFlags)

	db, err := sql.Open("sqlite3", "./test.db")

	if err != nil {
		return err
	}

	defer db.Close()

	southboundServiceImpl := southbound.NewSouthboundService("root", "clab123", logger)
	deviceStore := namsql.NewDeviceService(db)

	serverPort := getenv("PORT")

	if serverPort == "" {
		serverPort = "8080"
	}

	srv := namhttp.NewServer(logger, deviceStore, southboundServiceImpl)

	httpServer := http.Server{
		Addr:    "127.0.0.1:" + serverPort,
		Handler: srv,
	}

	go func() {
		logger.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()

	return err
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Stderr, os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	fmt.Println("server successfully shutdown")
}
