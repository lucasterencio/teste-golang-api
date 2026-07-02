package models

import (
	"time"

)

type Booking struct{
	ID            		int			`json:"id"`
	RoomID            	int 		`json:"room_id"`
	UserID            	int 		`json:"user_id"`
	Title         		string     	`json:"title"`
	CapacitityPeople    int8    	`json:"capacitity_people"`
	StartDate         	time.Time 	`json:"start_date"`
	EndDate           	time.Time 	`json:"end_date"`
	Status      		bool    	`json:"status"`
	CreatedAt 	  		time.Time 	`json:"created_at"`
}

const (
	CreateTableBooking = `
	CREATE TABLE IF NOT EXISTS bookins (
	id SERIAL PRIMARY KEY,
	room_id INT,
	user_id INT,
	title VARCHAR(100) NOT NULL,
	capacitity_people INT NOT NULL,
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
	status BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

	CONSTRAINT fk_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
	CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`
)