package router

import (
	"github.com/CookiesChen/cloudgo-io/controller"
	"github.com/gorilla/mux"
)

var R *mux.Router

func init()  {
	R = mux.NewRouter()

	R.HandleFunc("/", controller.HomeHandler)
}