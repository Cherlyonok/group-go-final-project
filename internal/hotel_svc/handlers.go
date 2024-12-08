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
    	location VARCHAR(255) NOT NULL
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
	add_hotel_query := "INSERT INTO Hotels (owner_id, name, location) VALUES ($1, $2, $3) RETURNING id"
	add_room_query := "INSERT INTO Rooms (hotel_id, price, available) VALUES ($1, $2, $3)"
	var id int
	err = service.db.QueryRow(add_hotel_query, hotel_data.OwnerId, hotel_data.Name, hotel_data.Location).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
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
	w.Write([]byte("Finding hotels"))
}

func (service *HotelService) GetRoomPrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Getting room price"))
}
