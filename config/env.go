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

	// 1. Dit à Viper de lire les variables d'environnement système (Kubernetes)
	viper.AutomaticEnv()
	for _, key := range []string{
		"APP_ENV",
		"APP_PORT",
		"HOST_NAME",
		"API_KEY",
		"API_GROUP",
		"GIN_MODE",
		"ALLOWED_HOSTS",
		"MINIO_ENDPOINT",
		"MINIO_ACCESS_KEY_ID",
		"MINIO_ACCESS_KEY_SECRET",
		"MINIO_USE_SSL",
		"INSECURE_SKIP_VERIFY",
	} {
		if err := viper.BindEnv(key); err != nil {
			return err
		}
	}

	// 2. Tente de lire le fichier de config local (si présent)
	if err := viper.ReadInConfig(); err != nil {
		// Si le fichier n'est pas trouvé, on ne bloque pas (normal en production)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// 3. CORRECTIF CRUCIAL : On force le mapping des variables vers la structure Env
	if err := viper.Unmarshal(Env); err != nil {
		return err
	}

	return nil
}
