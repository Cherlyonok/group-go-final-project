package main

import (
	"project/internal/hotel_svc"
	"project/pkg/server"

	_ "github.com/lib/pq"
)

func main() {
	hotel_server := server.CreateServer(":8080")
	hotel_service, _ := hotel_svc.CreateHotelService("postgres://postgres:12345@localhost:5555?sslmode=disable")

	requests := []server.Request{
		{Handler: hotel_service.AddHotel, Path: "/hotels/add"},
		{Handler: hotel_service.FindHotels, Path: "/hotels/find"},
		{Handler: hotel_service.GetRoomPrice, Path: "/hotels/get_room_price"},
	}
	for _, val := range requests {
		hotel_server.AddRequest(val)
	}
	hotel_server.Start()
}
