package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romain-h/gone-fishing/internal/cache"
	"github.com/romain-h/gone-fishing/internal/config"
	"github.com/romain-h/gone-fishing/internal/expenses"
	"github.com/romain-h/gone-fishing/internal/oauth"
	"golang.org/x/oauth2"
)

type Server struct {
	cfg       config.Config
	engine    *gin.Engine
	cache     cache.CacheManager
	monzo     oauth.AuthProvider
	splitwise oauth.AuthProvider
}

func New() *Server {
	oauth2.RegisterBrokenAuthHeaderProvider("https://api.monzo.com/oauth2/token")
	oauth2.RegisterBrokenAuthHeaderProvider("https://secure.splitwise.com/oauth/token")

	cfg := config.New()
	r := gin.Default()
	cache := cache.New(*cfg)
	monzo := oauth.NewMonzo(*cfg, cache)
	splitwise := oauth.NewSplitwise(*cfg, cache)

	s := &Server{
		engine:    r,
		cache:     cache,
		monzo:     monzo,
		splitwise: splitwise,
	}
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/assets", "./public")

	// Service worker need to be served at root...
	r.GET("/sw.js", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "public/sw.js")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/monzo/link", func(c *gin.Context) {
		url := s.monzo.GetRedirectURL()
		c.Redirect(302, url)
	})

	r.GET("/monzo/callback", func(c *gin.Context) {
		code := c.Query("code")
		if err := monzo.Callback(code); err != nil {
			c.String(401, err.Error())
			return
		}
		c.String(200, "ok")
	})

	r.GET("/splitwise/link", func(c *gin.Context) {
		url := splitwise.GetRedirectURL()
		c.Redirect(302, url)
	})
	r.GET("/splitwise/callback", func(c *gin.Context) {
		code := c.Query("code")
		if err := splitwise.Callback(code); err != nil {
			c.String(401, err.Error())
			return
		}
		c.String(200, "ok")
	})

	r.GET("/", func(c *gin.Context) {
		exps := expenses.GetAllExpenses(*cfg, cache, monzo, splitwise)
		mean, median := expenses.GetStats(exps, true)
		all := expenses.GetExpensesByWeek(exps)

		c.HTML(http.StatusOK, "expenses.html", gin.H{
			"mean":     mean,
			"median":   median,
			"expenses": all,
		})
	})

	r.GET("/refresh", func(c *gin.Context) {
		expenses.FetchMonzoTransactions(*cfg, cache, monzo)
		expenses.FetchSplitwiseExpenses(cache, splitwise)
		c.Redirect(302, "/")
	})

	return s
}

func (s *Server) Run() {
	s.engine.Run()
}
