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

func PostComment(c echo.Context) error {
	var (
		formFileName  = "file"
		formValueName = "comment"
	)
	formFile, err := c.FormFile(formFileName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	file, err := formFile.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	formValue := c.FormValue(formValueName)
	cm := &dto.Comment{Media: &dto.Media{Bytes: file}}
	err = json.Unmarshal([]byte(formValue), cm)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	if err = c.Bind(cm); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	cm.User = c.Request().Context().Value(ctxType).(uuid.UUID)
	cm, err = internal.CommentService.Create(c.Request().Context(), cm)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, err))
	}
	return c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, cm))
}
