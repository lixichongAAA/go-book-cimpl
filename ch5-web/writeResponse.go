package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Name   string
	Habits []string
}

func write(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-custom-header", "custome")
	w.WriteHeader(201)
	user := &User{
		Name:   "Lxc",
		Habits: []string{"睡觉", "发呆", "走神", "进字节"},
	}
	json, _ := json.Marshal(&user)
	w.Write(json)
}

// func main() {
// 	http.HandleFunc("/write", write)
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServer: ", err)
// 	}
// }
