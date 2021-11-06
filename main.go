package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zalando/gin-oauth2/google"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	//URL コントローラー的な
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/login", google.LoginHandler)

	private := router.Group("/auth")
	private.Use(google.Auth())
	private.GET("/MyPage", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	//GoogleLogin
	redirectURL := "https://inimaching.herokuapp.com/auth"
	credFile := os.Getenv("credFile") //これだとバグるからファイルじゃなくてオブジェクトから読み込む方法模索する
	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
	}
	secret := []byte("secret")
	sessionName := "goquestsession"
	google.Setup(redirectURL, credFile, scopes, secret)
	router.Use(google.Session(sessionName))

	port := os.Getenv("PORT")
	router.Run(port)
}
