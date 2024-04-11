package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"

	"github.com/EnsurityTechnologies/enscrypt"
)

func EncryptData(ss string, data string) (string, error) {
	eb, err := enscrypt.Seal(ss, []byte(data))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(eb), nil
}

func DecryptData(ss string, data string) (string, error) {
	ds, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	db, err := enscrypt.UnSeal(ss, ds)
	if err != nil {
		return "", err
	}
	return string(db), nil
}

func RandString(l int) string {
	b := make([]byte, l/2)
	rand.Read(b)
	return hex.EncodeToString(b)
}
