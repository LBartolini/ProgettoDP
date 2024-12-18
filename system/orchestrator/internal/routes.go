package internal

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func isLoggedIn(c *gin.Context) bool {
	return sessions.Default(c).Get("username") != nil
}

func Authorized(c *gin.Context) {
	if !isLoggedIn(c) {
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
	if isLoggedIn(c) {
		c.Redirect(http.StatusFound, "/private")
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func (r *MyRoutes) LoginRoute(c *gin.Context) {
	session := sessions.Default(c)

	username, password := c.PostForm("username"), c.PostForm("password")
	res, err := r.orchestrator.Login(username, password)

	if res && err == nil {
		session.Set("username", username)
		session.Save()
		c.Redirect(http.StatusFound, "/private")
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (r *MyRoutes) RegisterRoute(c *gin.Context) {
	session := sessions.Default(c)

	username, password, email, phone := c.PostForm("username"), c.PostForm("password"), c.PostForm("email"), c.PostForm("phone")
	res, err := r.orchestrator.Register(username, password, email, phone)

	if res && err == nil {
		session.Set("username", username)
		session.Save()
		c.Redirect(http.StatusFound, "/private")
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (r *MyRoutes) LogoutRoute(c *gin.Context) {
	session := sessions.Default(c)

	session.Delete("username")
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func (r *MyRoutes) HomeRoute(c *gin.Context) {
	// TODO: fetch position in leaderboard

	c.HTML(http.StatusOK, "home.html", gin.H{
		"username": sessions.Default(c).Get("username"),
	})
}

func (r *MyRoutes) GarageRoute(c *gin.Context) {
	// TODO: fetch motorcycles and money

	c.HTML(http.StatusOK, "garage.html", gin.H{})
}

func (r *MyRoutes) LeaderboardRoute(c *gin.Context) {
	// TODO:  fetch top 25 of leaderboard

	c.HTML(http.StatusOK, "leaderboard.html", gin.H{})
}

func (r *MyRoutes) RaceHistoryRoute(c *gin.Context) {
	// TODO:  fetch all races (all the users with motorcycle info that partecipated)

	c.HTML(http.StatusOK, "race_history.html", gin.H{})
}
