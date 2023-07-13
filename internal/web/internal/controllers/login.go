package controllers

import (
	"net/http"
	"time"

	"github.com/amleonc/tabula/config"
	"github.com/amleonc/tabula/internal/dto"
	"github.com/amleonc/tabula/internal/helpers/tokens"
	"github.com/amleonc/tabula/internal/web/internal"
	"github.com/amleonc/tabula/internal/web/internal/responses"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	u := new(dto.User)
	var err error
	if err = c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	u, err = internal.UserService.Login(c.Request().Context(), u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	token, err := buildJWTToken(u.ID, u.Role, u.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	tc := buildJWTCookie(token)
	c.SetCookie(tc)
	return c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, u))
}

func buildJWTToken(id, role, name any) (string, error) {
	t, err := tokens.TokenWithClaims(map[string]any{
		"uid":  id,
		"role": role,
		"name": name,
	})
	return t, err
}

func buildJWTCookie(t string) *http.Cookie {
	exp := time.Now().Add(time.Hour * 24 * 30)
	c := http.Cookie{
		Name:    "jwt",
		Value:   t,
		Expires: exp,
	}
	if config.AppEnv() != "dev" {
		c.Secure = true
	}
	return &c
}
