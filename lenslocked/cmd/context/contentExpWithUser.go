package main

import (
	stdctx "context"
	"fmt"

	"github.com/keisn1/lenslocked/context"
	"github.com/keisn1/lenslocked/models"
)

func main() {
	ctx := stdctx.Background()

	user := models.User{
		Email: "kay@email.context",
	}

	ctx = context.WithUser(ctx, &user)
	retreivedUser := context.User(ctx)
	fmt.Println(retreivedUser)
}
