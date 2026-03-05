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

	Webhook *struct {
		Domain *string
		Secret *string
		Path   *string

		Host *string
		Port *int64
	}
}
