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
	testDBUsername        = "test"
	testDBPassword        = "test"
)

func setenv() {
	os.Setenv("DOMAIN_HQ_LISTEN_ADDR", testAPIListenAddr)
	os.Setenv("DOMAIN_HQ_LISTEN_PORT", strconv.Itoa(testAPIListenPort))
	os.Setenv("DOMAIN_HQ_WEB_FINGER_DOMAIN", testWebFingerDomain)
	os.Setenv("DOMAIN_HQ_WEB_FINGER_RESOURCE", testWebFingerResource)
	os.Setenv("DOMAIN_HQ_DB_USERNAME", testDBUsername)
	os.Setenv("DOMAIN_HQ_DB_PASSWORD", testDBPassword)
	os.Setenv("DOMAIN_HQ_DB_HOST", "")
	os.Setenv("DOMAIN_HQ_DB_PORT", strconv.Itoa(constants.DefaultDBPort))
	os.Setenv("DOMAIN_HQ_DB_NAME", constants.DefaultDBName)
	os.Setenv("DOMAIN_HQ_API_KEYS", testAPIKeys)
	os.Setenv("DOMAIN_HQ_LOG_LEVEL", commonLogger.DefaultLoggerLevel)
	os.Setenv("DOMAIN_HQ_LOG_MODE", commonLogger.DefaultLoggerMode)
}

func unsetEnv() {
	os.Unsetenv("DOMAIN_HQ_LISTEN_ADDR")
	os.Unsetenv("DOMAIN_HQ_LISTEN_PORT")
	os.Unsetenv("DOMAIN_HQ_WEB_FINGER_DOMAIN")
	os.Unsetenv("DOMAIN_HQ_WEB_FINGER_RESOURCE")
	os.Unsetenv("DOMAIN_HQ_DB_USERNAME")
	os.Unsetenv("DOMAIN_HQ_DB_PASSWORD")
	os.Unsetenv("DOMAIN_HQ_DB_HOST")
	os.Unsetenv("DOMAIN_HQ_DB_PORT")
	os.Unsetenv("DOMAIN_HQ_DB_NAME")
	os.Unsetenv("DOMAIN_HQ_API_KEYS")
	os.Unsetenv("DOMAIN_HQ_LOG_LEVEL")
	os.Unsetenv("DOMAIN_HQ_LOG_MODE")
}

func TestEnvLoadedConfig(t *testing.T) {
	setenv()

	LoadConfig()

	assert.Equal(t, testAPIListenAddr, Current.Server.ListenAddr)
	assert.Equal(t, testAPIListenPort, Current.Server.ListenPort)
	assert.Equal(t, testWebFingerDomain, Current.WebFinger.Domain)
	assert.Equal(t, testWebFingerResource, Current.WebFinger.Resource)
	assert.Equal(t, testDBUsername, Current.DB.Username)
	assert.Equal(t, testDBPassword, Current.DB.Password)
	assert.Equal(t, []string{testAPIKeys}, Current.API.APIKeys)
	assert.Equal(t, commonLogger.DefaultLoggerLevel, Current.Logger.Level)
	assert.Equal(t, commonLogger.DefaultLoggerMode, Current.Logger.Mode)
	unsetEnv()
}

func TestDefaultConfig(t *testing.T) {
	unsetEnv()
	os.Setenv("DOMAIN_HQ_DB_USERNAME", testDBUsername)
	os.Setenv("DOMAIN_HQ_DB_PASSWORD", testDBPassword)

	LoadConfig()

	assert.Equal(t, constants.DefaultAPIListenAddr, Current.Server.ListenAddr)
	assert.Equal(t, constants.DefaultAPIListenPort, Current.Server.ListenPort)
	assert.Equal(t, constants.DefaultWebFingerDomain, Current.WebFinger.Domain)
	assert.Equal(t, constants.DefaultWebFingerResource, Current.WebFinger.Resource)
	assert.Equal(t, testDBUsername, Current.DB.Username)
	assert.Equal(t, testDBPassword, Current.DB.Password)
	assert.NotEmpty(t, Current.API.APIKeys)
	assert.Equal(t, commonLogger.DefaultLoggerLevel, Current.Logger.Level)
	assert.Equal(t, commonLogger.DefaultLoggerMode, Current.Logger.Mode)
}
