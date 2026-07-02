package main

import (
	"fmt"
	"os"
	"teste-golang-api/src/config"
	"teste-golang-api/src/handlers"
	"teste-golang-api/src/services"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	_"teste-golang-api/docs"
)


// @title           API de Reservas de Salas
// @version         1.0
// @description     Esta é uma API de reservas de salas usando Gin e Swagger.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Lucas Terêncio
// @contact.email  lucasterencio145@gmail.com

// @host      localhost:3000
func main(){
	port := os.Getenv("PORT")
	db := config.InitDB()
	defer db.Close()

	dbURL := fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	config.RunMigrations(dbURL)

	userService := services.NewUserService(db)
	userController := handlers.NewUserController(userService)

	roomService := services.NewRoomService(db)
	roomController := handlers.NewRoomController(roomService)

	bookingService := services.NewBookingService(db)
	bookingController := handlers.NewBookingController(bookingService)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userRoute := router.Group("/user")
	{
		userRoute.POST("/", userController.CreateUser)
		userRoute.DELETE("/delete/:id", userController.DeleteUser)
		userRoute.GET("/", userController.ListUsers)
	}

	roomRoute := router.Group("/room")
	{
		roomRoute.POST("/", roomController.CreateRoom)
		roomRoute.DELETE("/delete/:id", roomController.DeleteRoom)
		roomRoute.PUT("/:id", roomController.EditCapacitityRoom)
		roomRoute.GET("/", roomController.ListRooms)
	}

	bookingRoute := router.Group("/booking")
	{
		bookingRoute.POST("/", bookingController.CreateBooking)
		bookingRoute.GET("/", bookingController.ListBookings)
		bookingRoute.GET("/by-user", bookingController.ListBookingsByUser)
		bookingRoute.GET("/by-room", bookingController.ListBookingsByRoom)
		bookingRoute.PUT("/:id", bookingController.UpdateBookingStatus)
	}


	router.Run(":" + port)

}