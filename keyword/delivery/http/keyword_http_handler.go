package http

import (
	"context"
	"net/http"
	"olx-crawler/errors"
	_i18n "olx-crawler/i18n"
	"olx-crawler/keyword"
	"olx-crawler/models"
	"olx-crawler/utils"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/labstack/echo/v4"
)

type handler struct {
	ucase keyword.Usecase
}

func NewKeywordHandler(e *echo.Group, ucase keyword.Usecase) {
	handler := &handler{
		ucase,
	}
	e.DELETE("/keywords", handler.DeleteKeywords)
}

func (h *handler) DeleteKeywords(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	ids := []uint{}
	for _, idStr := range strings.Split(c.QueryParam("id"), ",") {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			e := errors.Wrap(errors.ErrCannotDeleteKeywords, []error{err})
			return c.JSON(getStatusCode(e), models.Response{Errors: []error{formatError(e)}})
		}
		ids = append(ids, uint(id))
	}
	keywords, err := h.ucase.Delete(ids...)
	if err != nil {
		return c.JSON(getStatusCode(err), models.Response{Errors: []error{formatError(err)}})
	}
	return c.JSON(http.StatusOK, models.Response{Data: keywords})
}

func getStatusCode(err error) int {
	e := errors.ToErrorModel(err)
	switch e.Message {
	case errors.ErrInvalidKeywordType,
		errors.ErrInvalidKeywordFor,
		errors.ErrInvalidKeywordValue:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func formatError(err error) error {
	e := errors.ToErrorModel(err)
	defaultMsg := utils.NewI18NMsg(e.Message, e.Message)
	localizer := _i18n.NewLocalizer()
	switch e.Message {
	default:
		e.Message = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID:      e.Message,
			DefaultMessage: defaultMsg,
		})
		return e
	}
}
