package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-kit/log"

	sharedLib "github.com/jumaroar-globant/go-bootcamp/shared"
	"github.com/jumaroar-globant/go-bootcamp/user/shared"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
)

// UserRepository defines a user repository
type UserRepository interface {
	Authenticate(ctx context.Context, username string, password string) error
	CreateUser(ctx context.Context, user *sharedLib.User) error
	UpdateUser(ctx context.Context, user *sharedLib.User) error
	GetUser(ctx context.Context, userID string) (*sharedLib.User, error)
	DeleteUser(ctx context.Context, userID string) error
}

type userRepository struct {
	db     *sql.DB
	logger log.Logger
}

//NewUserRepository is the UserRepository constructor
func NewUserRepository(db *sql.DB, logger log.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: log.With(logger, "userRepository", "sql"),
	}
}

// Authenticate is the userRepository method to authenticate a user
func (r *userRepository) Authenticate(ctx context.Context, username string, password string) error {
	userSQL := "SELECT password_hash FROM users WHERE name=?"

	user := &sharedLib.User{}

	err := r.db.QueryRowContext(ctx, userSQL, username).Scan(&user.Password)
	if err == sql.ErrNoRows {
		return ErrUserNotFound
	}

	if err != nil {
		return err
	}

	if shared.CheckPasswordHash(password, user.Password) {
		return ErrWrongPassword
	}

	return nil
}

// CreateUser is the userRepository method to create a user
func (r *userRepository) CreateUser(ctx context.Context, user *sharedLib.User) error {
	userSQL := "INSERT INTO users (id, name, password_hash, age, additional_information) VALUES(?, ?, ?, ?, ?)"

	_, err := r.db.ExecContext(ctx, userSQL, user.ID, user.Name, user.Password, user.Age, user.AdditionalInformation)
	if err != nil {
		return err
	}

	parentsSQL := "INSERT INTO user_parents (user_id, name) VALUES(?, ?)"

	for _, parent := range user.Parents {
		_, err := r.db.ExecContext(ctx, parentsSQL, user.ID, parent)
		if err != nil {
			return err
		}
	}

	return err
}

// UpdateUser is the userRepository method to update a user
func (r *userRepository) UpdateUser(ctx context.Context, user *sharedLib.User) error {
	sql := "UPDATE users SET name=?, age=?, additional_information=?  WHERE id = ?"

	_, err := r.db.ExecContext(ctx, sql, user.Name, user.Age, user.AdditionalInformation, user.ID)

	return err
}

// GetUser is the userRepository method to get a user
func (r *userRepository) GetUser(ctx context.Context, userID string) (*sharedLib.User, error) {
	sqlString := "SELECT * FROM users WHERE id=?"

	user := &sharedLib.User{}

	err := r.db.QueryRowContext(ctx, sqlString, userID).Scan(user)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}

	parentsSQLString := "SELECT name FROM user_parents WHERE user_id=?"

	rows, err := r.db.QueryContext(ctx, parentsSQLString, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var parent sharedLib.Parent
		if err := rows.Scan(&parent.UserID, &parent.Name); err != nil {
			return user, err
		}

		user.Parents = append(user.Parents, parent.Name)
	}

	return user, rows.Err()
}

// DeleteUser is the userRepository method to delete a user
func (r *userRepository) DeleteUser(ctx context.Context, userID string) error {
	parentsSQLString := "DELETE FROM user_parents WHERE user_id=?"

	_, err := r.db.ExecContext(ctx, parentsSQLString, userID)
	if err != nil {
		return err
	}

	sqlString := "DELETE FROM users WHERE id=?"

	_, err = r.db.ExecContext(ctx, sqlString, userID)
	if err != nil {
		return err
	}

	return err
}
