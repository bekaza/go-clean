package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	_userDelivery "github.com/bekaza/go-clean/user/delivery"
	_userRepo "github.com/bekaza/go-clean/user/repository/postgres"
	_userUsecase "github.com/bekaza/go-clean/user/usecase"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var runEnv string

func init() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("configs")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("cannot read in viper config:%s", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	runEnv = viper.GetString("run.env")
	if runEnv == "" {
		runEnv = "development"
	}
}

func main() {
	e := echo.New()
	e.HideBanner = false

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	dbConn, err := newPotgresDB()
	if err != nil {
		log.Fatal(err)
	}

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	tokenTTLExpired := viper.GetDuration("auth.token.ttl")

	// init Repository
	userRepo := _userRepo.NewPostgrelUserRepository(dbConn)

	// init Service
	userService := _userUsecase.NewUserService(userRepo, []byte(viper.GetString("auth.signing.key")), tokenTTLExpired, timeoutContext)

	userMiddleware := _userDelivery.NewAuthMiddleware(userService)
	e.Use(userMiddleware.CORS)

	v1 := e.Group("/api/v1")
	// v1UserAuth := e.Group("/api/v1", userMiddleware.AuthRequire)

	// init route
	_userDelivery.NewHandler(v1, userService)

	// Start server
	go func() {
		if err := e.Start(viper.GetString("port")); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func newPotgresDB() (*gorm.DB, error) {
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok",
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.name"),
	)
	return gorm.Open(postgres.Open(connString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

}
