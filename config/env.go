package config

import (
	"github.com/spf13/viper"
)

type Environment struct {
	// Application config
	AppEnv   string `mapstructure:"APP_ENV"`
	AppPort  int    `mapstructure:"APP_PORT"`
	Hostname string `mapstructure:"HOST_NAME"`

	// API config
	ApiKey       string `mapstructure:"API_KEY"`
	ApiGroup     string `mapstructure:"API_GROUP"`
	GinMode      string `mapstructure:"GIN_MODE"`
	AllowedHosts string `mapstructure:"ALLOWED_HOSTS"`

	// Minio config
	MinioEndpoint        string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKeyID     string `mapstructure:"MINIO_ACCESS_KEY_ID"`
	MinioAccessKeySecret string `mapstructure:"MINIO_ACCESS_KEY_SECRET"`
	MinioUseSSL          bool   `mapstructure:"MINIO_USE_SSL"`
	InsecureSkipVerify   bool   `mapstructure:"INSECURE_SKIP_VERIFY"`
}

var Env = &Environment{}

// LoadEnv Loads environment variables.
func LoadEnv() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		err = viper.Unmarshal(Env)
	}
	return err
}
