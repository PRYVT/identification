package helper

import (
	"github.com/google/uuid"
)

func GenerateGUID(input string) uuid.UUID {
	namespace := uuid.MustParse("6b37b120-9dad-1ad1-8fb4-00c04fd430c8")
	guid := uuid.NewSHA1(namespace, []byte(input))
	return guid
}
