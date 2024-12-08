package hotel_svc

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type HotelDB struct {
	Db *sql.DB
}

func CreateHotelService(dbURL string) (HotelDB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return HotelDB{}, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return HotelDB{}, fmt.Errorf("failed to ping database: %w", err)
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
	return HotelDB{db}, nil
}
