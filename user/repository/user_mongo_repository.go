package repository

import (
	"context"

	"github.com/go-kit/log"
	"github.com/jumaroar-globant/go-bootcamp/shared"
	"github.com/jumaroar-globant/go-bootcamp/user/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	userShared "github.com/jumaroar-globant/go-bootcamp/user/shared"
)

const (
	userCollectionName = "users"
)

type userMongoRepository struct {
	db     config.MongoDatabaseHelper
	logger log.Logger
}

func NewMongoUserRepository(db config.MongoDatabaseHelper, logger log.Logger) UserRepository {
	return &userMongoRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userMongoRepository) Authenticate(ctx context.Context, username string, password string) error {
	user := shared.User{}

	filter := bson.M{"username": username}

	result := r.db.Collection(userCollectionName).FindOne(ctx, filter)
	err := result.Err()
	if err == mongo.ErrNoDocuments {
		return ErrUserNotFound
	}

	if err != nil {
		return err
	}

	err = result.Decode(&user)
	if err != nil {
		return err
	}

	if !userShared.CheckPasswordHash(password, user.Password) {
		return ErrWrongPassword
	}

	return nil
}

func (r *userMongoRepository) CreateUser(ctx context.Context, user shared.User) error {
	_, err := r.db.Collection(userCollectionName).InsertOne(ctx, user)

	return err
}

func (r *userMongoRepository) GetUser(ctx context.Context, userID string) (shared.User, error) {
	user := shared.User{}

	filter := bson.D{primitive.E{Key: "user_id", Value: userID}}

	result := r.db.Collection(userCollectionName).FindOne(ctx, filter)
	err := result.Err()
	if err == mongo.ErrNoDocuments {
		return user, ErrUserNotFound
	}

	if err != nil {
		return user, err
	}

	err = result.Decode(&user)
	if err != nil {
		return shared.User{}, err
	}

	return user, nil
}

func (r *userMongoRepository) UpdateUser(ctx context.Context, user shared.User) (shared.User, error) {
	filter := bson.M{"user_id": user.ID}

	update := bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "name", Value: user.Name},
				primitive.E{Key: "age", Value: user.Age},
				primitive.E{Key: "additional_information", Value: user.AdditionalInformation},
				primitive.E{Key: "parents", Value: user.Parents},
			},
		},
	}

	after := options.After
	upsert := false
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := r.db.Collection(userCollectionName).FindOneAndUpdate(ctx, filter, update, opts)
	err := result.Err()
	if err == mongo.ErrNoDocuments {
		return shared.User{}, ErrUserNotFound
	}

	if err != nil {
		return shared.User{}, err
	}

	err = result.Decode(&user)
	if err != nil {
		return shared.User{}, err
	}

	return user, nil
}

func (r *userMongoRepository) DeleteUser(ctx context.Context, userID string) error {
	_, err := r.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	filter := bson.M{"user_id": userID}

	_, err = r.db.Collection(userCollectionName).DeleteOne(ctx, filter)

	return err
}
