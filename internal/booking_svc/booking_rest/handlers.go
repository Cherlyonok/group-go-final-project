package booking_rest

import (
	"encoding/json"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"project/internal/booking_svc"
	"project/internal/hotel_svc/hotel_grpc/proto/hotel_grpc"
	"project/models"
	"strconv"
)

type BookingService struct {
	BookingDB   *booking_svc.BookingDB
	HotelClient hotel_grpc.HotelServiceClient
}

func (service *BookingService) NewBookingService(bookingDB *booking_svc.BookingDB, hotelServiceAddr string) *BookingService {

	conn, err := grpc.Dial(hotelServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can't start grpc connection between booking svc and hotel svc: %v", err)
	}

	return &BookingService{
		BookingDB:   bookingDB,
		HotelClient: hotel_grpc.NewHotelServiceClient(conn),
	}
}

func (service *BookingService) AddBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var booking models.BookingInfo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&booking); err != nil {
		http.Error(w, "invalid booking data - unable to decode", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO Bookings (client_id, room_id, hotel_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	id := 0
	err := service.BookingDB.Db.QueryRow(query, booking.ClientId, booking.RoomId, booking.HotelId, booking.StartDate, booking.EndDate).Scan(&id)
	if err != nil {
		http.Error(w, "can't add this booking to database", http.StatusNotAcceptable)
		return
	}

	booking.Id = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
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
