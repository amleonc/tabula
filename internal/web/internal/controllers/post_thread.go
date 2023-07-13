package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/amleonc/tabula/internal/dto"
	"github.com/amleonc/tabula/internal/web/internal"
	"github.com/amleonc/tabula/internal/web/internal/responses"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

func PostThread(c echo.Context) error {
	const (
		formFileName  = "file"
		formValueName = "thread"
	)
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	formFile := form.File[formFileName][0]
	file, err := formFile.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	t := &dto.Thread{Media: &dto.Media{Bytes: file}}
	formValue := c.FormValue(formValueName)
	err = json.Unmarshal([]byte(formValue), t)
	if err != nil {
		return c.String(400, "bad value")
	}
	if err = c.Bind(t); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	t.User = c.Request().Context().Value(ctxType).(uuid.UUID)
	t, err = internal.ThreadService.Create(c.Request().Context(), t)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	return c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, t))
}
