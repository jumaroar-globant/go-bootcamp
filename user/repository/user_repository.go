package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-kit/log"
)

var (
	ErrMissingUserName = errors.New("missing username")
	ErrMissingPassword = errors.New("missing password")
	ErrMissingUserID   = errors.New("missing user id")
	ErrUserNotFound    = errors.New("user not found")
)

type UserRepository interface {
	Authenticate(ctx context.Context, username string, password string) error
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, userID string) (*User, error)
	DeleteUser(ctx context.Context, userID string) error
}

type userRepository struct {
	db     *sql.DB
	logger log.Logger
}

func NewUserRepository(db *sql.DB, logger log.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: log.With(logger, "userRepository", "sql"),
	}
}

func (r *userRepository) Authenticate(ctx context.Context, username string, password string) error {

	return nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *User) error {
	if user.Name == "" {
		return ErrMissingUserName
	}

	if user.PasswordHash == "" {
		return ErrMissingPassword
	}

	sql := "INSERT INTO users VALUES($1, $2, $3, $4, $5)"

	_, err := r.db.ExecContext(ctx, sql, user.ID, user.Name, user.PasswordHash, user.Age, user.AdditionalInformation)
	if err != nil {
		return err
	}

	parentsSQL := "INSERT INTO user_parents VALUES($1, $2)"

	for _, parent := range user.Parents {
		_, err := r.db.ExecContext(ctx, parentsSQL, user.ID, parent)
		if err != nil {
			return err
		}
	}

	return err
}

func (r *userRepository) UpdateUser(ctx context.Context, user *User) error {
	if user.Name == "" {
		return ErrMissingUserName
	}

	sql := "UPDATE users SET name=$1, age=$2, additional_information=$3  WHERE id = $4"

	_, err := r.db.ExecContext(ctx, sql, user.Name, user.Age, user.AdditionalInformation, user.ID)

	return err
}

func (r *userRepository) GetUser(ctx context.Context, userID string) (*User, error) {
	if userID == "" {
		return nil, ErrMissingUserID
	}

	sqlString := "SELECT * FROM users WHERE id=?"

	user := &User{}

	err := r.db.QueryRowContext(ctx, sqlString, userID).Scan(user)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}

	parentsSQLString := "SELECT name FROM user_parents WHERE user_id=$1"

	rows, err := r.db.QueryContext(ctx, parentsSQLString, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var parent Parent
		if err := rows.Scan(&parent.UserID, &parent.Name); err != nil {
			return user, err
		}

		user.Parents = append(user.Parents, parent.Name)
	}

	return user, rows.Err()
}

func (r *userRepository) DeleteUser(ctx context.Context, userID string) error {
	if userID == "" {
		return ErrMissingUserID
	}

	parentsSQLString := "DELETE FROM user_parents WHERE user_id=$1"

	_, err := r.db.ExecContext(ctx, parentsSQLString, userID)
	if err != nil {
		return err
	}

	sqlString := "DELETE FROM users WHERE id=$1"

	_, err = r.db.ExecContext(ctx, sqlString, userID)
	if err != nil {
		return err
	}

	return err
}
