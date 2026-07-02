package services

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"teste-golang-api/src/models"
	"teste-golang-api/src/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingService struct {
	db *pgxpool.Pool
}

func NewBookingService(db *pgxpool.Pool) *BookingService {
	return &BookingService{db: db}
}

//Função para criar uma nova reserva no banco de dados
//Recebe room_id, user_id, title, capacitity_people, startDate e endDate
//Devolve a reserva criada ou erro de validação/conflito
func (service *BookingService) CreateBookingService(ctx context.Context, room_id int, user_id int, title string, capacitity_people int8, startDate time.Time, endDate time.Time) (models.Booking, error) {

	if !endDate.After(startDate) {
		return models.Booking{}, utils.ErrValidation("end_date deve ser posterior a start_date")
	}

	queryUser := `SELECT id from users where id = $1`
	var idUser int32

	errUser := service.db.QueryRow(ctx, queryUser, user_id).Scan(&idUser)

	if errUser != nil {
		if errors.Is(errUser, sql.ErrNoRows) {
			return models.Booking{}, utils.ErrNotFound("Usuario com id %d nao existe", user_id)
		}
		return models.Booking{}, utils.WrapInternal(errUser, "Erro ao consultar usuario")
	}

	queryRoom := `SELECT id, capacitity, is_active from rooms where id = $1`
	var idRoom int32
	var roomCapacitity int32
	var isActive bool

	errRoom := service.db.QueryRow(ctx, queryRoom, room_id).Scan(&idRoom, &roomCapacitity, &isActive)

	if errRoom != nil {
		if errors.Is(errRoom, sql.ErrNoRows) {
			return models.Booking{}, utils.ErrNotFound("Quarto com id %d nao existe", room_id)
		}
		return models.Booking{}, utils.WrapInternal(errRoom, "Erro ao consultar sala")
	}

	if int8(capacitity_people) > int8(roomCapacitity) {
		return models.Booking{}, utils.ErrValidation("A sala permite no maximo %d pessoas", roomCapacitity)
	}

	queryConflict := `SELECT id FROM bookins WHERE room_id = $1 AND status = true AND start_date < $3 AND end_date > $2`
	var conflictID int32
	errConflict := service.db.QueryRow(ctx, queryConflict, room_id, startDate, endDate).Scan(&conflictID)
	if errConflict == nil {
		return models.Booking{}, utils.ErrConflict("Sala ja possui reserva neste periodo")
	}
	if !errors.Is(errConflict, sql.ErrNoRows) {
		return models.Booking{}, utils.WrapInternal(errConflict, "Erro ao verificar conflito de reserva")
	}

	queryCreate := `INSERT INTO bookins (room_id, user_id, title, capacitity_people, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id int32

	err := service.db.QueryRow(ctx, queryCreate, room_id, user_id, title, capacitity_people, startDate, endDate).Scan(&id)

	if err != nil {
		return models.Booking{}, utils.WrapInternal(err, "Erro ao criar reserva")
	}

	booking := models.Booking{
		ID:               int(id),
		RoomID:           room_id,
		UserID:           user_id,
		Title:            title,
		CapacitityPeople: capacitity_people,
		StartDate:        startDate,
		EndDate:          endDate,
	}

	return booking, err

}

//Função para listar reservas com paginação e filtro opcional por status
//Recebe status ("active" ou vazio), page e limit
//Devolve PaginatedResponse com lista de reservas
func (service *BookingService) ListBookingsService(ctx context.Context, status string, page int, limit int) (utils.PaginatedResponse, error) {
	whereClause := ""
	if status == "active" {
		whereClause = " WHERE status = true"
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM bookins` + whereClause
	err := service.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return utils.PaginatedResponse{}, utils.WrapInternal(err, "Erro ao contar reservas")
	}

	offset := (page - 1) * limit
	dataQuery := `SELECT * FROM bookins` + whereClause + ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := service.db.Query(ctx, dataQuery, limit, offset)

	if err != nil {
		return utils.PaginatedResponse{}, utils.WrapInternal(err, "Erro ao listar reservas")
	}
	defer rows.Close()

	var bookings []models.Booking

	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.RoomID,
			&booking.UserID,
			&booking.Title,
			&booking.CapacitityPeople,
			&booking.StartDate,
			&booking.EndDate,
			&booking.Status,
			&booking.CreatedAt,
		)

		if err != nil {
			return utils.PaginatedResponse{}, utils.WrapInternal(err, "Erro ao ler reserva")
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return utils.PaginatedResponse{}, utils.WrapInternal(err, "Erro ao listar reservas")
	}

	if bookings == nil {
		bookings = []models.Booking{}
	}

	totalPages := (total + limit - 1) / limit

	return utils.PaginatedResponse{
		Data:       bookings,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

//Função para listar reservas de um usuário pelo nome
//Recebe user_name
//Devolve slice de Booking ou erro
func (service *BookingService) ListBookingByUserService(ctx context.Context, user_name string) ([]models.Booking, error) {
	query := `SELECT b.id, b.room_id, b.user_id, b.title, b.capacitity_people, b.start_date, b.end_date, b.status, b.created_at
	          FROM bookins b
	          INNER JOIN users u ON b.user_id = u.id
	          WHERE u.name = $1`

	rows, err := service.db.Query(ctx, query, user_name)

	if err != nil {
		return nil, utils.WrapInternal(err, "Erro ao listar reservas por usuario")
	}
	defer rows.Close()

	var bookings []models.Booking

	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.RoomID,
			&booking.UserID,
			&booking.Title,
			&booking.CapacitityPeople,
			&booking.StartDate,
			&booking.EndDate,
			&booking.Status,
			&booking.CreatedAt,
		)

		if err != nil {
			return nil, utils.WrapInternal(err, "Erro ao ler reserva")
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.WrapInternal(err, "Erro ao listar reservas por usuario")
	}

	return bookings, nil
}

//Função para listar reservas de uma sala pelo ID
//Recebe room_id
//Devolve slice de Booking ou erro
func (service *BookingService) ListBookingsByRoomService(ctx context.Context, room_id int) ([]models.Booking, error) {
	query := `SELECT id, room_id, user_id, title, capacitity_people, start_date, end_date, status, created_at
	          FROM bookins
	          WHERE room_id = $1`

	rows, err := service.db.Query(ctx, query, room_id)

	if err != nil {
		return nil, utils.WrapInternal(err, "Erro ao listar reservas por sala")
	}
	defer rows.Close()

	var bookings []models.Booking

	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.RoomID,
			&booking.UserID,
			&booking.Title,
			&booking.CapacitityPeople,
			&booking.StartDate,
			&booking.EndDate,
			&booking.Status,
			&booking.CreatedAt,
		)

		if err != nil {
			return nil, utils.WrapInternal(err, "Erro ao ler reserva")
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.WrapInternal(err, "Erro ao listar reservas por sala")
	}

	return bookings, nil
}

//Função para desativar uma reserva (soft delete)
//Recebe booking_id
//Devolve true se desativado ou erro se não encontrada
func (service *BookingService) UpdateBookingStatusService(ctx context.Context, booking_id int) (bool, error) {
	query := `UPDATE bookins SET status = false WHERE id = $1 AND status = true`
	result, err := service.db.Exec(ctx, query, booking_id)

	if err != nil {
		return false, utils.WrapInternal(err, "Erro ao desativar reserva")
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return false, utils.ErrNotFound("Reserva nao existe ou ja esta inativa")
	}

	return true, nil
}
