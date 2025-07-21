package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Run initializes and starts the URL shortener application.
// It performs the following steps:
//  1. Loads configuration using config.NewConfig()
//  2. Initializes the storage backend using store.NewStore()
//  3. Creates the router with all endpoints using router.New()
//  4. Starts the HTTP server on the configured address
//
// The function blocks until the server exits and returns any error that occurs.
// In case of a fatal server error, it logs the error and terminates the program.
//
// Example usage:
//
//	if err := app.Run(); err != nil {
//	    log.Fatal("Application failed:", err)
//	}
func Run() error {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Printf("Received signal: %v\n", sig)

	return nil
}
