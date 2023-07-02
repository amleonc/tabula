package controllers

import (
	"net/http"

	"github.com/amleonc/tabula/internal/dto"
	"github.com/amleonc/tabula/internal/web/internal"
	"github.com/amleonc/tabula/internal/web/internal/responses"
	"github.com/labstack/echo/v4"
)

func Signup(c echo.Context) error {
	u := new(dto.User)
	var err error
	if err = c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	u, err = internal.UserService.Signup(c.Request().Context(), u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	return c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, u))
}

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
	return c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, u))
}
