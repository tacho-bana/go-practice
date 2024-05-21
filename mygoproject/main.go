package main

import (
	"html/template"
	"net/http"
)

type PageData struct {
	Title   string
	Content string
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:   "Hello, Go!",
		Content: "Welcome to the Go web server.",
	}
	tmpl, _ := template.ParseFiles("template.html")
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}
