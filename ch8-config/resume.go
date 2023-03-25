package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/longjoy/micro-go-book/ch8-config/conf"
	"github.com/spf13/viper"
)

func main() {
	http.HandleFunc("/resumes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "个人信息: \n")
		fmt.Fprintf(w, "姓名: %s\n性别: %s\n年龄: %d\n", viper.GetString("resume.name"), conf.Resume.Sex, conf.Resume.Age)
	})
	log.Fatal(http.ListenAndServe("8081", nil))
}
