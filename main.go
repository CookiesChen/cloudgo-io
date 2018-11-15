package main

import (
	"fmt"
	"github.com/CookiesChen/cloudgo-io/router"
	"github.com/urfave/negroni"
	"net/http"
)

const port = "9090"

func main() {
	r := router.R

	n := negroni.Classic()

	n.UseFunc(sayhi)
	n.UseFunc(sayhello)

	n.UseHandler(r)
	n.Run(": " + port)
}

func sayhi(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc){
	fmt.Println("Hi")
	next(rw,r)
}

func sayhello(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc){
	fmt.Println("Hello")
	next(rw,r)
}