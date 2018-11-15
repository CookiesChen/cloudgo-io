package controller

import (
	"fmt"
	"github.com/unrolled/render"
	"html/template"
	"log"
	"net/http"
	"time"
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
		t, _ := template.ParseFiles("./public/template/form.html")
		log.Println(t.Execute(w, struct {
			Name      string `json:"name"`
			Password string `json:"password"`
		}{Name: r.Form["username"][0], Password: r.Form["password"][0]}))
	}
}

func UnknownHandler(w http.ResponseWriter, r *http.Request)  {
	http.Error(w, "no such directory", 500)
}

func TimeHandler(w http.ResponseWriter, r *http.Request)  {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	formatter.JSON(w, http.StatusOK, struct {
		Time string `json:"time"`
	}{Time: time.Now().String()})
}