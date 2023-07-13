package controllers

import (
	"net/http"

	"github.com/amleonc/tabula/internal/web/internal"
	"github.com/amleonc/tabula/internal/web/internal/responses"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

func GetTopicByID(c echo.Context) error {
	id := c.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	t, err := internal.TopicService.ReadOneView(c.Request().Context(), uid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	return c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, t))
}