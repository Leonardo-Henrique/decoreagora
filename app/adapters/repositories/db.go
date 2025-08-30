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

func (r *MySQLRepository) CheckIfEmailIsRegistered(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *MySQLRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := `SELECT id, public_id, name, email, last_login FROM users WHERE email = ?`
	row := r.db.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.PublicID, &user.Name, &user.Email, &user.LastLogin)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *MySQLRepository) CreateAccessCode(accessCode models.AccessCode) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO access_codes (users_id, used, code, expire_at)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(accessCode.UserID, accessCode.IsUsed, accessCode.Code, accessCode.ExpireAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLRepository) GetAccessCodeByUserID(userID int) (models.AccessCode, error) {
	var accessCode models.AccessCode
	query := `SELECT code, expire_at, used FROM access_codes WHERE users_id = ?`
	row := r.db.QueryRow(query, userID)

	err := row.Scan(&accessCode.Code, &accessCode.ExpireAt, &accessCode.IsUsed)
	if err != nil {
		return models.AccessCode{}, err
	}
	accessCode.UserID = userID
	return accessCode, nil
}

func (r *MySQLRepository) DeleteAccessCode(userID int, code string) error {
	stmt, err := r.db.Prepare(`DELETE FROM access_codes WHERE users_id = ? AND code = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, code)
	return err
}
