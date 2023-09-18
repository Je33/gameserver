package service

import (
	"context"
	"server/internal/domain"
	"server/pkg/sign"
	"time"

	"github.com/lithammer/shortuuid/v3"
	"github.com/pkg/errors"
)

var (
	userErrorPrefix = "[service.user]"
)

//go:generate mockery --dir . --name UserRepository --output ./mocks
type UserRepository interface {
	GetByWallet(context.Context, string) (*domain.User, error)
	Create(context.Context, *domain.User) error
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{repository}
}

func (s *UserService) Auth(ctx context.Context, req *domain.UserAuthReq) error {

	// check wallet signature
	err := sign.VerifySignature(
		req.Wallet,
		req.Message,
		req.Sign,
	)
	if err != nil {
		return errors.Wrapf(err, "%s: signature check fail", userErrorPrefix)
	}

	_, err = s.repository.GetByWallet(ctx, req.Wallet)

	// if user not exists
	if errors.Is(err, domain.ErrNoDocuments) {

		// create new user and issue tokens
		newUser := &domain.User{
			ID:        shortuuid.New(),
			Wallet:    req.Wallet,
			Nickname:  "", // TODO: generate random nickname
			CreatedAt: time.Now(),
		}
		err = s.repository.Create(ctx, newUser)
		if err != nil {
			return errors.Wrapf(err, "%s: repo save error", userErrorPrefix)
		}
	}

	return nil
}

func (s *UserService) GetByWallet(ctx context.Context, wallet string) (*domain.User, error) {
	user, err := s.repository.GetByWallet(ctx, wallet)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: get by wallet", userErrorPrefix)
	}
	return user, nil
}
