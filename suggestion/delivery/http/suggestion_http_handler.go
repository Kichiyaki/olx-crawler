package http

import (
	"context"
	"net/http"
	"olx-crawler/errors"
	_i18n "olx-crawler/i18n"
	"olx-crawler/models"
	"olx-crawler/suggestion"
	"olx-crawler/utils"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/go-pg/urlstruct"
	"github.com/labstack/echo/v4"
)

type handler struct {
	ucase suggestion.Usecase
}

func NewSuggestionHandler(e *echo.Group, ucase suggestion.Usecase) {
	handler := &handler{
		ucase,
	}
	e.GET("/suggestions", handler.FetchSuggestions)
	e.GET("/suggestions/:id", handler.GetSuggestionByID)
	e.DELETE("/suggestions", handler.DeleteSuggestions)
}

func (h *handler) FetchSuggestions(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	filter := new(models.SuggestionFilter)
	urlstruct.Unmarshal(ctx, c.Request().URL.Query(), filter)
	pagination, err := h.ucase.Fetch(filter)
	if err != nil {
		return c.JSON(getStatusCode(err), models.Response{Error: formatError(err)})
	}
	return c.JSON(http.StatusOK, models.Response{Data: pagination})
}

func (h *handler) GetSuggestionByID(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		e := errors.Wrap(errors.ErrSuggestionNotFound, []error{err})
		return c.JSON(getStatusCode(e), models.Response{Error: formatError(e)})
	}
	suggestion, err := h.ucase.GetByID(uint(id))
	if err != nil {
		return c.JSON(getStatusCode(err), models.Response{Error: formatError(err)})
	}
	return c.JSON(http.StatusOK, models.Response{Data: suggestion})
}

func (h *handler) DeleteSuggestions(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	ids := []uint{}
	for _, idStr := range strings.Split(c.QueryParam("id"), ",") {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			e := errors.Wrap(errors.ErrCannotDeleteSuggestions, []error{err})
			return c.JSON(getStatusCode(e), models.Response{Error: formatError(e)})
		}
		ids = append(ids, uint(id))
	}
	suggestions, err := h.ucase.Delete(ids...)
	if err != nil {
		return c.JSON(getStatusCode(err), models.Response{Error: formatError(err)})
	}
	return c.JSON(http.StatusOK, models.Response{Data: suggestions})
}

func getStatusCode(err error) int {
	e := errors.ToErrorModel(err)
	switch e.Message {
	case errors.ErrSuggestionNotFound:
		return http.StatusNotFound
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
