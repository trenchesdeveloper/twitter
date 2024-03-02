package domain

import (
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	passwordCost = bcrypt.MinCost
	os.Exit(m.Run())
}
