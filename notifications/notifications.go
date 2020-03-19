package notifications

import (
	"olx-crawler/config"
	"olx-crawler/errors"

	"github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"

	"github.com/bwmarrin/discordgo"
)

type Manager interface {
	Notify(string) error
	Close() error
}

type manager struct {
	configManager config.Manager
	discord       *discordgo.Session
	logrus        *logrus.Entry
}

func NewManager(configManager config.Manager) (Manager, error) {
	var firstErr error
	m := &manager{
		configManager: configManager,
		logrus:        logrus.WithField("package", "notifications"),
	}
	configManager.OnConfigChange(m.handleConfigChange)
	cfg, err := configManager.Config()
	if err != nil {
		return nil, err
	}
	if cfg.DiscordNotifications.Enabled && cfg.DiscordNotifications.Token != "" {
		var err error
		m.discord, err = discordgo.New("Bot " + cfg.DiscordNotifications.Token)
		if err != nil {
			firstErr = errors.Wrap(errors.ErrCannotConnectToDiscord, []error{err})
		}
	}
	return m, firstErr
}

func (m *manager) Notify(msg string) error {
	cfg, err := m.configManager.Config()
	if err != nil {
		return err
	}
	if m.discord != nil && cfg.DiscordNotifications.Enabled && cfg.DiscordNotifications.ChannelID != "" {
		_, err := m.discord.ChannelMessageSend(cfg.DiscordNotifications.ChannelID, msg)
		withFields := m.logrus.WithFields(logrus.Fields{
			"token":      m.discord.Token,
			"channel_id": cfg.DiscordNotifications.ChannelID,
			"msg":        msg,
		})
		if err != nil {
			withFields.Debugf("cannot send discord notification: %s", err.Error())
			return errors.Wrap(errors.ErrCannotSendDiscordNotification, []error{err})
		} else {
			withFields.Debug("send discord notification")
		}
	}
	return nil
}

func (m *manager) Close() error {
	if m.discord != nil {
		return m.discord.Close()
	}
	return nil
}

func (m *manager) handleConfigChange(fsnotify.Event) {
	var err error
	cfg, err := m.configManager.Config()
	if err != nil {
		m.Close()
		return
	}

	if m.discord == nil && cfg.DiscordNotifications.Enabled && cfg.DiscordNotifications.Token != "" {
		m.discord, err = discordgo.New("Bot " + cfg.DiscordNotifications.Token)

	} else if m.discord != nil && (!cfg.DiscordNotifications.Enabled || cfg.DiscordNotifications.Token == "") {
		m.discord.Close()
		m.discord = nil
	} else if m.discord != nil && m.discord.Token != "Bot "+cfg.DiscordNotifications.Token {
		m.discord.Close()
		m.discord, err = discordgo.New("Bot " + cfg.DiscordNotifications.Token)
	}

	withFields := m.logrus.WithFields(logrus.Fields{
		"token":      cfg.DiscordNotifications.Token,
		"channel_id": cfg.DiscordNotifications.ChannelID,
	})
	if err != nil {
		withFields.Debugf("cannot create new discordgo session: %s", err.Error())
	} else if m.discord != nil {
		withFields.Debug("new discordgo session")
	}
}
