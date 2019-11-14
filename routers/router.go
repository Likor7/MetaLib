package routers

import (
	"MetaLib/templmanager"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/about", aboutHandler).Methods("GET")

	r.HandleFunc("/profile", profileHandler).Methods("GET")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	r.HandleFunc("/func/auth", authHandler).Methods("POST")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")

	r.HandleFunc("/books", booksHandler).Methods("GET")
	r.HandleFunc("/authors", authorsHandler).Methods("GET")
	r.HandleFunc("/libraries", librariesHandler).Methods("GET")

	r.HandleFunc("/book/{id:[0-9]+}", bookHandler).Methods("GET")
	r.HandleFunc("/author/{id:[0-9]+}", authorHandler).Methods("GET")
	r.HandleFunc("/library/{id:[0-9]+}", LibraryHandler).Methods("GET")

	r.HandleFunc("/search", searchHandler).Methods("GET")

	r.HandleFunc("/func/star", starHandler).Methods("POST")

	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	//models.GetBookById(1) // REMOVE

	return r
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	err := templmanager.RenderTemplate(w, r, "index.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	err := templmanager.RenderTemplate(w, r, "about.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := templmanager.RenderTemplate(w, r, "404.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Error(err)
		http.Redirect(w, r, "/", 302)
		return
	}

	query := r.FormValue("oq")
	fmt.Println(query)

	err := templmanager.RenderTemplate(w, r, "search.html", struct {
		SearchFor string
	}{SearchFor: query})
	if err != nil {
		log.Fatal(err)
	}
}
