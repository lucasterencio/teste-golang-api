package models

import (
	"time"
)

type User struct{
	ID   		int 		`json:"id"`
	Name 		string      `json:"name" binding:"required,min=3,max=100"`
	Email 		string		`json:"email" binding:"required,email"`
	CreatedAt 	time.Time 	`json:"created_at"`
}


const (
	CreateTableUser = `CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY, 
	name VARCHAR(100) NOT NULL, 
	email VARCHAR(100) NOT NULL UNIQUE, 
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP)`
)