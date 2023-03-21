package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	ID       int
	Name     string
	Password string
}

var UserById = make(map[int]*User)
var UserByName = make(map[string][]*User)

var tpl = `<html>
<head>
<title></title>
</head>
<body>
<form action="/login" method="post">
	用户名:<input type="text" name="username">
	密码:<input type="password" name="password">
	<input type="submit" value="登录">
</form>
</body>
</html>`

func loginMemory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {
		var t *template.Template
		t = template.New("Products")
		t, _ = t.Parse(tpl)
		log.Println(t.Execute(w, nil))
	} else {
		r.ParseForm()
		fmt.Println("username: ", r.Form["username"])
		fmt.Println("password: ", r.Form["password"])
		user1 := User{ID: 1, Name: r.Form.Get("username"), Password: r.Form.Get("password")}
		store(user1)

		if pwd := r.Form.Get("password"); pwd != "123456" {
			fmt.Fprintf(w, "密码错误，请重新输入!")
		} else {
			fmt.Fprintf(w, "欢迎登录， %s", r.Form.Get("username"))
		}
	}
}

func store(user User) {
	UserById[user.ID] = &user
	UserByName[user.Name] = append(UserByName[user.Name], &user)
}

func userInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println(UserById[1])
	r.ParseForm()

	for _, user := range UserByName[r.Form.Get("username")] {
		fmt.Println(user)
		fmt.Fprintf(w, " %v", user) //这个写入到w的是输出到客户端的
	}
}

// func main() {
// 	http.HandleFunc("/login", loginMemory)   //设置访问的路由
// 	http.HandleFunc("/info", userInfo)       //设置访问的路由
// 	err := http.ListenAndServe(":8080", nil) //设置监听的端口
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }
