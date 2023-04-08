package config

import (
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
			ListenAddr: env.MustString("LISTEN_ADDR", "0.0.0.0"),
			ListenPort: env.MustInt("LISTEN_PORT", 5000),
		},
		WebFinger: WebFingerConfig{
			Resource: env.MustString("WEB_FINGER_RESOURCE", "https://auth.example.com"),
			Domain:   env.MustString("WEB_FINGER_DOMAIN", "example.com"),
		},
	}

	log.Infof("WebFinger config: %+v", Current.WebFinger)
}
