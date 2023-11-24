package livechat

type Config struct {
	ClientID     string `env:"client_id"`
	ClientSecret string `env:"client_secret"`
	RedirectUrl  string `env:"redirect_url"`
}
