package config

import (
	"github.com/hibare/DomainHQ/internal/constants"
	"github.com/hibare/DomainHQ/internal/env"
	log "github.com/sirupsen/logrus"
)

type ServerConfig struct {
	ListenAddr string
	ListenPort int
}

type WebFingerConfig struct {
	Resource string
	Domain   string
}

type Config struct {
	Server    ServerConfig
	WebFinger WebFingerConfig
}

var Current *Config

func LoadConfig() {

	env.Load()

	Current = &Config{
		Server: ServerConfig{
			ListenAddr: env.MustString("LISTEN_ADDR", constants.DefaultAPIListenAddr),
			ListenPort: env.MustInt("LISTEN_PORT", constants.DefaultAPIListenPort),
		},
		WebFinger: WebFingerConfig{
			Resource: env.MustString("WEB_FINGER_RESOURCE", constants.DefaultWebFingerResource),
			Domain:   env.MustString("WEB_FINGER_DOMAIN", constants.DefaultWebFingerDomain),
		},
	}

	log.Infof("WebFinger config: %+v", Current.WebFinger)
}
