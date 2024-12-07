package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Hotel struct {
	ID       int
	Name     string
	Location string
	Price    float64
}

type Booking struct {
	ID        int
	HotelID   int
	ClientID  int
	StartDate time.Time
	EndDate   time.Time
}

func connectDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	fmt.Println("Connected to database successfully!")
	return db, nil
}

func addHotel(db *sql.DB, hotel Hotel) error {
	query := `INSERT INTO Hotels (name, location, price) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, hotel.Name, hotel.Location, hotel.Price)
	if err != nil {
		return fmt.Errorf("failed to add hotel: %w", err)
	}
	fmt.Println("Hotel added successfully!")
	return nil
}

func updateHotel(db *sql.DB, hotel Hotel) error {
	query := `UPDATE Hotels SET name = $1, location = $2, price = $3 WHERE id = $4`
	_, err := db.Exec(query, hotel.Name, hotel.Location, hotel.Price, hotel.ID)
	if err != nil {
		return fmt.Errorf("failed to update hotel: %w", err)
	}
	fmt.Println("Hotel updated successfully!")
	return nil
}

func getHotels(db *sql.DB) ([]Hotel, error) {
	query := `SELECT id, name, location, price FROM Hotels`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hotels: %w", err)
	}
	defer rows.Close()

	var hotels []Hotel
	for rows.Next() {
		var hotel Hotel
		err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Location, &hotel.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan hotel: %w", err)
		}
		hotels = append(hotels, hotel)
	}
	return hotels, nil
}

func getHotelPrice(db *sql.DB, hotelID int) (float64, error) {
	query := `SELECT price FROM Hotels WHERE id = $1`
	var price float64
	err := db.QueryRow(query, hotelID).Scan(&price)
	if err != nil {
		return 0, fmt.Errorf("failed to get hotel price: %w", err)
	}
	return price, nil
}

func getBookingsByClient(db *sql.DB, clientID int) ([]Booking, error) {
	query := `SELECT id, hotel_id, client_id, start_date, end_date FROM Bookings WHERE client_id = $1`
	rows, err := db.Query(query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bookings: %w", err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var booking Booking
		err := rows.Scan(&booking.ID, &booking.HotelID, &booking.ClientID, &booking.StartDate, &booking.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func getBookingsByHotel(db *sql.DB, hotelID int) ([]Booking, error) {
	query := `SELECT id, hotel_id, client_id, start_date, end_date FROM Bookings WHERE  hotel_id = $1`
	rows, err := db.Query(query, hotelID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bookings: %w", err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var booking Booking
		err := rows.Scan(&booking.ID, &booking.HotelID, &booking.ClientID, &booking.StartDate, &booking.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func addBooking(db *sql.DB, booking Booking) error {
	query := `INSERT INTO Bookings (hotel_id, client_id, start_date, end_date) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, booking.HotelID, booking.ClientID, booking.StartDate, booking.EndDate)
	if err != nil {
		return fmt.Errorf("failed to add booking: %w", err)
	}
	fmt.Println("Booking added successfully!")
	return nil
}

func main() {
	dbURL := "postgres://postgres:12345@localhost:13500?sslmode=disable"
	db, err := connectDB(dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	defer db.Close()

	// Пример работы с отелями
	newHotel := Hotel{Name: "Radisson", Location: "Zavidovo", Price: 20000.0}
	err = addHotel(db, newHotel)
	if err != nil {
		log.Fatalf("Error adding hotel: %v\n", err)
	}

	hotels, err := getHotels(db)
	if err != nil {
		log.Fatalf("Error retrieving hotels: %v\n", err)
	}
	fmt.Println("Hotels:", hotels)

	updatedHotel := Hotel{ID: 3, Name: "Radisson2", Location: "Zavidovo", Price: 20000.0}
	err = updateHotel(db, updatedHotel)
	if err != nil {
		log.Fatalf("Error updating hotel: %v\n", err)
	}

	hotels, err = getHotels(db)
	if err != nil {
		log.Fatalf("Error retrieving hotels: %v\n", err)
	}
	fmt.Println("Hotels:", hotels)

	price, err := getHotelPrice(db, 3)
	if err != nil {
		log.Fatalf("Error retrieving hotels: %v\n", err)
	}
	fmt.Println("Price:", price)

	newBooking := Booking{HotelID: 3, ClientID: 123, StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 2)}
	err = addBooking(db, newBooking)
	if err != nil {
		log.Fatalf("Error adding booking: %v\n", err)
	}

	bookings, err := getBookingsByClient(db, 123)
	if err != nil {
		log.Fatalf("Error retrieving bookings: %v\n", err)
	}
	fmt.Println("Bookings:", bookings)

	bookings, err = getBookingsByHotel(db, 3)
	if err != nil {
		log.Fatalf("Error retrieving bookings: %v\n", err)
	}
	fmt.Println("Bookings:", bookings)
}
