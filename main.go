package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var templates *template.Template

func init() {
	// call ParseFiles once at program initialization,
	// parsing all templates into a single *Template.
	// Then we can use the ExecuteTemplate method to render
	//a specific template.
	templates = template.Must(template.ParseFiles("edit.html", "view.html"))
}

type page struct {
	title string
	body  []byte
}

func (p *page) save() error {
	filename := p.title + ".txt"
	return ioutil.WriteFile(filename, p.body, 0600)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		// The http.Redirect function adds an HTTP status code
		// of http.StatusFound (302) and a Location header to
		// the HTTP response.
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	b := r.FormValue("body")
	p := &page{title: title, body: []byte(b)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &page{title: title}
	}
	renderTemplate(w, "edit", p)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	//http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderTemplate(w http.ResponseWriter, tpl string, p *page) {
	err := templates.ExecuteTemplate(w, tpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadPage(t string) (*page, error) {
	fn := t + ".txt"
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	return &page{title: t, body: b}, nil
}
