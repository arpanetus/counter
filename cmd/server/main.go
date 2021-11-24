package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arpanetus/counter/pkg/file"
	"github.com/arpanetus/counter/pkg/router"
	"github.com/arpanetus/counter/pkg/service"
)

const (
	Minute = time.Minute
	Path   = "/tmp/countertimer"
)

func CounterMux(path string) (*http.ServeMux, service.CounterServicer) {
	svc := service.New(file.New(path, Minute, file.FS))
	if err := svc.Parse(); err != nil {
		log.Fatalf("cannot parse times on init: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/count", router.New(svc))
	return mux, svc
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	mux, svc :=  CounterMux(Path)

	srv := &http.Server{
		Addr:    ":2021",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if closeErr := svc.Close(); closeErr != nil {
				log.Printf("cannot close counting file: %v", closeErr)
			}
			log.Fatalf("cannot listen: %+v", err)
		}
	}()

	log.Printf("started listening on: %s", srv.Addr)

	<-ctx.Done()

	log.Printf("stopped listening on: %s", srv.Addr)

	if closeErr := svc.Close(); closeErr != nil {
		log.Printf("cannot close counting file: %v", closeErr)
	}

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("cannot shutdown server: %+v", err)
	}

	log.Printf("server exited properly")
}
