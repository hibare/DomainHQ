package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/hibare/DomainHQ/internal/constants"
	"github.com/stretchr/testify/assert"
)

const (
	testAPIListenAddr     = "127.0.0.1"
	testAPIListenPort     = 10000
	testWebFingerDomain   = "example1.com"
	testWebFingerResource = "https://auth.example1.com"
)

func setenv() {
	os.Setenv("LISTEN_ADDR", testAPIListenAddr)
	os.Setenv("LISTEN_PORT", strconv.Itoa(testAPIListenPort))
	os.Setenv("WEB_FINGER_DOMAIN", testWebFingerDomain)
	os.Setenv("WEB_FINGER_RESOURCE", testWebFingerResource)
}

func unsetEnv() {
	os.Unsetenv("LISTEN_ADDR")
	os.Unsetenv("LISTEN_PORT")
	os.Unsetenv("WEB_FINGER_DOMAIN")
	os.Unsetenv("WEB_FINGER_RESOURCE")
}

func TestEnvLoadedConfig(t *testing.T) {
	setenv()

	LoadConfig()

	assert.Equal(t, testAPIListenAddr, Current.Server.ListenAddr)
	assert.Equal(t, testAPIListenPort, Current.Server.ListenPort)
	assert.Equal(t, testWebFingerDomain, Current.WebFinger.Domain)
	assert.Equal(t, testWebFingerResource, Current.WebFinger.Resource)

	unsetEnv()
}

func TestDefaultConfig(t *testing.T) {
	unsetEnv()

	LoadConfig()

	assert.Equal(t, constants.DefaultAPIListenAddr, Current.Server.ListenAddr)
	assert.Equal(t, constants.DefaultAPIListenPort, Current.Server.ListenPort)
	assert.Equal(t, constants.DefaultWebFingerDomain, Current.WebFinger.Domain)
	assert.Equal(t, constants.DefaultWebFingerResource, Current.WebFinger.Resource)
}
