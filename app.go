package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var (
	templates = template.Must(template.ParseGlob("./tmpl/*html"))
)

func main() {
	router := mux.NewRouter()
	configureRoutes(router)
	http.Handle("/", router)
	fmt.Println("Started on http://localhost:3000/")
	http.ListenAndServe(":3000", nil)
}

func configureRoutes(r *mux.Router) {
	r.HandleFunc("/{title:[a-zA-Z0-9_-]+}", ViewHandler).Methods("GET")
	r.HandleFunc("/{title:[a-zA-Z0-9_-]+}/edit", EditHandler).Methods("GET")
	r.HandleFunc("/{title:[a-zA-Z0-9_-]+}", ViewHandler).Methods("POST")
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	p := GetPage(title)
	if p == nil {
		http.Redirect(w, r, "/"+vars["title"]+"edit/", http.StatusFound)
		return
	}
	//TODO: handle err
	templates.ExecuteTemplate(w, "show.html", p)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	p := GetPage(title)
	if p == nil {
		p = &Page{Title: title}
	}
	//TODO: handle err
	templates.ExecuteTemplate(w, "edit.html", p)
}

type Page struct {
	Title string
	Body  []byte
}

func GetPage(title string) *Page {
	return &Page{Title: title}
}

//type PageStore interface {
//Get(url string) (*Page, error)
//}

//type FilePageStore struct {
//Dir string
//}
