package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tneuqole/snippetbox/internal/models"
	"github.com/tneuqole/snippetbox/internal/validator"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type changePasswordForm struct {
	CurrentPassword     string `form:"currentPassword"`
	NewPassword         string `form:"newPassword"`
	ConfirmNewPassword  string `form:"confirmNewPassword"`
	validator.Validator `form:"-"`
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (app *application) getUserSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *application) postUser(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "Cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "Cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Invalid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "Cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "You have signed up, please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) getUserLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "Cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Invalid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "Cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddMessage("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	redirectURL := app.sessionManager.PopString(r.Context(), "redirectURL")
	if redirectURL == "" {
		redirectURL = "/"
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully.")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) getAccount(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	user, err := app.users.Get(id)
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.User = user
	app.render(w, r, http.StatusOK, "account.tmpl", data)
}

func (app *application) getChangePasswordForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = changePasswordForm{}
	app.render(w, r, http.StatusOK, "password.tmpl", data)
}

func (app *application) postChangePasswordForm(w http.ResponseWriter, r *http.Request) {
	var form changePasswordForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "Cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "Cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "Must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.ConfirmNewPassword), "confirmNewPassword", "Cannot be blank")
	form.CheckField(validator.Equals(form.NewPassword, form.ConfirmNewPassword), "confirmNewPassword", "Passwords do not match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "password.tmpl", data)
		return
	}

	id := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	user, err := app.users.Get(id)
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	_, err = app.users.Authenticate(user.Email, form.CurrentPassword)
	if err != nil {
		form.AddError("currentPassword", "Current password is incorrect")
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "password.tmpl", data)
		return
	}

	err = app.users.UpdatePassword(user.ID, form.NewPassword)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your password has been updated.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "about.tmpl", data)
}

func (app *application) getSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *application) getSnippetForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365, // probably should be a constant
	}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) postSnippet(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.AllowedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
