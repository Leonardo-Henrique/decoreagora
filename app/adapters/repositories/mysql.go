package repositories

import (
	"database/sql"
	"time"

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

func (r *MySQLRepository) GetAccessCodeByUserID(userID int) ([]models.AccessCode, error) {
	var accessCodes []models.AccessCode
	query := `SELECT code, expire_at, used FROM access_codes WHERE users_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var acessCode models.AccessCode
		if rows.Scan(
			&acessCode.Code,
			&acessCode.ExpireAt,
			&acessCode.IsUsed,
		); err != nil {
			return nil, err
		}
		accessCodes = append(accessCodes, acessCode)
	}

	//accessCode.UserID = userID
	return accessCodes, nil
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

func (r *MySQLRepository) GetUserCredits(userID int) (int, error) {
	var qtdCredits int
	query := `SELECT total FROM available_credits WHERE users_id = ?`
	if err := r.db.QueryRow(query, userID).Scan(&qtdCredits); err != nil {
		return 0, err
	}
	return qtdCredits, nil
}

func (r *MySQLRepository) CreateNewImageEntry(public_id, imageKey, prompt_descr string, userID int, date time.Time) (int, error) {
	stmt, err := r.db.Prepare(`
		INSERT INTO generated_images 
		(public_id, original_image_key, prompt_description, users_id, created_at)
		VALUES (?,?,?,?,?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(public_id, imageKey, prompt_descr, userID, date)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *MySQLRepository) FinishImageEdition(generatedFileBucketKey, originalFilePublicKey string) error {
	stmt, err := r.db.Prepare(`
		UPDATE generated_images
		SET generated_image_key = ?
		WHERE public_id = ? 
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(generatedFileBucketKey, originalFilePublicKey); err != nil {
		return err
	}
	return nil
}

func (r *MySQLRepository) GetCurrentCredits(userID int) (int, error) {
	var qtd int
	query := `SELECT total FROM available_credits WHERE users_id = ?`
	if err := r.db.QueryRow(query, userID).Scan(&qtd); err != nil {
		return 0, err
	}
	return qtd, nil
}

func (r *MySQLRepository) AtomicDecrementCredit(userID int) (bool, error) {
	stmt, err := r.db.Prepare(`
        UPDATE available_credits 
        SET total = total - 1 
        WHERE users_id = ? AND total > 0
    `)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r *MySQLRepository) GetUserImages(userID int) ([]models.EditedImageResponse, error) {
	var images []models.EditedImageResponse
	query := `
		SELECT public_id, original_image_key, generated_image_key, prompt_description, created_at
		FROM generated_images
		WHERE users_id = ? ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var image models.EditedImageResponse
		var generatedImageNull sql.NullString
		var promptNull sql.NullString

		if err := rows.Scan(
			&image.PublicKey,
			&image.OriginalImageBucketKey,
			&generatedImageNull,
			&promptNull,
			&image.CreatedAt,
		); err != nil {
			return nil, err
		}

		if generatedImageNull.Valid {
			image.EditedImageBucketKey = generatedImageNull.String
		}

		if promptNull.Valid {
			image.Prompt = promptNull.String
		}

		images = append(images, image)
	}

	return images, nil
}

func (r *MySQLRepository) CreateNewSubscription(userID int, tier string, isActive bool, email string) error {
	stmt, err := r.db.Prepare(`
        INSERT INTO subscriptions (users_id, tier, is_active, user_email)
        VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(userID, tier, isActive, email); err != nil {
		return err
	}

	return nil
}

func (r *MySQLRepository) GetUserResume(userID int) (models.UserInfoMe, error) {
	var useInfo models.UserInfoMe
	query := `
		SELECT u.name, u.email, a.total, s.tier, s.is_active
		FROM users u
		INNER JOIN available_credits a
		ON u.id = a.users_id 
		INNER JOIN subscriptions s
		ON u.id = s.users_id
		WHERE u.id = ? 
	`
	if err := r.db.QueryRow(query, userID).Scan(
		&useInfo.Name,
		&useInfo.Email,
		&useInfo.AvailableCredits,
		&useInfo.Tier,
		&useInfo.IsPlanActive,
	); err != nil {
		return models.UserInfoMe{}, err
	}

	return useInfo, nil
}

func (r *MySQLRepository) UpdateUserCustomerID(userID int, customerID string) error {
	stmt, err := r.db.Prepare(`
		UPDATE subscriptions 
		SET stripe_costumer_id = ?
		WHERE users_id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(customerID, userID); err != nil {
		return err
	}
	return nil
}

func (r *MySQLRepository) GetSubscriptionByEmail(email string) (models.Subscription, error) {
	query := `
		SELECT id, stripe_costumer_id, stripe_subscription_id, 
		stripe_price_id, is_active, tier, user_email, users_id
		FROM subscriptions
		WHERE user_email = ?
	`

	var sub models.Subscription
	var customerIdNull sql.NullString
	var subIdNull sql.NullString
	var priceIdNull sql.NullString
	var userEmailNull sql.NullString

	if err := r.db.QueryRow(query, email).Scan(
		&sub.ID,
		&customerIdNull,
		&subIdNull,
		&priceIdNull,
		&sub.IsActive,
		&sub.Tier,
		&userEmailNull,
		&sub.UserID,
	); err != nil {
		if err != sql.ErrNoRows {
			return models.Subscription{}, err
		}
	}

	if customerIdNull.Valid {
		sub.StripeCostumerID = customerIdNull.String
	}

	if subIdNull.Valid {
		sub.StripeSubscriptionID = subIdNull.String
	}

	if priceIdNull.Valid {
		sub.StripePriceID = priceIdNull.String
	}

	if userEmailNull.Valid {
		sub.Email = userEmailNull.String
	}

	return sub, nil
}

func (r *MySQLRepository) IncrementUserCreditsByCustomerID(userID, qtd int) error {
	stmt, err := r.db.Prepare(`
		UPDATE available_credits
		SET total = total + ?
		WHERE users_id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(qtd, userID); err != nil {
		return err
	}
	return nil
}

func (r *MySQLRepository) GetSubscriptionByCustomerID(customerID string) (models.Subscription, error) {
	query := `
		SELECT id, stripe_costumer_id, stripe_subscription_id, 
		stripe_price_id, is_active, tier, user_email, users_id
		FROM subscriptions
		WHERE stripe_costumer_id = ?
	`

	var sub models.Subscription
	var customerIdNull sql.NullString
	var subIdNull sql.NullString
	var priceIdNull sql.NullString
	var userEmailNull sql.NullString

	if err := r.db.QueryRow(query, customerID).Scan(
		&sub.ID,
		&customerIdNull,
		&subIdNull,
		&priceIdNull,
		&sub.IsActive,
		&sub.Tier,
		&userEmailNull,
		&sub.UserID,
	); err != nil {
		if err != sql.ErrNoRows {
			return models.Subscription{}, err
		}
	}

	if customerIdNull.Valid {
		sub.StripeCostumerID = customerIdNull.String
	}

	if subIdNull.Valid {
		sub.StripeSubscriptionID = subIdNull.String
	}

	if priceIdNull.Valid {
		sub.StripePriceID = priceIdNull.String
	}

	if userEmailNull.Valid {
		sub.Email = userEmailNull.String
	}

	return sub, nil
}

func (r *MySQLRepository) UpdateUserTier(customerID, tier string) error {
	stmt, err := r.db.Prepare(`
		UPDATE subscriptions 
		SET tier = ?
		WHERE stripe_costumer_id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(tier, customerID); err != nil {
		return err
	}
	return nil
}

func (r *MySQLRepository) CreatePaymentHistoryEntry(paymentEntry models.PaymentHistory) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO payment_history
		(stripe_customer_id, processed_at, stripe_price_id, amount_paid, credits_received, public_id)
		VALUES (?,?,?,?,?,?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		paymentEntry.CustomerID,
		paymentEntry.ProcessedAt,
		paymentEntry.StripePriceID,
		paymentEntry.AmountPaid,
		paymentEntry.CreditsReceived,
		paymentEntry.PublicID,
	); err != nil {
		return err
	}

	return nil
}
