package main

import (
	"os"

	"orchestrator/internal"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	r := gin.Default()
	r.Use(sessions.Sessions("session", store))
	r.LoadHTMLGlob("./templates/*")

	routes := internal.NewMyRoutes()

	r.GET("/", routes.IndexRoute)

	r.POST("/login", routes.LoginRoute)

	private := r.Group("/private")
	private.Use(internal.Authorized)
	{
		private.GET("/", routes.HomeRoute)
		private.POST("/logout", routes.LogoutRoute)
		private.GET("/garage", routes.GarageRoute)
	}

	r.Run("0.0.0.0:8080")
}
