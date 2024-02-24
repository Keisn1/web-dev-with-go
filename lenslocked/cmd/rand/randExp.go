package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	n := 32
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	if nRead < n {
		panic("didn't read enough random bytes")
	}

	fmt.Println(base64.URLEncoding.EncodeToString(b))
}
