package authentication

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
	"gopkg.in/square/go-jose.v2"
)

const microsoftOpenIDConfigurationURL = "https://login.microsoftonline.com/common/.well-known/openid-configuration"

// OpenIDConfiguration contains the configuration for a servers OpenID endpoints
type OpenIDConfiguration struct {
	JWKSURI string `json:"jwks_uri"`
}

// GetMicrosoftOpenIDJWKS returns the JWKS for Microsoft authentication
func GetMicrosoftOpenIDJWKS() jose.JSONWebKeySet {
	configResp, err := http.Get(microsoftOpenIDConfigurationURL)
	if err != nil {
		log.Error().Err(err).Msg("error retrieving Microsoft OpenID Configuration")
		return jose.JSONWebKeySet{}
	}

	configBody, err := ioutil.ReadAll(configResp.Body)
	if err != nil {
		log.Error().Err(err).Msg("error reading Microsoft OpenID Configuration response body")
		return jose.JSONWebKeySet{}
	}

	if configResp.StatusCode != http.StatusOK {
		log.Error().Int("status", configResp.StatusCode).Bytes("body", configBody).Msg("error Microsoft OpenID Configuration returned an error")
		return jose.JSONWebKeySet{}
	}

	var config OpenIDConfiguration
	if err := json.Unmarshal(configBody, &config); err != nil {
		log.Error().Err(err).Msg("error parsing Microsoft OpenID Configuration response body")
		return jose.JSONWebKeySet{}
	}

	if config.JWKSURI == "" {
		log.Error().Err(err).Msg("error retrieving jwks uri from Microsoft OpenID Configuration response body")
		return jose.JSONWebKeySet{}
	}

	resp, err := http.Get(config.JWKSURI)
	if err != nil {
		log.Error().Err(err).Msg("error retrieving Microsoft OpenID JWKS keys")
		return jose.JSONWebKeySet{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("error reading Microsoft OpenID JWKS keys response body")
		return jose.JSONWebKeySet{}
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status", resp.StatusCode).Bytes("body", body).Msg("error Microsoft OpenID JWKS keys returned an error")
		return jose.JSONWebKeySet{}
	}

	var keys jose.JSONWebKeySet
	if err := json.Unmarshal(body, &keys); err != nil {
		log.Error().Err(err).Msg("error parsing Microsoft OpenID JWKS keys response body")
		return jose.JSONWebKeySet{}
	}

	return keys
}
