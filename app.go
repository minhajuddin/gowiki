package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/russross/blackfriday"
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
	r.HandleFunc("/", RootHandler).Methods("GET")
	r.HandleFunc("/{title:[a-zA-Z0-9_-]+}", ViewHandler).Methods("GET")
	r.HandleFunc("/{title:[a-zA-Z0-9_-]+}/edit", EditHandler).Methods("GET")
	r.HandleFunc("/{title:[a-zA-Z0-9_-]+}", UpdateHandler).Methods("POST")
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusFound)
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
	page := &Page{Title: vars["title"], Body: []byte(r.FormValue("body"))}
	SavePage(page)
	http.Redirect(w, r, "/"+vars["title"], http.StatusFound)
	return
}

type Page struct {
	Title string
	Body  []byte
}

func (self *Page) RenderedBody() template.HTML {
	return template.HTML(self.Body)
}

// file storage
func GetPage(title string) (*Page, error) {
	bytes, err := ioutil.ReadFile(PathForPage(&Page{Title: title}))
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: bytes}, nil
}

func SavePage(p *Page) error {
	return ioutil.WriteFile(PathForPage(p), p.Body, 0600)
}

func PathForPage(p *Page) string {
	return "./data/" + p.Title
}

//type PageStore interface {
//Get(url string) (*Page, error)
//}

//type FilePageStore struct {
//Dir string
//}
