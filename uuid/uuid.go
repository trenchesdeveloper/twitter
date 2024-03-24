package uuid

import "github.com/google/uuid"

func Generate() string {
	return uuid.NewString()
}

func Validate(id string) bool {
	_, err := uuid.Parse(id)

	return err == nil
}
