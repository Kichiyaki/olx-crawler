package models

type Config struct {
	Version string `json:"version" mapstructure:"version"`
	Colly   struct {
		Delay int `json:"delay"`
		Limit int `json:"limit"`
	} `mapstructure:"colly" json:"colly"`
	DiscordNotifications struct {
		ChannelID string `json:"channel_id" mapstructure:"channel_id"`
		Enabled   bool   `json:"enabled"`
		Token     string `json:"token"`
	} `mapstructure:"discord_notifications" json:"discord_notifications"`
	Port    uint     `json:"port"`
	Lang    string   `json:"lang" mapstructure:"lang"`
	Proxies []string `json:"proxies,omitempty" mapstructure:"proxies"`
	Debug   bool     `json:"debug" mapstructure:"debug"`
}
