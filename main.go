package main

import (
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"net/url"
	"os"
	"os/signal"
	"phsite/models"
	"phsite/web"
	"syscall"
)

var logger *loggo.Logger


func main() {
	// Collect Config
	cfg := CollectConfig()

	// Init Logging
	newLogger := loggo.GetLogger("main")
	logger = &newLogger

	err := loggo.ConfigureLoggers(cfg.LoggerConfig)
	if err != nil {
		logger.Errorf("Error configurting Logger: %s", err.Error())
		return
	}

	_, err = loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
	if err != nil {
		logger.Errorf("Error configurting Color Logger: %s", err.Error())
		return
	}

	logger.Infof("Module Init")
	// Init DB
	err = models.Init(cfg.DBEngine)
	if err !=nil {
		panic(err)
	}
	defer models.Close()

	// Init Web
	callbackURL := &url.URL{
		Scheme: "https",
		Host: cfg.OAuthCallbackHost,
		Path: "/oauth/callback",
	}
	if !cfg.OAuthCallbackHTTPS {callbackURL.Scheme = "http"}

	err = web.Init(cfg.SecretKey, cfg.OAuthProviderURL, cfg.OAuthClientID, cfg.OAuthClientSecret, callbackURL.String())
	if err !=nil {
		panic(err)
	}
	defer web.Close()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	logger.Infof("%s", <-nch)

}