package main

import (
	"context"
	"github.com/arpanetus/counter/pkg/file"
	"github.com/arpanetus/counter/pkg/router"
	"github.com/arpanetus/counter/pkg/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	Minute     = time.Minute
	TimeFormat = time.RFC3339
	Path       = "/tmp/countertimer"
)

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	svc := service.New(Path, file.New(Minute, TimeFormat))
	if err := svc.Parse(); err != nil {
		log.Fatalf("cannot parse times on init: %v", err)
	}

	handler := router.New(svc)
	mux := http.NewServeMux()
	mux.Handle("/", handler)

	srv := &http.Server{
		Addr:    ":2021",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err!=http.ErrServerClosed {
			if closeErr := svc.Close(); closeErr != nil {
				log.Printf("cannot close counting file: %v", closeErr)
			}
			log.Fatalf("cannot listen: %+v", err)
		}
	}()

	log.Printf("started listening on: %s", srv.Addr)

	<- ctx.Done()

	log.Printf("stopped listening on: %s", srv.Addr)

	if closeErr := svc.Close(); closeErr != nil {
		log.Printf("cannot close counting file: %v", closeErr)
	}

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err:= srv.Shutdown(ctxShutDown); err!=nil {
		log.Fatalf("cannot shutdown server: %+v", err)
	}

	log.Printf("server exited properly")

	return
}
