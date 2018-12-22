package main

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

var (
	testKey                         = "OoWee3ri1zeadethoopha1oo"
	testExpiredAccessToken          = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiMzUxNzc1IiwiYWxsb3dlZF9jaWRycyI6WyIwLjAuMC4wLzAiXSwiZXhwaXJlcyI6MTQ5MTU2MTM3MSwicHJvZmlsZV9pZCI6IjM1Mjc5NyIsInVzZXJfaWQiOiIzNTE3NzUiLCJ2IjoxfQ.7zV8Wzs55pujgdmY-zZhuBm7Z2GGSU4MFhEVOzRoDS8"
	testInvalidSignatureAccessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiNDc0MjU5IiwiYWxsb3dlZF9jaWRycyI6WyIwLjAuMC4wLzAiXSwiZXhwaXJlcyI6MTQ4NDA2MzAzNSwicHJvZmlsZV9pZCI6IjQ3NTI4MiIsInVzZXJfaWQiOiI0NzQyNTkiLCJ2IjoxfQ.EgP75wq7nUQZBIjQbIHfB1MlV6azlLyBhK8IskoERk4"
	testBadAccessToken              = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbGciOiJINiIsInR5cCI6IkpXVCJ9.EgP75wq7nUQZBIjQbIHfB1MlV6azlLyBhK8IskoERk4"
)

func makeAccessToken(signKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Second * 60).Unix()
	claims["sub"] = "351775"

	return token.SignedString([]byte(signKey))
}

func TestAutenticationJwt(t *testing.T) {
	token, err := makeAccessToken(testKey)
	assert.NoError(t, err)

	a, err := PluginInit("./testdata/config.yml")
	assert.NoError(t, err)

	user, err := a.Authenticate(token)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestAutenticationExpiredJwt(t *testing.T) {
	a, err := PluginInit("./testdata/config.yml")
	assert.NoError(t, err)

	user, err := a.Authenticate(testExpiredAccessToken)
	assert.Error(t, err)
	assert.Nil(t, user)

	//assert.Equal(t, errCredentialsExpired, err)
}

func TestAutenticationInvalidSignatureJwt(t *testing.T) {
	a, err := PluginInit("./testdata/config.yml")
	assert.NoError(t, err)

	user, err := a.Authenticate(testInvalidSignatureAccessToken)
	assert.Error(t, err)
	assert.Nil(t, user)

	assert.Equal(t, "Bad Credentials", err.Error())
}

func TestAutenticationBadJwt(t *testing.T) {
	a, err := PluginInit("./testdata/config.yml")
	assert.NoError(t, err)

	user, err := a.Authenticate(testBadAccessToken)
	assert.Error(t, err)
	assert.Nil(t, user)
}
