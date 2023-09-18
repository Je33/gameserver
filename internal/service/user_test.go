package service

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"server/internal/domain"
	"server/internal/service/mocks"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lithammer/shortuuid/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_Auth(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	messageString := "test"
	messageBytes := []byte(messageString)
	messageHash := crypto.Keccak256Hash(messageBytes)

	signature, err := crypto.Sign(messageHash.Bytes(), privateKey)
	require.NoError(t, err)

	userAuthReq := &domain.UserAuthReq{
		Wallet:  address,
		Message: messageString,
		Sign:    hexutil.Encode(signature),
	}

	userAuth := &domain.User{
		ID:        shortuuid.New(),
		Wallet:    address,
		Nickname:  "",
		CreatedAt: time.Now(),
	}

	testCases := []struct {
		name         string
		expectations func(context.Context, *mocks.UserRepository)
		input        *domain.UserAuthReq
		err          error
	}{
		{
			name:  "success auth",
			input: userAuthReq,
			expectations: func(ctx context.Context, userRepo *mocks.UserRepository) {
				userRepo.On("GetByWallet", ctx, userAuthReq.Wallet).Return(nil, nil)
				userRepo.On("Create", ctx, userAuth).Return(nil)
			},
		},
		{
			name:  "failed auth",
			input: userAuthReq,
			expectations: func(ctx context.Context, userRepo *mocks.UserRepository) {
				userRepo.On("GetByWallet", ctx, userAuthReq.Wallet).Return(nil, nil)
				userRepo.On("Create", ctx, userAuth).Return(errors.New("error"))
			},
		},
	}

	for _, test := range testCases {
		t.Logf("testing %s", test.name)

		ctx := context.Background()

		userRepo := mocks.NewUserRepository(t)
		userService := NewUserService(userRepo)

		test.expectations(ctx, userRepo)

		err := userService.Auth(ctx, test.input)

		if err != nil {
			if test.err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.NoError(t, err)
			}
		}

		userRepo.AssertExpectations(t)

	}
}

func TestUserService_GetByWallet(t *testing.T) {

}
