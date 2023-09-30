package pgconnector

type Config struct {
	Host     string `mapstructure:"HOST"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
	Port     uint16 `mapstructure:"PORT"`
	DB       string `mapstructure:"DB"`
}
