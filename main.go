package main

import (
	"log"
	"net/http"
	"os"

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
	if saveJson() == 1 {
		log.Fatal("could not creat json file")
	} else {
		credFile := "./credFile.json"
		scopes := []string{
			"https://www.googleapis.com/auth/userinfo.email",
		}
		secret := []byte("secret")
		sessionName := "goquestsession"
		google.Setup(redirectURL, credFile, scopes, secret)
		router.Use(google.Session(sessionName))
	}

	port := os.Getenv("PORT")
	router.Run(port)
}
