package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/keisn1/lenslocked/context"
	"github.com/keisn1/lenslocked/models"
)

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, CookieSession)
		if err != nil || token == "" {
			if token == "" {
				fmt.Println("token is empty string")
			} else {
				err = fmt.Errorf("readCookie in SetUser: %w", err)
				fmt.Println(err)
			}
			next.ServeHTTP(w, r)
			return
		}

		user, err := umw.SessionService.User(token)
		if err != nil {
			err = fmt.Errorf("readCookie in SetUser: %w", err)
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

// SetUser and RequireUser middleware are required
func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	nu := models.NewUser{Email: r.FormValue("email"), Password: r.FormValue("password")}
	user, err := u.UserService.Authenticate(nu)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Fprintf(w, "Email address not known")
		return
	}
	if err != nil {
		fmt.Fprintf(w, "Wrong password")
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	// render the signup page
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form submission.", http.StatusBadRequest)
	}

	// User creation logic
	var nu models.NewUser
	nu.Email = r.FormValue("email")
	nu.Password = r.FormValue("password")
	user, err := u.UserService.Create(nu)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}

	// session creation logic
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "users/me", http.StatusFound)
	fmt.Fprintf(w, "User created: %+v", user)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
	// TODO: Delete the user's cookie
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)

}
