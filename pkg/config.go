package BookQuest

import "github.com/spf13/viper"

var Configuration = Config{}

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DatabasePath string `mapstructure:"DATABASE_PATH"`
	OIDCServer   string `mapstructure:"OIDC_SERVER"`
	ClientID     string `mapstructure:"OIDC_CLIENT_ID"`
	SecretKey    string `mapstructure:"SECRET__KEY"`
	RedirectURI  string `mapstructure:"OIDC_REDIRECT_URI"`
	ListenPort   string `mapstructure:"LISTEN_PORT"`
}

func SetDefaults() {
	if Configuration.ListenPort == "" {
		Configuration.ListenPort = "8090"
	}
}

func LoadConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&Configuration)
	SetDefaults()
	return err
}
