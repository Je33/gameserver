package handler_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/internal/config"
	"server/internal/repository/db/mongodb"
	"server/internal/service"
	"server/internal/transport/rest/handler"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	userService *service.UserService
	userHandler *handler.UserHandler
	authToken   string
}

type UserAuthReq struct {
	Wallet  string `json:"wallet"`
	Message string `json:"message"`
	Sign    string `json:"sign"`
}

type UserAuthRes struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewUserTestSuite() (*UserTestSuite, error) {
	suite := &UserTestSuite{}

	ctx := context.Background()

	db, err := mongodb.Connect(ctx)
	if err != nil {
		return nil, err
	}

	userRepo := mongodb.NewUserRepo(db)

	suite.userService = service.NewUserService(userRepo)
	suite.userHandler = handler.NewUserHandler(suite.userService)

	return suite, nil
}

func TestHttpTestSuite(t *testing.T) {
	userTestSuite, err := NewUserTestSuite()
	require.NoError(t, err)
	suite.Run(t, userTestSuite)
}

func (suite *UserTestSuite) TestAuth_Success() {
	privateKey, err := crypto.GenerateKey()
	require.NoError(suite.T(), err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(suite.T(), ok)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	messageString := "test"
	messageBytes := []byte(messageString)
	messageHash := crypto.Keccak256Hash(messageBytes)

	signature, err := crypto.Sign(messageHash.Bytes(), privateKey)
	require.NoError(suite.T(), err)

	userAuthReq := UserAuthReq{
		Wallet:  address,
		Message: messageString,
		Sign:    hexutil.Encode(signature),
	}

	reqBody, err := json.Marshal(userAuthReq)
	require.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err = suite.userHandler.Auth(c)

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, rec.Code)

		authRes := UserAuthRes{}
		json.Unmarshal(rec.Body.Bytes(), &authRes)

		assert.NotEmpty(suite.T(), authRes.AuthToken)

		suite.authToken = authRes.AuthToken
	}
}

func (suite *UserTestSuite) TestAuth_Fail() {
	userAuthReq := UserAuthReq{
		Wallet:  "0x543A7060C8bB455294319b23D825478B2b798c0E",
		Message: "test",
		Sign:    "0x499cf8ce848eac151a49d23f95ce3fbfc7bf9bac709445458ca34de7afe8d98f2b9e2a34b4f7cf447145efef707355da665189cdea83a7464c804bda4cec556400",
	}

	reqBody, err := json.Marshal(userAuthReq)
	require.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err = suite.userHandler.Auth(c)

	assert.Error(suite.T(), err)
}

func (suite *UserTestSuite) TestMe_Success() {
	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set(echo.HeaderAuthorization, `Bearer `+suite.authToken)
	rec := httptest.NewRecorder()

	e := echo.New()

	c := e.NewContext(req, rec)

	cfg := config.Get()
	err := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.JWTSecret),
	})(suite.userHandler.Me)(c)

	assert.NoError(suite.T(), err)
}
