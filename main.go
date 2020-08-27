package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/meatballhat/negroni-logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
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
	//set var
	host := "localhost"
	port := "443"
	output := ""

	args := os.Args[1:]

	for {
		if len(args) < 2 {
			break
		} else if args[0] == "-h" || args[0] == "--host" {
			host = args[1]

			args = args[2:]
		} else if args[0] == "-p" || args[0] == "--port" {
			port = args[1]

			args = args[2:]
		} else if args[0] == "-l" || args[0] == "--log" {
			output = args[1]

			args = args[2:]
		} else {
			log.Fatalln(fmt.Sprintf("Unknown parameter: %s", args[0]))
		}

	}

	mux := httprouter.New()

	mux.GET("/", index)

	mux.NotFound = http.HandlerFunc(notFound)

	mux.PanicHandler = errorHandler

	l := log.New()

	var f *os.File
	var err error

	if output != "" {
		f, err = os.OpenFile(output, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		l.SetOutput(f)
	}

	n := negroni.New()
	n.Use(negronilogrus.NewMiddlewareFromLogger(l, "webserver"))
	n.UseHandler(mux)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: n,
	}
	l.Println(fmt.Sprintf("Run the Web Server for TLS at %s:%s", host, port))
	l.Fatal(server.ListenAndServeTLS("./ssl/localhost+1.pem", "./ssl/localhost+1-key.pem"))
	// server.ListenAndServeTLS("./ssl/localhost+1.pem", "./ssl/localhost+1-key.pem")
}
