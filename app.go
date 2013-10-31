package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
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
	r.HandleFunc("/{title:[a-zA-Z0-9_-]+}", UpdateHandler).Methods("POST")
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	p, err := GetPage(title)
	if p == nil {
		log.Println(err)
		http.Redirect(w, r, "/"+vars["title"]+"/edit", http.StatusFound)
		return
	}
	//TODO: handle err
	templates.ExecuteTemplate(w, "show.html", p)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	p, err := GetPage(title)
	if p == nil {
		log.Println(err)
		p = &Page{Title: title}
	}
	//TODO: handle err
	templates.ExecuteTemplate(w, "edit.html", p)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	body := r.FormValue("body")
	ioutil.WriteFile(title, []byte(body), 0600)
	http.Redirect(w, r, "/"+vars["title"], http.StatusFound)
	return
}

type Page struct {
	Title string
	Body  []byte
}

func GetPage(title string) (*Page, error) {
	bytes, err := ioutil.ReadFile(title)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: bytes}, nil
}

//type PageStore interface {
//Get(url string) (*Page, error)
//}

//type FilePageStore struct {
//Dir string
//}
