package internal

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Authorized(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("username")
	if user == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}

type MyRoutes struct {
	orchestrator *MyOrchestrator
}

func NewMyRoutes(orchestrator *MyOrchestrator) *MyRoutes {
	return &MyRoutes{orchestrator}
}

func (r *MyRoutes) IndexRoute(c *gin.Context) {
	session := sessions.Default(c)

	count := 0
	if session.Get("count") != nil {
		count = session.Get("count").(int)
	}

	session.Set("count", count+1)
	session.Save()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"message": count,
		"title":   "Welcome",
	})
}

func (r *MyRoutes) LoginRoute(c *gin.Context) {
	session := sessions.Default(c)

	// TODO: Contact Auth Service to perform authentication
	if 1 == 1 {
		session.Set("username", "my_user")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/private")
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

func (r *MyRoutes) LogoutRoute(c *gin.Context) {
	session := sessions.Default(c)

	session.Delete("username")
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func (r *MyRoutes) HomeRoute(c *gin.Context) {
	session := sessions.Default(c)

	count := 0
	if session.Get("count") != nil {
		count = session.Get("count").(int)
	}

	session.Set("count", count+1)
	session.Save()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"message": count,
		"title":   "Home",
	})
}

func (r *MyRoutes) GarageRoute(c *gin.Context) {
	session := sessions.Default(c)

	count := 0
	if session.Get("count") != nil {
		count = session.Get("count").(int)
	}

	session.Set("count", count+1)
	session.Save()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"message": count,
		"title":   "Home",
	})
}
