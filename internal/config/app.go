package config

type App struct {
	QuotesLimit int    `env:"QUOTES_LIMIT"`
	Port        string `env:"PORT"`
	LogLevel    string `env:"LOG_LEVEL"`
}
