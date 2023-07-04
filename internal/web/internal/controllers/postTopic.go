package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/amleonc/tabula/config"
	"github.com/amleonc/tabula/internal/dto"
	"github.com/amleonc/tabula/internal/helpers/tokens"
	"github.com/amleonc/tabula/internal/web/internal"
	"github.com/amleonc/tabula/internal/web/internal/responses"
	"github.com/labstack/echo/v4"
)

const (
	formFileName  = "file"
	formValueName = "topic"
)

var (
	ctxKey = config.UserIdKey()
)

func PostTopic(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.String(http.StatusBadRequest, "bad form")
	}
	formFile := form.File[formFileName][0]
	file, err := formFile.Open()
	if err != nil {
		return c.String(400, "bad file")
	}
	t := &dto.Topic{Media: &dto.Media{Bytes: file}}
	formValue := c.FormValue(formValueName)
	err = json.Unmarshal([]byte(formValue), t)
	if err != nil {
		return c.String(400, "bad value")
	}
	if err = c.Bind(t); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	uid, err := tokens.UserIDFromToken(c.Request().Context(), c.Get(ctxKey))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	t.User = uid
	t, err = internal.TopicService.Create(c.Request().Context(), t)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	return c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, t))
}
