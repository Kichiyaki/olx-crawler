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
	e.PATCH("/config/proxies", handler.UpdateProxies)
	e.PATCH("/config/colly", handler.UpdateColly)
	e.PATCH("/config/discord-notifications", handler.UpdateDiscordNotifications)
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

func (h *handler) UpdateProxies(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	proxies := []string{}
	if err := c.Bind(&proxies); err != nil {
		formatted := errors.Wrap(errors.ErrInvalidPayload, []error{err})
		return c.JSON(http.StatusBadRequest, models.Response{Errors: []error{formatError(formatted)}})
	}
	config, err := h.configManager.Config()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Errors: []error{formatError(err)}})
	}
	config.Proxies = proxies
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
