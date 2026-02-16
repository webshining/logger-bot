package settings

type Config struct {
	Database DatabaseConfig `koanf:"db"`
	Bot      BotConfig
}

type DatabaseConfig struct {
	Driver string
	Url    string
}

type BotConfig struct {
	Token string
	Admin int64
}
