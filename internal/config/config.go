package config

import (
	"github.com/google/uuid"
	"github.com/hibare/DomainHQ/internal/constants"
	"github.com/hibare/GoCommon/v2/pkg/env"
	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
	"github.com/rs/zerolog/log"
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

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Name     string
}

type LoggerConfig struct {
	Level string
	Mode  string
}

type Config struct {
	Server    ServerConfig
	WebFinger WebFingerConfig
	DB        DBConfig
	API       APIConfig
	Logger    LoggerConfig
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
		DB: DBConfig{
			Username: env.MustString("DB_USERNAME", ""),
			Password: env.MustString("DB_PASSWORD", ""),
			Host:     env.MustString("DB_HOST", constants.DefaultDBHost),
			Port:     env.MustInt("DB_PORT", constants.DefaultDBPort),
			Name:     env.MustString("DB_NAME", constants.DefaultDBName),
		},
		API: APIConfig{
			APIKeys: env.MustStringSlice("API_KEYS", token),
		},
		Logger: LoggerConfig{
			Level: env.MustString("LOG_LEVEL", commonLogger.DefaultLoggerLevel),
			Mode:  env.MustString("LOG_MODE", commonLogger.DefaultLoggerMode),
		},
	}

	if !commonLogger.IsValidLogLevel(Current.Logger.Level) {
		log.Fatal().Str("level", Current.Logger.Level).Msg("Error invalid logger level")
	}

	if !commonLogger.IsValidLogMode(Current.Logger.Mode) {
		log.Fatal().Str("mode", Current.Logger.Mode).Msg("Error invalid logger mode")
	}

	commonLogger.SetLoggingLevel(Current.Logger.Level)
	commonLogger.SetLoggingMode(Current.Logger.Mode)

	log.Info().Msgf("WebFinger config: %+v", Current.WebFinger)
}
