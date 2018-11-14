package main

import (
	"github.com/CookiesChen/cloudgo-io/router"
	"github.com/urfave/negroni"
)

const port = "9090"

func main() {
	r := router.R

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(": " + port)
}