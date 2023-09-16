package config

import (
	"github.com/google/uuid"
	"github.com/hibare/DomainHQ/internal/constants"
	"github.com/hibare/GoCommon/v2/pkg/env"
	log "github.com/sirupsen/logrus"
)

type ServerConfig struct {
	ListenAddr string
	ListenPort int
}

type APIConfig struct {
	APIKeys []string
}

type WebFingerConfig struct {
	Resource string
	Domain   string
}

type DB struct {
	Username string
	Password string
	Host     string
	Port     int
	Name     string
}

type Config struct {
	Server    ServerConfig
	WebFinger WebFingerConfig
	DB        DB
	APIConfig APIConfig
}

var Current *Config

func LoadConfig() {

	env.Load()

	token := []string{
		uuid.New().String(),
	}

	Current = &Config{
		Server: ServerConfig{
			ListenAddr: env.MustString("LISTEN_ADDR", constants.DefaultAPIListenAddr),
			ListenPort: env.MustInt("LISTEN_PORT", constants.DefaultAPIListenPort),
		},
		WebFinger: WebFingerConfig{
			Resource: env.MustString("WEB_FINGER_RESOURCE", constants.DefaultWebFingerResource),
			Domain:   env.MustString("WEB_FINGER_DOMAIN", constants.DefaultWebFingerDomain),
		},
		DB: DB{
			Username: env.MustString("DB_USERNAME", ""),
			Password: env.MustString("DB_PASSWORD", ""),
			Host:     env.MustString("DB_HOST", constants.DefaultDBHost),
			Port:     env.MustInt("DB_PORT", constants.DefaultDBPort),
			Name:     env.MustString("DB_NAME", constants.DefaultDBName),
		},
		APIConfig: APIConfig{
			APIKeys: env.MustStringSlice("API_KEYS", token),
		},
	}

	log.Infof("WebFinger config: %+v", Current.WebFinger)
}
