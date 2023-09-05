package service

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/storyicon/sigverify"
	"server/internal/domain"
	"time"
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
	valid, err := sigverify.VerifyEllipticCurveHexSignatureEx(
		common.HexToAddress(req.Wallet),
		[]byte(req.Message),
		req.Sign,
	)
	if err != nil {
		return errors.Wrapf(err, "%s: sign", userErrorPrefix)
	}
	if !valid {
		return errors.Wrapf(domain.ErrSignature, "%s: sign", userErrorPrefix)
	}

	_, err = s.repository.GetByWallet(ctx, req.Wallet)

	// if user not exists
	if errors.Is(err, domain.ErrNoDocuments) {

		// create new user and issue tokens
		newUser := &domain.User{
			ID:        uuid.New().String(),
			Wallet:    req.Wallet,
			Nickname:  "", // TODO: generate random nickname
			CreatedAt: time.Now(),
		}
		err = s.repository.Create(ctx, newUser)
		if err != nil {
			return errors.Wrapf(err, "%s: auth", userErrorPrefix)
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
