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

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/go-pg/urlstruct"
	"github.com/labstack/echo/v4"
)

type handler struct {
	ucase observation.Usecase
}

type ObservationInput struct {
	Name     string            `json:"name"`
	URL      string            `json:"url"`
	OneOf    []models.OneOf    `json:"one_of"`
	Excluded []models.Excluded `json:"excluded"`
	Checked  []models.Checked  `json:"checked"`
}

func (o ObservationInput) ToModel() models.Observation {
	return models.Observation{
		Name:     o.Name,
		URL:      o.URL,
		OneOf:    o.OneOf,
		Excluded: o.Excluded,
		Checked:  o.Checked,
	}
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
		return c.JSON(getStatusCode(err), models.Response{Error: formatError(err)})
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
		return c.JSON(getStatusCode(e), models.Response{Error: formatError(e)})
	}
	observation, err := h.ucase.GetByID(uint(id))
	if err != nil {
		return c.JSON(getStatusCode(err), models.Response{Error: formatError(err)})
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
		return c.JSON(getStatusCode(e), models.Response{Error: formatError(e)})
	}
	observation := input.ToModel()
	err = h.ucase.Store(&observation)
	if err != nil {
		return c.JSON(getStatusCode(err), models.Response{Error: formatError(err)})
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
		return c.JSON(getStatusCode(e), models.Response{Error: formatError(e)})
	}
	input := ObservationInput{}
	err = c.Bind(&input)
	if err != nil {
		e := errors.Wrap(errors.ErrCannotUpdateObservation, []error{err})
		return c.JSON(getStatusCode(e), models.Response{Error: formatError(e)})
	}
	o := input.ToModel()
	o.ID = uint(id)
	observation, err := h.ucase.Update(&o)
	if err != nil {
		return c.JSON(getStatusCode(err), models.Response{Error: formatError(err)})
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
			return c.JSON(getStatusCode(e), models.Response{Error: formatError(e)})
		}
		ids = append(ids, uint(id))
	}
	observations, err := h.ucase.Delete(ids...)
	if err != nil {
		return c.JSON(getStatusCode(err), models.Response{Error: formatError(err)})
	}
	return c.JSON(http.StatusOK, models.Response{Data: observations})
}

func getStatusCode(err error) int {
	e := errors.ToErrorModel(err)
	switch e.Message {
	case errors.ErrInvalidExcludedFor,
		errors.ErrInvalidExcludedValue,
		errors.ErrInvalidObservationName,
		errors.ErrInvalidObservationURL,
		errors.ErrInvalidOneOfFor,
		errors.ErrInvalidOneOfValue:
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
