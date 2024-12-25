package booking_svc

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type BookingDB struct {
	Db *sql.DB
}

func CreateBookingService(dbURL string) (BookingDB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return BookingDB{}, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return BookingDB{}, fmt.Errorf("failed to ping database: %w", err)
	}
	fmt.Println("Connected to database successfully!")

	create_bookings_query := `CREATE TABLE IF NOT EXISTS Bookings (
		id SERIAL PRIMARY KEY,
		client_id INT NOT NULL,
    	room_id INT NOT NULL,
    	hotel_id INT NOT NULL,
		start_date DATE NOT NULL,
		end_date DATE NOT NULL
	)`

	db.Exec(create_bookings_query)
	return BookingDB{db}, nil
}
