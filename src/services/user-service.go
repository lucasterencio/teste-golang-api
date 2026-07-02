package services

import (
	"context"
	"teste-golang-api/src/models"
	"teste-golang-api/src/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

//Função para criar um novo usuário no banco de dados
//Recebe name e email
//Devolve ID e name
func (service *UserService) CreateUserService(ctx context.Context, name string, email string) (models.User, error) {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`

	var id int

	err := service.db.QueryRow(ctx, query, name, email).Scan(&id)

	if err != nil {
		return models.User{}, utils.WrapInternal(err, "Erro ao criar usuario")
	}

	return models.User{ID: id, Name: name}, nil
}

//Função para deletar um usuário do banco de dados
//Recebe user_id
//Devolve true se deletado ou erro se não encontrado
func (service *UserService) DeleteUserService(ctx context.Context, user_id int) (bool, error){

	query := `DELETE FROM users WHERE id = $1`
	result, err := service.db.Exec(ctx, query, user_id)

	if err != nil{
		return false, utils.WrapInternal(err, "Erro ao apagar usuário")
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0{
		return false, utils.ErrNotFound("Usuário não existe!")
	}

	return true, nil
}

//Função para listar todos os usuários do banco de dados
//Devolve slice de User ou erro
func (service *UserService) ListUsersService(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, name, email, created_at FROM users ORDER BY created_at DESC`
	rows, err := service.db.Query(ctx, query)

	if err != nil {
		return nil, utils.WrapInternal(err, "Erro ao listar usuarios")
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

		if err != nil {
			return nil, utils.WrapInternal(err, "Erro ao ler usuario")
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.WrapInternal(err, "Erro ao listar usuarios")
	}

	if users == nil {
		users = []models.User{}
	}

	return users, nil
}