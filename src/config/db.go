package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


func InitDB() *pgxpool.Pool{
	err := godotenv.Load()

	if err != nil{
		log.Fatal("Erro ao carregar variáveis de ambiente")
	}

	host:=os.Getenv("DB_HOST")
	port:=os.Getenv("DB_PORT")
	user:=os.Getenv("DB_USER")
	password:=os.Getenv("DB_PASSWORD")
	dbname:=os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`, host, port, user, password, dbname)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, connStr)
	if err != nil{
		log.Fatal("Erro ao criar conexão com banco!" + err.Error())
	}

	err = db.Ping(ctx)
	if err != nil{
		log.Fatal("Não foi possível conectar ao banco de dados: \n", err.Error())
	}

	fmt.Println("Conectado com sucesso ao PostgreSQL usando pgx!")
	return db
}