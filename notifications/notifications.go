package notifications

import (
	"olx-crawler/config"
	"olx-crawler/errors"

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
}

func NewManager(configManager config.Manager) (Manager, error) {
	var firstErr error
	m := &manager{
		configManager: configManager,
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
		if err != nil {
			return errors.Wrap(errors.ErrCannotSendDiscordNotification, []error{err})
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
	cfg, err := m.configManager.Config()
	if err != nil {
		m.Close()
		return
	}
	if m.discord == nil && cfg.DiscordNotifications.Enabled && cfg.DiscordNotifications.Token != "" {
		m.discord, _ = discordgo.New("Bot " + cfg.DiscordNotifications.Token)
	} else if m.discord != nil && (!cfg.DiscordNotifications.Enabled || cfg.DiscordNotifications.Token == "") {
		m.discord.Close()
		m.discord = nil
	} else if m.discord != nil && m.discord.Token != "Bot "+cfg.DiscordNotifications.Token {
		m.discord.Close()
		m.discord, _ = discordgo.New("Bot " + cfg.DiscordNotifications.Token)
	}
}
