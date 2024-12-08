package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	address  string
	mux      *http.ServeMux
	requests []Request
}

func CreateServer(address string) Server {
	return Server{address: address, mux: http.NewServeMux(), requests: make([]Request, 0)}
}

func (server *Server) AddRequest(request Request) {
	server.requests = append(server.requests, request)
}

func (server *Server) Start() {
	for _, val := range server.requests {
		server.mux.HandleFunc(val.Path, val.Handler)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		<-c
		cancel()
	}()

	httpServer := http.Server{
		Addr:    server.address,
		Handler: server.mux,
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Printf("http server listening at %v", httpServer.Addr)
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("%s \n", err)
	}
}
