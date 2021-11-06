package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zalando/gin-oauth2/google"
)

func saveJson() int {
	file, err := os.Create("credFile.json")
	if err != nil {
		log.Fatal(err)
		return 1
	}
	defer file.Close()
	file.WriteString(os.Getenv("credFile"))
	return 0
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	/*
		store := sessions.NewCookieStore([]byte("secret"))
		router.Use(sessions.Sessions("session", store))
	*/

	//URL コントローラー的な
	router.GET("/index", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/beforeLogin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "beforeLogin.tmpl", gin.H{
			"title": "beforeLogin",
		})
	})

	router.GET("/auth", func(c *gin.Context) {
		c.HTML(http.StatusOK, "oauth.tmpl", gin.H{
			"title": "logined",
		})
	})

	router.GET("/login", google.LoginHandler)

	private := router.Group("/auth")
	private.Use(google.Auth())
	private.GET("/MyPage", func(c *gin.Context) {
		c.HTML(http.StatusOK, "myPage.tmpl", gin.H{
			"title": "myPage",
		})
	})

	//GoogleLogin
	redirectURL := "https://inimaching.herokuapp.com/auth"

	var credFile string
	if os.Getenv("USER") == "Knight-of-Skyrim" {
		credFile = "./json/credFileLocal.json"
	} else {
		if saveJson() == 1 {
			log.Fatal("could not creat json file")
		} else {
			credFile = "./credFile.json"
		}
	}
	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
	}
	secret := []byte("secret")
	sessionName := "goquestsession"
	google.Setup(redirectURL, credFile, scopes, secret)
	router.Use(google.Session(sessionName))

	var port string
	if os.Getenv("USER") == "Knight-of-Skyrim" {
		port = "5000"
	} else {
		port = os.Getenv("PORT")
	}
	router.Run(":" + port)
}
