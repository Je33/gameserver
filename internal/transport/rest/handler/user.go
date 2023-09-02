package handler

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"server/internal/config"
	"server/internal/domain"
	"server/internal/transport/rest/model"
	"time"
)

var (
	userErrorPrefix = "[transport.rest.handler.user]"
)

type jwtCustomClaims struct {
	Wallet string `json:"wallet"`
	jwt.RegisteredClaims
}

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	Auth(context.Context, *domain.UserAuthReq) error
	GetByWallet(ctx context.Context, wallet string) (*domain.User, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) Auth(ctx echo.Context) error {
	cfg := config.Get()

	restUserAuthReq := new(model.UserAuthReq)
	err := ctx.Bind(restUserAuthReq)
	if err != nil {
		return err
	}

	domainUserAuthReq := &domain.UserAuthReq{
		Wallet:  restUserAuthReq.Wallet,
		Sign:    restUserAuthReq.Sign,
		Message: restUserAuthReq.Message,
	}

	err = h.service.Auth(ctx.Request().Context(), domainUserAuthReq)
	if err != nil {
		return err
	}

	claims := &jwtCustomClaims{
		restUserAuthReq.Wallet,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenSign, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return err
	}

	// TODO: Implement refresh strategy

	restUserAuthRes := &model.UserAuthRes{
		AuthToken:    tokenSign,
		RefreshToken: "",
	}

	return ctx.JSON(http.StatusCreated, restUserAuthRes)
}

func (h *UserHandler) Me(ctx echo.Context) error {
	auth, ok := ctx.Get("user").(*jwt.Token)
	if !ok {
		return fmt.Errorf("%s: auth user error", userErrorPrefix)
	}
	claims, ok := auth.Claims.(*jwtCustomClaims)
	if !ok {
		return fmt.Errorf("%s: auth claims error", userErrorPrefix)
	}
	wallet := claims.Wallet
	user, err := h.service.GetByWallet(ctx.Request().Context(), wallet)
	if err != nil {
		return errors.Wrapf(err, "%s: user not found", userErrorPrefix)
	}

	return ctx.JSON(http.StatusCreated, user)
}
