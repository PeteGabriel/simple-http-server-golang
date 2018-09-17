package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var templates *template.Template
var validPath *regexp.Regexp

func init() {
	// call ParseFiles once at program initialization,
	// parsing all templates into a single *Template.
	// Then we can use the ExecuteTemplate method to render
	//a specific template.
	templates = template.Must(template.ParseFiles("edit.html", "view.html"))

	// small protection against user input
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request, ttl string) {
	p, err := loadPage(ttl)
	if err != nil {
		// The http.Redirect function adds an HTTP status code
		// of http.StatusFound (302) and a Location header to
		// the HTTP response.
		http.Redirect(w, r, "/edit/"+ttl, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, ttl string) {
	b := r.FormValue("body")
	p := &Page{Title: ttl, Body: []byte(b)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+ttl, http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request, ttl string) {
	p, err := loadPage(ttl)
	if err != nil {
		p = &Page{Title: ttl}
	}
	renderTemplate(w, "edit", p)
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request, ttl string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func renderTemplate(w http.ResponseWriter, tpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadPage(t string) (*Page, error) {
	fn := t + ".txt"
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	return &Page{Title: t, Body: b}, nil
}
