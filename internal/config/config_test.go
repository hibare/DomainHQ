package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/hibare/DomainHQ/internal/constants"
	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
	"github.com/stretchr/testify/assert"
)

const (
	testAPIListenAddr     = "127.0.0.1"
	testAPIListenPort     = 10000
	testWebFingerDomain   = "example1.com"
	testWebFingerResource = "https://auth.example1.com"
	testAPIKeys           = "test-api-key"
)

func setenv() {
	os.Setenv("LISTEN_ADDR", testAPIListenAddr)
	os.Setenv("LISTEN_PORT", strconv.Itoa(testAPIListenPort))
	os.Setenv("WEB_FINGER_DOMAIN", testWebFingerDomain)
	os.Setenv("WEB_FINGER_RESOURCE", testWebFingerResource)
	os.Setenv("DB_USERNAME", "")
	os.Setenv("DB_PASSWORD", "")
	os.Setenv("DB_HOST", "")
	os.Setenv("DB_PORT", strconv.Itoa(constants.DefaultDBPort))
	os.Setenv("DB_NAME", constants.DefaultDBName)
	os.Setenv("API_KEYS", testAPIKeys)
	os.Setenv("LOG_LEVEL", commonLogger.DefaultLoggerLevel)
	os.Setenv("LOG_MODE", commonLogger.DefaultLoggerMode)
}

func unsetEnv() {
	os.Unsetenv("LISTEN_ADDR")
	os.Unsetenv("LISTEN_PORT")
	os.Unsetenv("WEB_FINGER_DOMAIN")
	os.Unsetenv("WEB_FINGER_RESOURCE")
	os.Unsetenv("DB_USERNAME")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("API_KEYS")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("LOG_MODE")
}

func TestEnvLoadedConfig(t *testing.T) {
	setenv()

	LoadConfig()

	assert.Equal(t, testAPIListenAddr, Current.Server.ListenAddr)
	assert.Equal(t, testAPIListenPort, Current.Server.ListenPort)
	assert.Equal(t, testWebFingerDomain, Current.WebFinger.Domain)
	assert.Equal(t, testWebFingerResource, Current.WebFinger.Resource)
	assert.Equal(t, []string{testAPIKeys}, Current.API.APIKeys)
	assert.Equal(t, commonLogger.DefaultLoggerLevel, Current.Logger.Level)
	assert.Equal(t, commonLogger.DefaultLoggerMode, Current.Logger.Mode)
	unsetEnv()
}

func TestDefaultConfig(t *testing.T) {
	unsetEnv()

	LoadConfig()

	assert.Equal(t, constants.DefaultAPIListenAddr, Current.Server.ListenAddr)
	assert.Equal(t, constants.DefaultAPIListenPort, Current.Server.ListenPort)
	assert.Equal(t, constants.DefaultWebFingerDomain, Current.WebFinger.Domain)
	assert.Equal(t, constants.DefaultWebFingerResource, Current.WebFinger.Resource)
	assert.NotEmpty(t, Current.API.APIKeys)
	assert.Equal(t, commonLogger.DefaultLoggerLevel, Current.Logger.Level)
	assert.Equal(t, commonLogger.DefaultLoggerMode, Current.Logger.Mode)
}
