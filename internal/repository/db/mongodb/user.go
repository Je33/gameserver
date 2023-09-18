package mongodb

import (
	"context"
	"fmt"
	"server/internal/config"
	"server/internal/domain"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	userTable       = "user"
	userErrorPrefix = "[repository.db.mongodb.user]"
)

type UserMongoRepo struct {
	db *DB
}

type userDB struct {
	ID        string    `bson:"_id,omitempty"`
	Nickname  string    `bson:"nickname,omitempty"`
	Wallet    string    `bson:"wallet,omitempty"`
	CreatedAt time.Time `bson:"createdAt,omitempty"`
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
		FindOne(ctx, bson.D{{Key: "_id", Value: objId}}).Decode(&userDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.Wrapf(domain.ErrNoDocuments, "%s: get by id", userErrorPrefix)
		}
		return nil, errors.Wrapf(err, "%s: get by id", userErrorPrefix)
	}
	return &domain.User{
		ID:        userDb.ID,
		Nickname:  userDb.Nickname,
		Wallet:    userDb.Wallet,
		CreatedAt: userDb.CreatedAt,
	}, nil
}

func (repo *UserMongoRepo) GetByWallet(ctx context.Context, wallet string) (*domain.User, error) {
	userDb := &userDB{}
	cfg := config.Get()
	err := repo.db.Client.Database(cfg.MongoDB).Collection(userTable).
		FindOne(ctx, bson.D{{Key: "wallet", Value: wallet}}).Decode(&userDb)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.Wrapf(domain.ErrNoDocuments, "%s: get by wallet", userErrorPrefix)
		}
		return nil, errors.Wrapf(err, "%s: get by wallet", userErrorPrefix)
	}
	return &domain.User{
		ID:        userDb.ID,
		Nickname:  userDb.Nickname,
		Wallet:    userDb.Wallet,
		CreatedAt: userDb.CreatedAt,
	}, nil
}

func (repo *UserMongoRepo) Create(ctx context.Context, user *domain.User) error {
	userDb := &userDB{
		ID:        user.ID,
		Nickname:  user.Nickname,
		Wallet:    user.Wallet,
		CreatedAt: user.CreatedAt,
	}
	cfg := config.Get()
	result, err := repo.db.Client.Database(cfg.MongoDB).Collection(userTable).
		InsertOne(ctx, userDb)
	fmt.Println(result)
	if err != nil {
		return errors.Wrapf(err, "%s: create", userErrorPrefix)
	}
	return nil
}

func (repo *UserMongoRepo) Update(ctx context.Context, user *domain.User) error {
	userDb := &userDB{
		Nickname:  user.Nickname,
		Wallet:    user.Wallet,
		CreatedAt: user.CreatedAt,
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
		DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return errors.Wrapf(err, "%s: delete", userErrorPrefix)
	}
	if result.DeletedCount == 0 {
		return errors.Wrapf(domain.ErrNoDocuments, "%s: delete", userErrorPrefix)
	}
	return nil
}
