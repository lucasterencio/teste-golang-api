package services

import (
	"context"
	"teste-golang-api/src/models"
	"teste-golang-api/src/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)


type RoomService struct {
	db *pgxpool.Pool
}

func NewRoomService(db *pgxpool.Pool) *RoomService {
	return &RoomService{db: db}
}

//Função para criar uma nova sala no banco de dados
//Recebe name e capacitity
//Devolve ID, name, capacitity e is_active
func (service *RoomService) CreateRoomService(ctx context.Context, name string, capacitity int8) (models.Room, error) {
	query := `INSERT INTO rooms (name, capacitity) VALUES ($1, $2) RETURNING id`

	var id int

	err := service.db.QueryRow(ctx, query, name, capacitity).Scan(&id)

	if err != nil {
		return models.Room{}, utils.WrapInternal(err, "Erro ao criar sala")
	}

	return models.Room{ID: id, Name: name, Capacitity: capacitity}, nil
}

//Função para deletar uma sala do banco de dados
//Recebe room_id
//Devolve true se deletado ou erro se não encontrada
func (service *RoomService) DeleteRoomService(ctx context.Context, room_id int) (bool, error) {
	query := `DELETE FROM rooms WHERE id = $1`
	result, err := service.db.Exec(ctx, query, room_id)

	if err != nil {
		return false, utils.WrapInternal(err, "Erro ao apagar sala")
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return false, utils.ErrNotFound("Sala nao existe!")
	}

	return true, nil
}


//Função para editar a capacidade de uma sala no banco de dados
//Recebe room_id e nova capacitity
//Devolve a sala atualizada com todos os campos
func (service *RoomService) EditCapacitiyRoomService(ctx context.Context, room_id int, capacitity int8) (models.Room, error) {
	var room models.Room
	query := `UPDATE rooms SET capacitity = $1 WHERE id = $2 RETURNING id, name, capacitity, is_active, created_at`
	err := service.db.QueryRow(ctx, query, capacitity, room_id).Scan(&room.ID, &room.Name, &room.Capacitity, &room.IsActive, &room.CreatedAt)

	if err != nil {
		return models.Room{}, utils.WrapInternal(err, "Erro ao atualizar sala")
	}

	return room, nil
}

//Função para listar todas as salas do banco de dados
//Devolve slice de Room ou erro
func (service *RoomService) ListRoomsService(ctx context.Context) ([]models.Room, error) {
	query := `SELECT id, name, is_active, capacitity, created_at FROM rooms ORDER BY created_at DESC`
	rows, err := service.db.Query(ctx, query)

	if err != nil {
		return nil, utils.WrapInternal(err, "Erro ao listar salas")
	}
	defer rows.Close()

	var rooms []models.Room

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.Name, &room.IsActive, &room.Capacitity, &room.CreatedAt)

		if err != nil {
			return nil, utils.WrapInternal(err, "Erro ao ler sala")
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.WrapInternal(err, "Erro ao listar salas")
	}

	if rooms == nil {
		rooms = []models.Room{}
	}

	return rooms, nil
}