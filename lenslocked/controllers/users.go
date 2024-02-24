package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/keisn1/lenslocked/models"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService *models.UserService
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("email")
	if err != nil {
		fmt.Fprint(w, "The email cookie could not be read")
		return
	}
	fmt.Fprintf(w, "Email cookie: %s\n", email.Value)
	fmt.Fprintf(w, "Headers: %+v\n", r.Header)
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

	cookie := http.Cookie{
		Name:     "email",
		Value:    user.Email,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	fmt.Fprintf(w, "User account authenticated: %v", user)
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
	var nu models.NewUser
	nu.Email = r.FormValue("email")
	nu.Password = r.FormValue("password")
	user, err := u.UserService.Create(nu)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "User created: %+v", user)
}
