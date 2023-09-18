package rest

import (
	"context"
	"net/http"
	"server/internal/config"
	"server/internal/repository/db/mongodb"
	"server/internal/service"
	"server/internal/transport/rest/handler"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
)

func RunRest() error {

	ctx := context.Background()

	cfg := config.Get()

	// connect db
	db, err := mongodb.Connect(ctx)
	if err != nil {
		return err
	}

	// init repositories
	userRepo := mongodb.NewUserRepo(db)

	// init services
	userService := service.NewUserService(userRepo)

	// init handlers
	userHandler := handler.NewUserHandler(userService)

	// init echo
	e := echo.New()
	// Disable Echo JSON logger in debug mode
	if cfg.LogLevel == "debug" {
		if l, ok := e.Logger.(*echoLog.Logger); ok {
			l.SetHeader("${time_rfc3339} | ${level} | ${short_file}:${line}")
		}
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// API V1
	v1 := e.Group("/v1")

	// Auth jwt request
	v1.POST("/auth", userHandler.Auth)

	// User
	r := v1.Group("/user")
	r.Use(echojwt.JWT([]byte(cfg.JWTSecret)))
	r.GET("", userHandler.Me)

	// Start server
	s := &http.Server{
		Addr:         cfg.HTTPAddr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	e.Logger.Fatal(e.StartServer(s))

	return nil
}
