package booking_rest

import (
	"encoding/json"
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

	query := `SELECT id, client_id, room_id, hotel_id, start_date, end_date FROM Bookings WHERE client_id = $1`

	rows, err := service.BookingDB.Db.Query(query, clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}
	defer rows.Close()

	for rows.Next() {
		var booking models.BookingInfo

		if err := rows.Scan(&booking.Id, &booking.ClientId, &booking.RoomId, &booking.HotelId, &booking.StartDate, &booking.EndDate); err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}
		bookings = append(bookings, booking)
	}
	json.NewEncoder(w).Encode(bookings)
}

func (service *BookingService) GetBookingsByHotelID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hotelIDString := r.URL.Query().Get("hotel_id")
	if hotelIDString == "" {
		http.Error(w, "hotel_id is required", http.StatusBadRequest)
		return
	}

	hotelID, err := strconv.Atoi(hotelIDString)
	if err != nil {
		http.Error(w, "invalid hotel_id", http.StatusBadRequest)
		return
	}

	var bookings []models.BookingInfo

	query := `SELECT id, client_id, room_id, hotel_id, start_date, end_date FROM Bookings WHERE hotel_id = $1`

	rows, err := service.BookingDB.Db.Query(query, hotelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}
	defer rows.Close()

	for rows.Next() {
		var booking models.BookingInfo

		if err := rows.Scan(&booking.Id, &booking.ClientId, &booking.RoomId, &booking.HotelId, &booking.StartDate, &booking.EndDate); err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}
		bookings = append(bookings, booking)
	}
	json.NewEncoder(w).Encode(bookings)
}
