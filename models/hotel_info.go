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

type AddHotelJson struct {
	OwnerId  int    `json:"owner_id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Rooms    []int  `json:"room_prices"`
}
