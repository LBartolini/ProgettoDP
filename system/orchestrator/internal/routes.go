package internal

import (
	"net/http"
	"orchestrator/internal/services"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func isLoggedIn(c *gin.Context) bool {
	// user is logged if username key is in session
	return sessions.Default(c).Get("username") != nil
}

func Authorized(c *gin.Context) {
	// Handle that checks if the user is authorized
	if !isLoggedIn(c) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}

type MyRoutes struct {
	orchestrator *Orchestrator
}

func NewMyRoutes(orchestrator *Orchestrator) *MyRoutes {
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
		// if login is successful, save username in session
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
		// if registration is successful, save username in session
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
	username := sessions.Default(c).Get("username").(string)

	info, err := r.orchestrator.GetLeaderboardInfo(username)

	points := 0
	position := 0

	if err == nil {
		points = int(info.Points)
		position = int(info.Position)
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"username": username,
		"points":   points,
		"position": position,
	})
}

func (r *MyRoutes) GarageRoute(c *gin.Context) {
	username := sessions.Default(c).Get("username").(string)

	money, _ := r.orchestrator.GetUserMoney(username)
	owned, err := r.orchestrator.GetUserMotorcycles(username)
	if err != nil {
		owned = make([]*services.Ownership, 0)
	}
	not_owned, err := r.orchestrator.GetRemainingMotorcycles(username)
	if err != nil {
		not_owned = make([]*services.Motorcycle, 0)
	}

	c.HTML(http.StatusOK, "garage.html", gin.H{
		"money":     money,
		"not_owned": not_owned,
		"owned":     owned,
	})
}

func (r *MyRoutes) GarageBuyRoute(c *gin.Context) {
	username := sessions.Default(c).Get("username").(string)

	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/private")
		return
	}

	_ = r.orchestrator.BuyMotorcycle(username, id)

	c.Redirect(http.StatusSeeOther, "/private/garage")
}

func (r *MyRoutes) GarageUpgradeRoute(c *gin.Context) {
	username := sessions.Default(c).Get("username").(string)

	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/private")
		return
	}

	_ = r.orchestrator.UpgradeMotorcycle(username, id)

	c.Redirect(http.StatusSeeOther, "/private/garage")
}

func (r *MyRoutes) RaceStartRoute(c *gin.Context) {
	username := sessions.Default(c).Get("username").(string)
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/private")
		return
	}

	err = r.orchestrator.StartMatchmaking(username, id)
	if err != nil {
		c.Redirect(http.StatusFound, "/private")
		return
	}

	c.Redirect(http.StatusSeeOther, "/private/garage")
}

func (r *MyRoutes) LeaderboardRoute(c *gin.Context) {
	leaderboard, err := r.orchestrator.GetFullLeaderboard()
	if err != nil {
		leaderboard = make([]*services.LeaderboardPosition, 0)
	}

	c.HTML(http.StatusOK, "leaderboard.html", gin.H{
		"leaderboard": leaderboard,
	})
}

func (r *MyRoutes) RaceHistoryRoute(c *gin.Context) {
	username := sessions.Default(c).Get("username").(string)
	results, err := r.orchestrator.GetHistory(username)
	if err != nil {
		c.Redirect(http.StatusFound, "/private")
		return
	}

	c.HTML(http.StatusOK, "race_history.html", gin.H{"results": results})
}
