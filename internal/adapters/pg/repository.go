package postgres

import (
	"database/sql"
	"log"

	"github.com/francisco-alonso/go-template/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(user domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, email) VALUES ($1, $2)", user.Username, user.Email)
	log.Printf("Err: %v", err)
	return err
}
