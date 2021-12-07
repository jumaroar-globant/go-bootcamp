package config

import (
	"context"
	"encoding/json"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// ErrForcedFailure when forced failure is required
	ErrForcedFailure = errors.New("forced failure response")
	// ErrForcedConnectionFailure when forced connection failure is required
	ErrForcedConnectionFailure = errors.New("forced connection failure response")

	// ForceMockFail is variable to force a failure for mocked functions
	ForceMockFail = false
	// ForceNotFound is variable to cause empty results
	ForceNotFound = false
	// ForceConnectionFail is a variable to force a connection failure
	ForceConnectionFail = false
	//ForceDecodeError is a variable to force a decode error
	ForceDecodeError = false
	// ForceInsertFail is a variable to force an insert error
	ForceInsertFail = false
	// ForceInsertFail is a variable to force a delete error
	ForceDeleteFail = false

	MockedItem []byte
)

// MockMongoDbHelper defines a mock struct to be used
type MockMongoDbHelper struct {
	MongoDatabaseHelper
}

// MockMongoDbClient defines a mock struct to be used
type MockMongoDbClient struct {
	MongoClientHelper
}

// MockMongoDbCollection defines a mock struct to be used
type MockMongoDbCollection struct {
	MongoCollectionHelper
}

type MockSingleResultHelper struct {
	MongoSingleResultHelper
}

//MockSample is editable sample to mock item stored in mongo
var MockSample interface {
	Decode(v interface{}) error
	Err() error
}

func (c MockMongoDbHelper) Collection(collectionName string) MongoCollectionHelper {
	return MockMongoDbCollection{}
}

func (c MockMongoDbHelper) Client() MongoClientHelper {
	return MockMongoDbClient{}
}

func (c MockMongoDbClient) Connect() error {
	if ForceConnectionFail {
		return ErrForcedConnectionFailure
	}

	return nil
}

func (c MockMongoDbClient) Database(databaseName string) MongoDatabaseHelper {
	var db MongoDatabaseHelper
	return db
}

func (c MockSingleResultHelper) Err() error {
	if ForceMockFail {
		return ErrForcedFailure
	}

	if ForceNotFound {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (c MockSingleResultHelper) Decode(v interface{}) error {
	if ForceDecodeError {
		return ErrForcedFailure
	}

	return json.Unmarshal(MockedItem, v)
}

func (c MockMongoDbCollection) FindOne(ctx context.Context, filter interface{}) MongoSingleResultHelper {
	return MockSingleResultHelper{}
}

func (c MockMongoDbCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts *options.FindOneAndUpdateOptions) MongoSingleResultHelper {
	return MockSingleResultHelper{}
}

func (c MockMongoDbCollection) InsertOne(ctx context.Context, item interface{}) (interface{}, error) {
	if ForceMockFail {
		return nil, ErrForcedFailure
	}

	if ForceInsertFail {
		return nil, ErrForcedFailure
	}

	return &mongo.InsertOneResult{InsertedID: 123}, nil
}

func (c MockMongoDbCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	if ForceMockFail {
		return 0, ErrForcedFailure
	}

	if ForceDeleteFail {
		return 0, ErrForcedFailure
	}

	return 1, nil
}
