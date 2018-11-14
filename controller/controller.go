package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./public/template/index.html")
		log.Println(t.Execute(w, nil))
	}
}

func FormHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "POST" {
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

func UnknownHandler(w http.ResponseWriter, r *http.Request)  {
	w.WriteHeader(500)
}

func ApiTestHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./public/template/jsTest.html")
	log.Println(t.Execute(w, struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	}{ID: "8675309", Content: "Hello from Go!"}))
}