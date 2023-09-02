package mongodb

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"server/internal/config"
	"server/internal/domain"
	"time"
)

var (
	userTable       = "user"
	userErrorPrefix = "[repository.db.mongodb.user]"
)

type UserMongoRepo struct {
	db *DB
}

type userDB struct {
	id        string    `bson:"_id,omitempty"`
	nickname  string    `bson:"nickname,omitempty"`
	wallet    string    `bson:"wallet,omitempty"`
	createdAt time.Time `bson:"createdAt,omitempty"`
}

func NewUserRepo(db *DB) *UserMongoRepo {
	return &UserMongoRepo{db}
}

func (repo *UserMongoRepo) GetById(ctx context.Context, id string) (*domain.User, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrapf(domain.ErrConversion, "%s: user id from hex", userErrorPrefix)
	}
	userDb := &userDB{}
	cfg := config.Get()
	err = repo.db.Client.Database(cfg.MongoDB).Collection(userTable).
		FindOne(ctx, bson.D{{"_id", objId}}).Decode(&userDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.Wrapf(domain.ErrNotFound, "%s: get by id", userErrorPrefix)
		}
		return nil, errors.Wrapf(err, "%s: get by id", userErrorPrefix)
	}
	return &domain.User{
		ID:        userDb.id,
		Nickname:  userDb.nickname,
		Wallet:    userDb.wallet,
		CreatedAt: userDb.createdAt,
	}, nil
}

func (repo *UserMongoRepo) GetByWallet(ctx context.Context, wallet string) (*domain.User, error) {
	userDb := &userDB{}
	cfg := config.Get()
	err := repo.db.Client.Database(cfg.MongoDB).Collection(userTable).
		FindOne(ctx, bson.D{{"wallet", wallet}}).Decode(&userDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.Wrapf(domain.ErrNotFound, "%s: get by wallet", userErrorPrefix)
		}
		return nil, errors.Wrapf(err, "%s: get by wallet", userErrorPrefix)
	}
	return &domain.User{
		ID:        userDb.id,
		Nickname:  userDb.nickname,
		Wallet:    userDb.wallet,
		CreatedAt: userDb.createdAt,
	}, nil
}

func (repo *UserMongoRepo) Create(ctx context.Context, user *domain.User) error {
	userDb := &userDB{
		id:        user.ID,
		nickname:  user.Nickname,
		wallet:    user.Wallet,
		createdAt: user.CreatedAt,
	}
	cfg := config.Get()
	result, err := repo.db.Client.Database(cfg.MongoDB).Collection(userTable).
		InsertOne(ctx, userDb)
	if err != nil {
		return errors.Wrapf(err, "%s: create", userErrorPrefix)
	}
	_, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.Wrapf(domain.ErrConversion, "%s: user id to hex", userErrorPrefix)
	}
	return nil
}

func (repo *UserMongoRepo) Update(ctx context.Context, user *domain.User) error {
	userDb := &userDB{
		nickname:  user.Nickname,
		wallet:    user.Wallet,
		createdAt: user.CreatedAt,
	}
	cfg := config.Get()
	result, err := repo.db.Client.Database(cfg.MongoDB).Collection(userTable).
		UpdateByID(ctx, user.ID, userDb)
	if err != nil {
		return errors.Wrapf(err, "%s: update", userErrorPrefix)
	}
	if result.ModifiedCount == 0 {
		return errors.Wrapf(domain.ErrNoDocuments, "%s: update", userErrorPrefix)
	}
	return nil
}

func (repo *UserMongoRepo) Delete(ctx context.Context, id string) error {
	cfg := config.Get()
	result, err := repo.db.Client.Database(cfg.MongoDB).Collection(userTable).
		DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return errors.Wrapf(err, "%s: delete", userErrorPrefix)
	}
	if result.DeletedCount == 0 {
		return errors.Wrapf(domain.ErrNoDocuments, "%s: delete", userErrorPrefix)
	}
	return nil
}
