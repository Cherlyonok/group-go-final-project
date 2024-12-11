package booking_rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/internal/booking_svc"
	"project/models"
	"strconv"
)

type BookingService struct {
	BookingDB *booking_svc.BookingDB
}

func (service *BookingService) GetBookingsByClientID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clientIDString := r.URL.Query().Get("client_id")
	if clientIDString == "" {
		http.Error(w, "client_id is required", http.StatusBadRequest)
		return
	}

	clientID, err := strconv.Atoi(clientIDString)
	if err != nil {
		http.Error(w, "invalid client_id", http.StatusBadRequest)
		return
	}

	var bookings []models.BookingInfo

	query := `SELECT id, client_id, room_id, start_date, end_date FROM Bookings WHERE client_id = $1`

	rows, err := service.BookingDB.Db.Query(query, clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}
	defer rows.Close()

	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(&booking.ID, &booking.ClientID, &booking.RoomID, &booking.StartDate, &booking.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err != nil {
		http.Error(w, "failed to retrieve bookings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(bookings)
}

func (service *BookingService) AddHotel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var hotel_data models.AddHotelJson
	err := json.NewDecoder(r.Body).Decode(&hotel_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	add_hotel_query := "INSERT INTO Hotels (owner_id, name, description) VALUES ($1, $2, $3) RETURNING id"
	add_room_query := "INSERT INTO Rooms (hotel_id, price, available) VALUES ($1, $2, $3)"
	var id int
	err = service.HotelDB.Db.QueryRow(add_hotel_query, hotel_data.OwnerId, hotel_data.Name, hotel_data.Description).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	for _, room_price := range hotel_data.Rooms {
		service.HotelDB.Db.Exec(add_room_query, id, room_price, true)
	}
	w.Write([]byte("Succesfully added hotel"))
}

func (service *HotelService) FindHotels(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var hotels []models.HotelInfo
	rows, err := service.HotelDB.Db.Query("SELECT id, owner_id, name, description FROM Hotels")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}
	for rows.Next() {
		var hotel models.HotelInfo
		if err := rows.Scan(&hotel.Id, &hotel.OwnerId, &hotel.Name, &hotel.Description); err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}
		hotels = append(hotels, hotel)
	}
	json.NewEncoder(w).Encode(hotels)
}

func (service *HotelService) GetAvailableHotelRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request_json models.GetAvailableHotelRoomsJson
	err := json.NewDecoder(r.Body).Decode(&request_json)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var rooms []models.RoomInfo
	rows, err := service.HotelDB.Db.Query("SELECT id, hotel_id, price, available FROM Rooms WHERE hotel_id = $1 AND available = TRUE", request_json.HotelId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}
	for rows.Next() {
		var room models.RoomInfo
		if err := rows.Scan(&room.Id, &room.HotelId, &room.Price, &room.Available); err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}
		rooms = append(rooms, room)
	}
	json.NewEncoder(w).Encode(rooms)
}
