package mongodb

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/internal/config"
	"server/internal/domain"
)

var (
	mongodbErrorPrefix = "[repository.db.mongodb]"
)

type DB struct {
	*mongo.Client
}

func Connect(ctx context.Context) (*DB, error) {
	cfg := config.Get()
	if cfg.MongoURL == "" {
		return nil, errors.Wrapf(domain.ErrConfig, "%s: connection url not provided", mongodbErrorPrefix)
	}

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.MongoURL).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: disconnected", mongodbErrorPrefix)
	}

	// Send a ping to confirm a successful connection
	res := &DB{client}
	if err := res.Ping(ctx); err != nil {
		return nil, errors.Wrapf(err, "%s: disconnected", mongodbErrorPrefix)
	}

	return res, nil
}

func (db *DB) Ping(ctx context.Context) error {
	// Send a ping to check connection
	var result bson.M
	if err := db.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return errors.Wrapf(err, "%s: disconnected", mongodbErrorPrefix)
	}
	return nil
}
