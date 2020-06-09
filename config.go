package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DBEngine string

	LoggerConfig string

	OAuthCallbackHost  string
	OAuthCallbackHTTPS bool
	OAuthClientID      string
	OAuthClientSecret  string
	OAuthProviderURL   string

	SecretKey string
}

func CollectConfig() (config Config) {
	var missingEnv []string

	// DB_ENGINE
	config.DBEngine = os.Getenv("DB_ENGINE")
	if config.DBEngine == "" {
		missingEnv = append(missingEnv, "DB_ENGINE")
	}

	// LOG_LEVEL
	var envLoggerLevel = os.Getenv("LOG_LEVEL")
	if envLoggerLevel == "" {
		config.LoggerConfig = "<root>=INFO"
	} else {
		config.LoggerConfig = fmt.Sprintf("<root>=%s", strings.ToUpper(envLoggerLevel))
	}

	// OAUTH_CALLBACK_HOST
	config.OAuthCallbackHost = os.Getenv("OAUTH_CALLBACK_HOST")
	if config.OAuthCallbackHost == "" {
		missingEnv = append(missingEnv, "OAUTH_CALLBACK_HOST")
	}

	// OAUTH_CALLBACK_HTTPS
	var envOAuthCallbackHTTPS = os.Getenv("OAUTH_CALLBACK_HTTPS")
	if strings.ToUpper(envOAuthCallbackHTTPS) == "FALSE" {
		config.OAuthCallbackHTTPS = false
	} else {
		config.OAuthCallbackHTTPS = true
	}

	// OAUTH_CLIENT_ID
	config.OAuthClientID = os.Getenv("OAUTH_CLIENT_ID")
	if config.OAuthClientID == "" {
		missingEnv = append(missingEnv, "OAUTH_CLIENT_ID")
	}

	// OAUTH_CLIENT_SECRET
	config.OAuthClientSecret = os.Getenv("OAUTH_CLIENT_SECRET")
	if config.OAuthClientSecret == "" {
		missingEnv = append(missingEnv, "OAUTH_CLIENT_SECRET")
	}

	// OAUTH_PROVIDER_URL
	config.OAuthProviderURL = os.Getenv("OAUTH_PROVIDER_URL")
	if config.OAuthProviderURL == "" {
		missingEnv = append(missingEnv, "OAUTH_PROVIDER_URL")
	}

	// SECRET_KEY
	config.SecretKey = os.Getenv("SECRET_KEY")
	if config.SecretKey == "" {
		missingEnv = append(missingEnv, "SECRET_KEY")
	}

	// Validation
	if len(missingEnv) > 0 {
		msg := fmt.Sprintf("Environment variables missing: %v", missingEnv)
		panic(fmt.Sprint(msg))
	}

	return
}
