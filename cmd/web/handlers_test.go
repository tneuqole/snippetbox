package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/tneuqole/snippetbox/internal/models/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, body := ts.get(t, "/ping")
	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, body, "pong")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		path     string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			path:     "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "World",
		},
		{
			name:     "Non-existent ID",
			path:     "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			path:     "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			path:     "/snippet/view/1.0",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			path:     "/snippet/view/abc",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.path)
			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	const (
		validName     = "John Smith"
		validPassword = "password"
		validEmail    = "john@example.com"
		formTag       = `<form action="/user/signup" method="POST" novalidate>`
	)

	tests := []struct {
		name        string
		userName    string
		email       string
		password    string
		csrfToken   string
		wantCode    int
		wantFormTag string
	}{
		{
			name:      "Valid submission",
			userName:  validName,
			email:     validEmail,
			password:  validPassword,
			csrfToken: csrfToken,
			wantCode:  http.StatusSeeOther,
		},
		{
			name:      "Invalid CSRF Token",
			userName:  validName,
			email:     validEmail,
			password:  validPassword,
			csrfToken: "wrongToken",
			wantCode:  http.StatusBadRequest,
		},
		{
			name:        "Empty name",
			userName:    "",
			email:       validEmail,
			password:    validPassword,
			csrfToken:   csrfToken,
			wantCode:    http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:        "Empty email",
			userName:    validName,
			email:       "",
			password:    validPassword,
			csrfToken:   csrfToken,
			wantCode:    http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:        "Empty password",
			userName:    validName,
			email:       validEmail,
			password:    "",
			csrfToken:   csrfToken,
			wantCode:    http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:        "Invalid email",
			userName:    validName,
			email:       "bob@example.",
			password:    validPassword,
			csrfToken:   csrfToken,
			wantCode:    http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:        "Short password",
			userName:    validName,
			email:       validEmail,
			password:    "pa$$",
			csrfToken:   csrfToken,
			wantCode:    http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:        "Duplicate email",
			userName:    validName,
			email:       "dupe@example.com",
			password:    validPassword,
			csrfToken:   csrfToken,
			wantCode:    http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}

func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, header, _ := ts.get(t, "/snippet/create")
		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, "/user/login", header.Get("Location"))
	})

	t.Run("Authenticated", func(t *testing.T) {
		_, _, body := ts.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body)

		form := url.Values{}
		form.Add("email", "user@example.com")
		form.Add("password", "password")
		form.Add("csrf_token", csrfToken)
		code, header, _ := ts.postForm(t, "/user/login", form)

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, "/snippet/create", header.Get("Location")) // uses redirectPath from session

		code, _, body = ts.get(t, "/snippet/create")
		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, `<form action="/snippet/create" method="POST">`)
	})
}
