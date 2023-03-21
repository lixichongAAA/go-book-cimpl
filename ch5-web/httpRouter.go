package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello, Welcome back home!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello, %s \n", ps.ByName("name"))
}

// func main() {
// 	router := httprouter.New()
// 	router.GET("/", Home)
// 	router.GET("/hello/:name", Hello)

// 	log.Fatal(http.ListenAndServe(":8080", router))
// }
