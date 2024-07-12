package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snipit.bikraj.net/internal/models"
	Validator "snipit.bikraj.net/internal/validator"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	Validator.Validator `form:"-"`
}
type userSignForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	Validator.Validator `form:"-"`
}
type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	Validator.Validator `form:"-"`
}
type userPasswordChangeForm struct {
	OldPassword         string `form:"OldPassword"`
	NewPassword         string `form:"NewPassword"`
	ConfirmPassword     string `form:"ConfirmPassword"`
	Validator.Validator `form:"-"`
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	snippets, err := app.snippets.Latest()
	// fmt.Println("here")
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(Validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(
		Validator.MaxChars(form.Title, 100),
		"title",
		"This field cannot be more than 100 characters long",
	)
	form.CheckField(Validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(
		Validator.PermittedValue(form.Expires, 1, 7, 365),
		"expires",
		"This field must equal 1, 7 or 365",
	)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Snippet  Successfully created.!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignUpPost(w http.ResponseWriter, r *http.Request) {
	var form userSignForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, err)
		return
	}
	form.CheckField(Validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(Validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(Validator.Matches(form.Email, Validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(Validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(Validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}
	err = app.users.Insert(form.Name, form.Email, form.Password)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.CheckField(true, "email", "Email address already in use")
			data := app.newTemplateData(r)
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(Validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(Validator.Matches(form.Email, Validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(Validator.NotBlank(form.Password), "password", "This field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	path := app.sessionManager.PopString(r.Context(), "redirectedPathAfterLogin")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You have been logged out Successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "about.tmpl.html", data)

}

func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	id := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	user, err := app.users.Get(id)
	if errors.Is(err, models.ErrNoRecord) {
		http.Redirect(w, r, "user/login", http.StatusSeeOther)
	} else if err != nil {
		app.serverError(w, err)
	}
	data.User = user
	app.render(w, http.StatusOK, "account.tmpl.html", data)

}
func (app *application) changePassword(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userPasswordChangeForm{}
	app.render(w, http.StatusOK, "changePassword.tmpl.html", data)
}

func (app *application) userChangePasswordPost(w http.ResponseWriter, r *http.Request) {
	var form userPasswordChangeForm
	// data :=app.newTemplateData(r)
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, err)
	}
	// Blank Field
	form.CheckField(Validator.NotBlank(form.OldPassword), "CurrentPassword", "Current Password cannot be blank")
	form.CheckField(Validator.NotBlank(form.NewPassword), "NewPassword", "New Password cannot be blank")
	form.CheckField(Validator.NotBlank(form.OldPassword), "OldPassword", "Old Password cannot be blank")
	// Validations
	form.CheckField(Validator.MinChars(form.NewPassword, 8), "NewPassword", "Password Length must be at least 8 characters")
	fmt.Println(form.NewPassword, form.ConfirmPassword)
	form.CheckField(Validator.Same(form.NewPassword, form.ConfirmPassword), "NewPassword", "Password Must be Same")
	form.CheckField(!Validator.Same(form.OldPassword, form.NewPassword), "NewPassword", "New Password Cannot be same as the new Password")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "changePassword.tmpl.html", data)
		return
	}
	id := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	_, err = app.users.ChangePassword(id, form.OldPassword, form.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			data := app.newTemplateData(r)
			data.Form = form

			form.AddNonFieldError("The password you have entered is incorrect")
			fmt.Println(form)
			app.render(w, http.StatusUnprocessableEntity, "changePassword.tmpl.html", data)
		} else if errors.Is(err, models.ErrNoRecord) {
			form.AddNonFieldError("Please Login Again")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
