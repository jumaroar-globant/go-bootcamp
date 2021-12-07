package config

import (
	"context"

	"github.com/jumaroar-globant/go-bootcamp/user/shared"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoURI          = shared.GetStringEnvVar("MONGODB_URI", "mongodb://localhost:27017")
	mongoDatabaseName = shared.GetStringEnvVar("MONGODB_DB_NAME", "bootcamp")
)

type MongoDatabaseHelper interface {
	Collection(name string) MongoCollectionHelper
	Client() MongoClientHelper
}

type MongoCollectionHelper interface {
	FindOne(context.Context, interface{}) MongoSingleResultHelper
	FindOneAndUpdate(context.Context, interface{}, interface{}, *options.FindOneAndUpdateOptions) MongoSingleResultHelper
	InsertOne(context.Context, interface{}) (interface{}, error)
	DeleteOne(ctx context.Context, filter interface{}) (int64, error)
}

type MongoSingleResultHelper interface {
	Decode(v interface{}) error
	Err() error
}

type MongoClientHelper interface {
	Database(string) MongoDatabaseHelper
	Connect() error
}

type mongoClient struct {
	cl *mongo.Client
}
type mongoDatabase struct {
	db *mongo.Database
}
type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

func NewClient() (MongoClientHelper, error) {
	c, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	return &mongoClient{cl: c}, err
}

func NewDatabase(client MongoClientHelper) MongoDatabaseHelper {
	return client.Database(mongoDatabaseName)
}

func (mc *mongoClient) Database(dbName string) MongoDatabaseHelper {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) Connect() error {
	return mc.cl.Connect(context.TODO())
}

func (md *mongoDatabase) Collection(colName string) MongoCollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() MongoClientHelper {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) MongoSingleResultHelper {
	singleResult := mc.coll.FindOne(ctx, filter)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts *options.FindOneAndUpdateOptions) MongoSingleResultHelper {
	singleResult := mc.coll.FindOneAndUpdate(ctx, filter, update, opts)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.coll.InsertOne(ctx, document)
	return id.InsertedID, err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mc.coll.DeleteOne(ctx, filter)
	return count.DeletedCount, err
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

func (sr *mongoSingleResult) Err() error {
	return sr.sr.Err()
}

func ConnectToMongoDB() (MongoDatabaseHelper, error) {
	client, err := NewClient()
	if err != nil {
		return MockMongoDbHelper{}, err
	}

	err = client.Connect()
	if err != nil {
		return MockMongoDbHelper{}, err
	}

	return client.Database(mongoDatabaseName), nil
}
