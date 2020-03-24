package http

import (
	"context"
	"net/http"
	"olx-crawler/config"
	"olx-crawler/errors"
	_i18n "olx-crawler/i18n"
	"olx-crawler/models"
	"olx-crawler/utils"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type handler struct {
	configManager config.Manager
}

func NewConfigHandler(e *echo.Group, configManager config.Manager) {
	handler := &handler{
		configManager,
	}
	e.GET("/config", handler.GetConfig)
	e.PATCH("/config/proxy", handler.UpdateProxyList)
	e.PATCH("/config/colly", handler.UpdateColly)
	e.PATCH("/config/discord-notifications", handler.UpdateDiscordNotifications)
	e.PATCH("/config/lang/:lang", handler.SetLanguage)
	e.PATCH("/config/debug/:value", handler.SetDebugMode)
}

func (h *handler) GetConfig(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	config, err := h.configManager.Config()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Errors: []error{formatError(err)}})
	}
	return c.JSON(http.StatusOK, models.Response{Data: config})
}

func (h *handler) UpdateProxyList(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	proxy := []string{}
	if err := c.Bind(&proxy); err != nil {
		formatted := errors.Wrap(errors.ErrInvalidPayload, []error{err})
		return c.JSON(http.StatusBadRequest, models.Response{Errors: []error{formatError(formatted)}})
	}
	config, err := h.configManager.Config()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Errors: []error{formatError(err)}})
	}
	config.Proxy = proxy
	if err := h.configManager.Save(config); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Errors: []error{formatError(err)}})
	}

	return c.JSON(http.StatusOK, models.Response{Data: config})
}

func (h *handler) UpdateColly(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	payload := models.Config{}
	if err := c.Bind(&payload.Colly); err != nil {
		formatted := errors.Wrap(errors.ErrInvalidPayload, []error{err})
		return c.JSON(http.StatusBadRequest, models.Response{Errors: []error{formatError(formatted)}})
	}
	config, err := h.configManager.Config()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Errors: []error{formatError(err)}})
	}
	if payload.Colly.Delay == 0 {
		payload.Colly.Delay = config.Colly.Delay
	}
	if payload.Colly.Limit == 0 {
		payload.Colly.Limit = config.Colly.Limit
	}
	config.Colly = payload.Colly
	if err := h.configManager.Save(config); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Errors: []error{formatError(err)}})
	}

	return c.JSON(http.StatusOK, models.Response{Data: config})
}

func (h *handler) UpdateDiscordNotifications(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	payload := make(map[string]interface{})
	if err := c.Bind(&payload); err != nil {
		formatted := errors.Wrap(errors.ErrInvalidPayload, []error{err})
		return c.JSON(http.StatusBadRequest, models.Response{Errors: []error{formatError(formatted)}})
	}
	config, err := h.configManager.Config()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Errors: []error{formatError(err)}})
	}
	if channelID, ok := payload["channel_id"].(string); ok {
		config.DiscordNotifications.ChannelID = channelID
	}
	if token, ok := payload["token"].(string); ok {
		config.DiscordNotifications.Token = token
	}
	if enabled, ok := payload["enabled"].(bool); ok {
		config.DiscordNotifications.Enabled = enabled
	}
	if err := h.configManager.Save(config); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Errors: []error{formatError(err)}})
	}

	return c.JSON(http.StatusOK, models.Response{Data: config})
}

func (h *handler) SetLanguage(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	lang := c.Param("lang")
	if err := _i18n.SetLanguage(lang); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.Response{Errors: []error{formatError(errors.Wrap(errors.ErrInvalidLanguage, []error{err}))}})
	}
	config, err := h.configManager.Config()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Errors: []error{formatError(err)}})
	}
	config.Lang = lang
	if err := h.configManager.Save(config); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Errors: []error{formatError(err)}})
	}

	return c.JSON(http.StatusOK, models.Response{Data: config})
}

func (h *handler) SetDebugMode(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	debug := c.Param("value")
	config, err := h.configManager.Config()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Errors: []error{formatError(err)}})
	}
	if debug == "true" {
		config.Debug = true
	} else if debug == "false" {
		config.Debug = false
	} else {
		return c.JSON(http.StatusBadRequest,
			models.Response{Errors: []error{formatError(errors.Wrap(errors.ErrInvalidPayload, []error{}))}})
	}
	if err := h.configManager.Save(config); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Errors: []error{formatError(err)}})
	}

	return c.JSON(http.StatusOK, models.Response{Data: config})
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
