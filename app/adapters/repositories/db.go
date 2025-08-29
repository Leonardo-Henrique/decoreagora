package repositories

import (
	"database/sql"

	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{
		db: db,
	}
}

func (r *MySQLRepository) CreateUser(user models.User) (int64, error) {
	stmt, err := r.db.Prepare(`
        INSERT INTO users (public_id, name, email, last_login)
        VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.PublicID, user.Name, user.Email, user.LastLogin)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (r *MySQLRepository) CreateUserCredit(userID int64) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO available_credits (total, users_id)
		VALUES (?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(0, userID)
	if err != nil {
		return err
	}

	return nil
}
