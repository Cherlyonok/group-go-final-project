package models

type BookingInfo struct {
	Id        int
	ClientId  int
	RoomId    int
	HotelId   int
	StartDate string
	EndDate   string
}
