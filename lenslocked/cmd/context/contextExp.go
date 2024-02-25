package main

import (
	"context"
	"fmt"
	"strings"
)

type ctxKey string

const (
	favoriteColorKey ctxKey = "favorite-color"
)

func contextExp() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, favoriteColorKey, "blue")
	anyValue := ctx.Value(favoriteColorKey)
	// This .(string) format attempts to assert that anyValue has a type of string
	// If it succeeds, ok will be true. Otherwise ok will be false.
	stringValue, ok := anyValue.(string)
	if !ok {
		// anyValue is not a string!
		fmt.Println(anyValue, "is not a string")
		return
	}
	// if we are here, stringValue has a type of string!
	fmt.Println(stringValue, "is a string!")

	value := ctx.Value(favoriteColorKey)
	strValue, ok := value.(string)
	if !ok {
		fmt.Println("not a string")
	} else {
		fmt.Println(strings.HasPrefix(strValue, "b"))
	}
}
