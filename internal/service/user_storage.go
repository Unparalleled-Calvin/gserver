package service

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/Unparalleled-Calvin/gserver/internal/schema"
	"github.com/Unparalleled-Calvin/gserver/internal/storage"
)

func RegisterUser(userRegister schema.UserRegister) error {
	sum := sha256.Sum256([]byte(userRegister.Password))
	userRegister.Password = hex.EncodeToString(sum[:])

	err := storage.InsertUser(userRegister)
	return err
}
