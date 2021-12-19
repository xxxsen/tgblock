package bot

type Option func(c *TGBotConfig)

type TGBotConfig struct {
	Token  string `json:"token"`
	ChatID int64  `json:"chatid"`
}

func WithToken(token string) Option {
	return func(c *TGBotConfig) {
		c.Token = token
	}
}

func WithChatId(chatid int64) Option {
	return func(c *TGBotConfig) {
		c.ChatID = chatid
	}
}
