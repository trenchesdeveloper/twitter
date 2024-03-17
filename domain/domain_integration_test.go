package domain

import (
	"context"
	"log"
	"os"
	"testing"

	twitter "github.com/trenchesdeveloper/tweeter"
	"github.com/trenchesdeveloper/tweeter/config"
	"github.com/trenchesdeveloper/tweeter/jwt"
	"github.com/trenchesdeveloper/tweeter/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	conf             *config.Config
	db               *postgres.DB
	authService      *AuthService
	authTokenService twitter.AuthTokenService
	userRepo         twitter.UserRepo
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	passwordCost = bcrypt.MinCost

	// load config
	config.LoadEnv(".env.test")

	conf = config.New()
	log.Println(conf)
	db = postgres.New(ctx, conf)
	defer db.Close()

	err := db.Drop()

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	userRepo = postgres.NewUserRepo(db)

	authTokenService = jwt.NewTokenService(conf)

	authService = NewAuthService(userRepo, authTokenService)

	os.Exit(m.Run())
}
