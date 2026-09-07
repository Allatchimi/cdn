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
    viper.SetConfigName("app")
    viper.SetConfigType("env")
    viper.AddConfigPath(".")
    viper.AddConfigPath("/app")

    // 1. DIT À VIPER DE LIRE LES VARIABLES D'ENVIRONNEMENT SYSTÈME
    viper.AutomaticEnv() 

    // 2. SI LE FICHIER N'EST PAS LÀ, CE N'EST PAS GRAVE EN PRODUCTION (KUBERNETES)
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            // Le fichier n'est pas là, mais c'est normal si les variables sont dans le Secret K8s
            return nil 
        }
        return err
    }
    return nil
}

