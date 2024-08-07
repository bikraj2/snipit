package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"snipit.bikraj.net/ui"
)

func (app *application) routes() http.Handler {
router := httprouter.New()
router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { app.notFound(w)
})
// Leave the static files route unchanged.
fileServer := http.FileServer(http.FS(ui.Files))
router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

dynamic := alice.New(app.sessionManager.LoadAndSave,noSurf,app.authenticate)
  router.HandlerFunc(http.MethodGet,  "/ping",ping) 
  router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home)) 
  router.Handler(http.MethodGet, "/about",dynamic.ThenFunc(app.about)) 
  router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
  // Routes for Authentication

 
  router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignUp)) 
  router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignUpPost)) 
  router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin)) 
  router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost)) 
  router.Handler(http.MethodGet, "/account/password/update",dynamic.ThenFunc(app.changePassword))
  router.Handler(http.MethodPost,"/account/password/update" , dynamic.ThenFunc(app.userChangePasswordPost))
  protected :=dynamic.Append(app.requireAuthentication)
  // Routes for Snippets
  router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate)) 
  router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
  router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
  // Routes for Account Viewing
  router.Handler(http.MethodGet,  "/account/view", protected.ThenFunc(app.accountView))
standard := alice.New(app.recoverPanic, app.logRequest, secureHeader )
return standard.Then(router)
}
