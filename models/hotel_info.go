package models

type RoomInfo struct {
	Price     int
	Available bool
}

type HotelInfo struct {
	Name  string
	Stars int
	Rooms []RoomInfo
}
