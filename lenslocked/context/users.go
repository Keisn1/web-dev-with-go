package context

import (
	"context"

	"github.com/keisn1/lenslocked/models"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, value *models.User) context.Context {
	return context.WithValue(ctx, userKey, value)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)

	if user, ok := val.(*models.User); ok {
		// The most likely case is that nothing was ever stored in the context,
		// so it doesn't have a type of *models.User. It is also possible that
		// other code in this package wrote an invalid value using the user key,
		// so it is important to review code changes in this package.
		return user
	}

	return nil
}
