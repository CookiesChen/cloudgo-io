package router

import (
	"github.com/CookiesChen/cloudgo-io/controller"
	"github.com/gorilla/mux"
	"net/http"
)

var R *mux.Router

func init()  {
	R = mux.NewRouter()

	R.HandleFunc("/", controller.HomeHandler).Methods("GET")
	R.HandleFunc("/form", controller.FormHandler).Methods("POST")
	R.HandleFunc("/time", controller.TimeHandler).Methods("GET")
	R.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	R.PathPrefix("/").HandlerFunc(controller.UnknownHandler)
}