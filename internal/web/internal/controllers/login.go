package controllers

import (
	"net/http"
	"time"

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
	token, err := tokens.TokenWithClaims(map[string]any{
		"id":   u.ID,
		"role": u.Role,
		"name": u.Name,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	expires := time.Now().Add(time.Hour * 24 * 30)
	tc := http.Cookie{
		Name:    "jwt",
		Value:   token,
		Expires: expires,
		Secure:  true,
	}
	c.SetCookie(&tc)
	return c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, u))
}
