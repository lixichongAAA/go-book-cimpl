package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getApi(c *gin.Context) {
	fmt.Println(c.Query("id"))
	c.String(http.StatusOK, "ok")
}

func postApi(c *gin.Context) {
	fmt.Println(c.PostForm("id"))
	c.String(http.StatusOK, "ok")
}

func postjson(c *gin.Context) {
	var data = &struct {
		Name string `json:"title"`
	}{}

	c.BindJSON(data)
	fmt.Println(data)
	c.String(http.StatusOK, "ok")
}

// 全局中间件,允许跨域
func GlobalMiddleWare(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Next()
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		authorized := check(token)
		if authorized {
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "UnAuthorized",
		})
		c.Abort()
		return
	}
}

func check(tk string) bool {
	return true
}

func main() {
	r := gin.Default()

	r.GET("/home", AuthMiddleWare(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "home"})
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/home/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	r.GET("/home/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	r.GET("/getApi", getApi)     //注册接口
	r.POST("/postApi", postApi)  //注册接口
	r.POST("postjson", postjson) //注册接口
	r.Use(GlobalMiddleWare)
	_ = r.Run(":8000")
}
