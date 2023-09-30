package rsconnector

type Config struct {
	Host     string `mapstructure:"HOST"`
	Password string `mapstructure:"PASSWORD"`
	DB       int    `mapstructure:"DB"`
}
