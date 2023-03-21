package main

import (
	"fmt"
	"net/http"
	"strings"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("Path: ", r.URL.Path)
	fmt.Println("Host: ", r.Host)
	for k, v := range r.Form {
		fmt.Println("Key: ", k)
		fmt.Println("Val: ", strings.Join(v, ""))
	}
	_, _ = fmt.Fprintf(w, "Hello web, %s", r.Form.Get("name"))
}

// func main() {
// 	http.HandleFunc("/", sayHello)
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServer: ", err)
// 	}
// }
