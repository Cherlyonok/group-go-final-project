package models

type RoomInfo struct {
	Id        int
	HotelId   int
	Price     int
	Available bool
}

type HotelInfo struct {
	Id          int
	OwnerId     int
	Name        string
	Description string
}

type AddHotelJson struct {
	OwnerId     int    `json:"owner_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Rooms       []int  `json:"room_prices"`
}
