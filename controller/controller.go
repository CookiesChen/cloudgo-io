package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./public/html/index.html")
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