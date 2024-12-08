package hotel_svc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"project/models"
)

type HotelService struct {
	db *sql.DB
}

func CreateHotelService(dbURL string) (HotelService, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return HotelService{}, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return HotelService{}, fmt.Errorf("failed to ping database: %w", err)
	}
	fmt.Println("Connected to database successfully!")
	create_hotels_query := `CREATE TABLE IF NOT EXISTS Hotels (
		id SERIAL PRIMARY KEY,
		owner_id INT,
    	name VARCHAR(100) NOT NULL,
    	description VARCHAR(255) NOT NULL
	)`
	create_rooms_query := `CREATE TABLE IF NOT EXISTS Rooms (
		id SERIAL PRIMARY KEY,
		hotel_id INT,
		price INT,
    	available BOOL
	)`
	db.Exec(create_hotels_query)
	db.Exec(create_rooms_query)
	return HotelService{db}, nil
}

func (service *HotelService) AddHotel(w http.ResponseWriter, r *http.Request) {
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
	err = service.db.QueryRow(add_hotel_query, hotel_data.OwnerId, hotel_data.Name, hotel_data.Description).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	for _, room_price := range hotel_data.Rooms {
		service.db.Exec(add_room_query, id, room_price, true)
	}
	w.Write([]byte("Succesfully added hotel"))
}

func (service *HotelService) FindHotels(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var hotels []models.HotelInfo
	rows, err := service.db.Query("SELECT id, owner_id, name, description FROM Hotels")
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
	rows, err := service.db.Query("SELECT id, hotel_id, price, available FROM Rooms WHERE hotel_id = $1 AND available = TRUE", request_json.HotelId)
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
