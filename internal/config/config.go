package config

import (
	"log"

	"github.com/google/uuid"
	"github.com/hibare/DomainHQ/internal/constants"
	"github.com/hibare/GoCommon/v2/pkg/env"
	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
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
			ListenAddr: env.MustString("DOMAIN_HQ_LISTEN_ADDR", constants.DefaultAPIListenAddr),
			ListenPort: env.MustInt("DOMAIN_HQ_LISTEN_PORT", constants.DefaultAPIListenPort),
		},
		WebFinger: WebFingerConfig{
			Resource: env.MustString("DOMAIN_HQ_WEB_FINGER_RESOURCE", constants.DefaultWebFingerResource),
			Domain:   env.MustString("DOMAIN_HQ_WEB_FINGER_DOMAIN", constants.DefaultWebFingerDomain),
		},
		DB: DBConfig{
			Username: env.MustString("DOMAIN_HQ_DB_USERNAME", ""),
			Password: env.MustString("DOMAIN_HQ_DB_PASSWORD", ""),
			Host:     env.MustString("DOMAIN_HQ_DB_HOST", constants.DefaultDBHost),
			Port:     env.MustInt("DOMAIN_HQ_DB_PORT", constants.DefaultDBPort),
			Name:     env.MustString("DOMAIN_HQ_DB_NAME", constants.DefaultDBName),
		},
		API: APIConfig{
			APIKeys: env.MustStringSlice("DOMAIN_HQ_API_KEYS", token),
		},
		Logger: LoggerConfig{
			Level: env.MustString("DOMAIN_HQ_LOG_LEVEL", commonLogger.DefaultLoggerLevel),
			Mode:  env.MustString("DOMAIN_HQ_LOG_MODE", commonLogger.DefaultLoggerMode),
		},
	}

	if Current.DB.Username == "" {
		log.Fatal("Error missing DB username")
	}

	if Current.DB.Password == "" {
		log.Fatal("Error missing DB password")
	}

	if !commonLogger.IsValidLogLevel(Current.Logger.Level) {
		log.Fatal("Error invalid logger level")
	}

	if !commonLogger.IsValidLogMode(Current.Logger.Mode) {
		log.Fatal("Error invalid logger mode")
	}

	commonLogger.InitLogger(&Current.Logger.Level, &Current.Logger.Mode)

}
