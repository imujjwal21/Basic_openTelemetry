package httptransport

import (
	"html/template"
	"log"
	"login/users"
	"net/http"

	"github.com/gorilla/mux"
)

var templates *template.Template

func NewHandler(storage users.Storage) http.Handler {

	router := mux.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", "text/html")
			next.ServeHTTP(w, r)
		})
	})

	templates = template.Must(template.ParseGlob("templates/*.html"))

	router.HandleFunc("/", indexGetHandler(storage)).Methods("GET")
	//router.HandleFunc("/", indexPostHandler(storage)).Methods("POST")

	router.HandleFunc("/login", loginGetHandler(storage)).Methods("GET")
	router.Path("/login").Methods(http.MethodPost).HandlerFunc(loginPostHandler(storage))

	router.HandleFunc("/register", registerGetHandler(storage)).Methods("GET")
	router.HandleFunc("/register", registerPostHandler(storage)).Methods("POST")

	fs := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	return router
}

func indexGetHandler(storage users.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		templates.ExecuteTemplate(w, "index.html", nil)
	}
}

func loginGetHandler(storage users.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		templates.ExecuteTemplate(w, "login.html", nil)
	}
}

func loginPostHandler(storage users.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		name := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		err := storage.Check(r.Context(), name, password)

		if err != nil {
			log.Printf("id or pass may be incorrect : %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", 302)

	}
}

func registerGetHandler(storage users.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		templates.ExecuteTemplate(w, "register.html", nil)

	}
}

func registerPostHandler(storage users.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		name := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		err := storage.Create(r.Context(), name, password)

		if err != nil {
			log.Printf("cannot able Created Account : %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", 302)

	}
}
