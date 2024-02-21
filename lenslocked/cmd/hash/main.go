package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	secretKeyForHash := "secret-key"
	password := "this is a totally secret password nobody will guess"

	h := hmac.New(sha256.New, []byte(secretKeyForHash))

	h.Write([]byte(password))

	result := h.Sum(nil)

	fmt.Println(hex.EncodeToString(result))
}
