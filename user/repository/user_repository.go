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
	CreateUser(ctx context.Context, user sharedLib.User) error
	UpdateUser(ctx context.Context, user sharedLib.User) (sharedLib.User, error)
	GetUser(ctx context.Context, userID string) (sharedLib.User, error)
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
	user := sharedLib.User{}

	err := r.db.QueryRowContext(ctx, PasswordHashQuery, username).Scan(&user.Password)
	if err == sql.ErrNoRows {
		return ErrUserNotFound
	}

	if err != nil {
		return err
	}

	if !shared.CheckPasswordHash(password, user.Password) {
		return ErrWrongPassword
	}

	return nil
}

// CreateUser is the userRepository method to create a user
func (r *userRepository) CreateUser(ctx context.Context, user sharedLib.User) error {
	_, err := r.db.ExecContext(ctx, InsertUserStatement, user.ID, user.Name, user.Password, user.Age, user.AdditionalInformation)
	if err != nil {
		return err
	}

	for _, parent := range user.Parents {
		_, err := r.db.ExecContext(ctx, InsertParentStatement, user.ID, parent)
		if err != nil {
			return err
		}
	}

	return err
}

// UpdateUser is the userRepository method to update a user
func (r *userRepository) UpdateUser(ctx context.Context, user sharedLib.User) (sharedLib.User, error) {
	_, err := r.db.ExecContext(ctx, UpdateUserStatement, user.Name, user.Age, user.AdditionalInformation, user.ID)
	if err != nil {
		return sharedLib.User{}, err
	}

	_, err = r.db.ExecContext(ctx, DeleteUserParentsStatement, user.ID)
	if err != nil {
		return sharedLib.User{}, err
	}

	for _, parent := range user.Parents {
		_, err := r.db.ExecContext(ctx, InsertParentStatement, user.ID, parent)
		if err != nil {
			return sharedLib.User{}, err
		}
	}

	return r.GetUser(ctx, user.ID)
}

// GetUser is the userRepository method to get a user
func (r *userRepository) GetUser(ctx context.Context, userID string) (sharedLib.User, error) {
	user := sharedLib.User{}

	err := r.db.QueryRowContext(ctx, UserDataQuery, userID).Scan(&user.ID, &user.Name, &user.Age, &user.AdditionalInformation)
	if err == sql.ErrNoRows {
		return sharedLib.User{}, ErrUserNotFound
	}

	if err != nil {
		return sharedLib.User{}, err
	}

	rows, err := r.db.QueryContext(ctx, UserParentsQuery, userID)
	if err != nil {
		return sharedLib.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var parent sharedLib.Parent
		if err := rows.Scan(&parent.Name); err != nil {
			return user, err
		}

		user.Parents = append(user.Parents, parent.Name)
	}

	return user, rows.Err()
}

// DeleteUser is the userRepository method to delete a user
func (r *userRepository) DeleteUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, DeleteUserParentsStatement, userID)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, DeleteUserStatement, userID)
	if err != nil {
		return err
	}

	return err
}
