package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iammarkps/tu83-announcer-api/app"
)

func main() {
	e, db := app.New()

	go func() {
		if err := e.Start(":1323"); err != nil {
			e.Logger.Info(err)
			e.Logger.Fatal("Shutting down HTTP server 🔥🔥🔥")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	} else {
		log.Print("HTTP server gracefully shutdown")
	}

	if err := db.Close(); err != nil {
		e.Logger.Fatal(err)
	} else {
		log.Print("Successfully disconnected from database")
	}
}
