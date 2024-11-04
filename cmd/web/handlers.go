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
	Title   string               `form:"title"`
	Content string               `form:"content"`
	Expires int                  `form:"expires"`
	V       *validator.Validator `form:"-"`
}

func NewSnippetCreateForm() *snippetCreateForm {
	return &snippetCreateForm{
		Expires: 365,
		V:       validator.New(),
	}
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, r, http.StatusOK, "home.tmpl", data)
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
	data.Form = NewSnippetCreateForm()
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) postSnippet(w http.ResponseWriter, r *http.Request) {
	form := NewSnippetCreateForm()
	err := app.decodePostForm(r, form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.V.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.V.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.V.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.V.CheckField(validator.AllowedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")

	if !form.V.Valid() {
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

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
