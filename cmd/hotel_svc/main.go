package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"project/internal/hotel_svc"
	"project/internal/hotel_svc/hotel_grpc"
	"project/internal/hotel_svc/hotel_rest"
	"project/pkg/server"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func startHttpServer(wg *sync.WaitGroup, port string, hotel_db *hotel_svc.HotelDB) {
	defer wg.Done()

	hotel_server := server.CreateServer(":" + port)
	hotel_service := hotel_rest.HotelService{HotelDB: hotel_db}

	requests := []server.Request{
		{Handler: hotel_service.AddHotel, Path: "/hotels/add"},
		{Handler: hotel_service.FindHotels, Path: "/hotels/find"},
		{Handler: hotel_service.GetAvailableHotelRooms, Path: "/hotels/available_rooms"},
	}
	for _, val := range requests {
		hotel_server.AddRequest(val)
	}
	hotel_server.Start()
}

func startGrpcServer(wg *sync.WaitGroup, port string, hotel_db *hotel_svc.HotelDB) {
	defer wg.Done()

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	hotel_grpc.RegisterHotelServiceServer(s, &hotel_grpc.GrpcServer{HotelDB: hotel_db})

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		<-c
		cancel()
	}()

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Printf("grpc server listening at %v", lis.Addr())
		return s.Serve(lis)
	})
	g.Go(func() error {
		<-gCtx.Done()
		s.GracefulStop()
		fmt.Printf("grpc: Server closed. \n")
		return nil
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("GRPC Server exit reason: %s \n", err)
	}
}

func main() {
	hotel_service, _ := hotel_svc.CreateHotelService("postgres://postgres:12345@localhost:5555?sslmode=disable")
	var wg sync.WaitGroup

	wg.Add(2)

	go startHttpServer(&wg, "8080", &hotel_service)
	go startGrpcServer(&wg, "8081", &hotel_service)

	wg.Wait()
	fmt.Println("Program executed.")
}
