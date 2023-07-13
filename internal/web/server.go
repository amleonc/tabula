package web

import (
	"context"
	"fmt"
	"net/http"

	jwx "github.com/amleonc/echo-middleware-jwx"
	"github.com/amleonc/tabula/config"
	"github.com/amleonc/tabula/internal/helpers/keys"
	"github.com/amleonc/tabula/internal/helpers/tokens"
	"github.com/amleonc/tabula/internal/web/internal/controllers"
	"github.com/amleonc/tabula/internal/web/internal/responses"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
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

	// just for fun
	r.server.Static("/uploads", config.UploadsDir())
}

func setUpPrivateRoutes(r *echo.Group) {
	cfg := jwx.DefaultConfig
	cfg.ContextKey = ctxKey
	cfg.SignatureAlgorithm = jwa.RS256
	cfg.Key = keys.PublicJWKKey()
	cfg.TokenLookup = "cookie:jwt"
	r.Use(jwx.WithConfig(cfg))
	r.Use(setUserID)

	r.POST("/comment", controllers.PostComment)

	r.POST("/thread", controllers.PostThread)
	r.GET("/thread/:id", controllers.GetThreadByID)

	r.POST("/topic", controllers.PostTopic)
	r.GET("/topic/:id", controllers.GetTopicByID)
}

const (
	errJWTAssertion = "cannot extract jwt token"
)

func setUserID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		i := c.Get(ctxKey)
		var t jwt.Token
		var ok bool
		if t, ok = i.(jwt.Token); !ok {
			return c.JSON(
				http.StatusInternalServerError,
				responses.NewErrorResponse(http.StatusInternalServerError, newMiddlewareError(errJWTAssertion)),
			)
		}
		uid, err := tokens.UserIDFromToken(c.Request().Context(), t)
		if err != nil {
			// TODO
			panic("implement me")
		}
		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, config.ContextKey(), uid)
		req := c.Request().WithContext(ctx)
		c.SetRequest(req)
		return next(c)
	}
}

type middlewareError struct {
	msg string
}

func (e middlewareError) Error() string {
	return e.msg
}

func newMiddlewareError(msg string) middlewareError {
	return middlewareError{msg}
}
