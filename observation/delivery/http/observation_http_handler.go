package http

import (
	"context"
	"net/http"
	"olx-crawler/errors"
	_i18n "olx-crawler/i18n"
	"olx-crawler/models"
	"olx-crawler/observation"
	"olx-crawler/utils"
	"strconv"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/go-pg/urlstruct"
	"github.com/labstack/echo/v4"
)

type handler struct {
	ucase observation.Usecase
}

type ObservationInput struct {
	Name        string           `json:"name"`
	URL         string           `json:"url"`
	Keywords    []models.Keyword `json:"keywords"`
	LastCheckAt time.Time        `json:"last_check_at"`
	Started     *bool            `json:"started"`
}

func (input ObservationInput) ToModel() models.Observation {
	o := models.Observation{
		Name:        input.Name,
		URL:         input.URL,
		Keywords:    input.Keywords,
		LastCheckAt: input.LastCheckAt,
	}
	if input.Started != nil {
		o.Started = input.Started
	}
	return o
}

func NewObservationHandler(e *echo.Group, ucase observation.Usecase) {
	handler := &handler{
		ucase,
	}
	e.POST("/observations", handler.CreateObservation)
	e.PATCH("/observations/:id", handler.UpdateObservation)
	e.DELETE("/observations", handler.DeleteObservations)
	e.GET("/observations", handler.FetchObservations)
	e.GET("/observations/:id", handler.GetObservationByID)
}

func (h *handler) FetchObservations(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	filter := new(models.ObservationFilter)
	urlstruct.Unmarshal(ctx, c.Request().URL.Query(), filter)
	pagination, err := h.ucase.Fetch(filter)
	if err != nil {
		return c.JSON(getStatusCode(err), models.Response{Errors: []error{formatError(err)}})
	}
	return c.JSON(http.StatusOK, models.Response{Data: pagination})
}

func (h *handler) GetObservationByID(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		e := errors.Wrap(errors.ErrObservationNotFound, []error{err})
		return c.JSON(getStatusCode(e), models.Response{Errors: []error{formatError(e)}})
	}
	observation, err := h.ucase.GetByID(uint(id))
	if err != nil {
		return c.JSON(getStatusCode(err), models.Response{Errors: []error{formatError(err)}})
	}
	return c.JSON(http.StatusOK, models.Response{Data: observation})
}

func (h *handler) CreateObservation(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	input := ObservationInput{}
	err := c.Bind(&input)
	if err != nil {
		e := errors.Wrap(errors.ErrCannotCreateObservation, []error{err})
		return c.JSON(getStatusCode(e), models.Response{Errors: []error{formatError(e)}})
	}
	observation := input.ToModel()
	err = h.ucase.Store(&observation)
	if err != nil {
		return c.JSON(getStatusCode(err), models.Response{Errors: []error{formatError(err)}})
	}
	return c.JSON(http.StatusOK, models.Response{Data: observation})
}

func (h *handler) UpdateObservation(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		e := errors.Wrap(errors.ErrObservationNotFound, []error{err})
		return c.JSON(getStatusCode(e), models.Response{Errors: []error{formatError(e)}})
	}
	input := ObservationInput{}
	err = c.Bind(&input)
	if err != nil {
		e := errors.Wrap(errors.ErrCannotUpdateObservation, []error{err})
		return c.JSON(getStatusCode(e), models.Response{Errors: []error{formatError(e)}})
	}
	o := input.ToModel()
	o.ID = uint(id)
	observation, err := h.ucase.Update(&o)
	if err != nil {
		return c.JSON(getStatusCode(err), models.Response{Errors: []error{formatError(err)}})
	}
	return c.JSON(http.StatusOK, models.Response{Data: observation})
}

func (h *handler) DeleteObservations(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	ids := []uint{}
	for _, idStr := range strings.Split(c.QueryParam("id"), ",") {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			e := errors.Wrap(errors.ErrCannotDeleteObservations, []error{err})
			return c.JSON(getStatusCode(e), models.Response{Errors: []error{formatError(e)}})
		}
		ids = append(ids, uint(id))
	}
	observations, err := h.ucase.Delete(ids...)
	if err != nil {
		return c.JSON(getStatusCode(err), models.Response{Errors: []error{formatError(err)}})
	}
	return c.JSON(http.StatusOK, models.Response{Data: observations})
}

func getStatusCode(err error) int {
	e := errors.ToErrorModel(err)
	switch e.Message {
	case errors.ErrInvalidObservationName,
		errors.ErrInvalidObservationURL,
		errors.ErrInvalidKeywordType,
		errors.ErrInvalidKeywordFor,
		errors.ErrInvalidKeywordValue:
		return http.StatusBadRequest
	case errors.ErrObservationNotFound:
		return http.StatusNotFound
	case errors.ErrObservationNameMustBeUnique,
		errors.ErrObservationURLMustBeUnique:
		return http.StatusConflict
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
