package config

type MemcachedConfig struct {
	Address           string `json:"address,omitempty"`
	Password          string `json:"password,omitempty"`
	Timeout           int    `json:"timeout,omitempty"`
	MaxIdleConnection int    `json:"maxIdleConnection,omitempty"`
}

type TelegramConfig struct {
	Host   string `json:"host,omitempty"`
	ChatID string `json:"chatId,omitempty"`
	Token  string `json:"token,omitempty"`
}

// Config holds the plugin configuration.
type Config struct {
	MaxRequestInWindow int             `json:"maxRequestInWindow,omitempty"`
	WindowTime         int             `json:"windowTime,omitempty"`
	Env                string          `json:"env,omitempty"`
	Memcached          MemcachedConfig `json:"memcached"`
	Telegram           TelegramConfig  `json:"telegram"`
}
