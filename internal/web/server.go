package web

import (
	"fmt"

	jwx "github.com/amleonc/echo-middleware-jwx"
	"github.com/amleonc/tabula/config"
	"github.com/amleonc/tabula/internal/helpers/keys"
	"github.com/amleonc/tabula/internal/web/internal/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/lestrrat-go/jwx/v2/jwa"
)

type Router struct {
	server *echo.Echo
	port   int
}

var (
	ctxKey = config.UserIdKey()
)

func (r *Router) Start() {
	mountPublicRoutes(r)
	private := r.server.Group("/private")
	setUpPrivateRoutes(private)
	r.server.Logger.Fatal(r.server.Start(fmt.Sprintf(":%d", r.port)))
}

func New() Router {
	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 2 << 10,
		LogLevel:  log.ERROR,
	}))
	p := config.ServerPort()
	return Router{
		server: e,
		port:   p,
	}
}

func mountPublicRoutes(r *Router) {
	r.server.POST("/signup", controllers.Signup)
	r.server.POST("/login", controllers.Login)
}

func setUpPrivateRoutes(r *echo.Group) {
	cfg := jwx.DefaultConfig
	cfg.ContextKey = ctxKey
	cfg.SignatureAlgorithm = jwa.RS256
	cfg.Key = keys.PublicJWKKey()
	cfg.TokenLookup = "cookie:jwt"
	r.Use(jwx.WithConfig(cfg))
	r.POST("/topic", controllers.PostTopic)
}
