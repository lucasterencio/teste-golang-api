package models

import (
	"time"

)

type Room struct{
	ID            int 			`json:"id"`
	Name          string     	`json:"name"`
	IsActive      bool    		`json:"is_active"`
	Capacitity    int8    		`json:"capacitity"`
	CreatedAt 	  time.Time 	`json:"created_at"`
}

const (
	CreateTableRoom = `
	CREATE TABLE IF NOT EXISTS rooms (
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	is_active BOOLEAN NOT NULL DEFAULT FALSE,
	capacitity INT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`
)