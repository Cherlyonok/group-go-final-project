package main

import (
	"fmt"
	"project/internal/booking_svc"
	"project/internal/booking_svc/booking_rest"
	"project/pkg/server"
	"sync"
)

func startHttpServer(wg *sync.WaitGroup, port string, bookingDb *booking_svc.BookingDB, hotelServiceAddr string) {
	defer wg.Done()

	bookingServer := server.CreateServer(":" + port)
	bookingService := booking_rest.NewBookingService(bookingDb, hotelServiceAddr)

	requests := []server.Request{
		{Handler: bookingService.AddBooking, Path: "/bookings/add"},
		{Handler: bookingService.GetBookingsByClientID, Path: "/bookings/get_by_client_id"},
		{Handler: bookingService.GetBookingsByHotelID, Path: "/bookings/get_by_hotel_id"},
	}
	for _, val := range requests {
		bookingServer.AddRequest(val)
	}
	bookingServer.Start()
}

func main() {
	bookingService, _ := booking_svc.CreateBookingService("postgres://postgres:12345@localhost:5555?sslmode=disable")
	var wg sync.WaitGroup

	wg.Add(2)

	go startHttpServer(&wg, "8082", &bookingService, "localhost:8081")

	wg.Wait()
	fmt.Println("Program executed.")
}
