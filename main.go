package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, "Hello World")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "Page Not Found")
}

func errorHandler(w http.ResponseWriter, r *http.Request, p interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Internal Server Error!!")
}

func main() {
	//
	mux := httprouter.New()

	mux.GET("/", index)

	mux.NotFound = http.HandlerFunc(notFound)

	mux.PanicHandler = errorHandler

	server := http.Server{
		Addr:    "0.0.0.0:443",
		Handler: mux,
	}

	server.ListenAndServeTLS("./ssl/localhost+1.pem", "./ssl/localhost+1-key.pem")
}
