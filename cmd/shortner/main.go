package main

import (
	"flag"
	"fmt"
	"github.com/hex_microservice_template/pkg/adapters/driven"
	"github.com/hex_microservice_template/pkg/adapters/driving"
	"github.com/hex_microservice_template/pkg/core/ports/outbound"
	"github.com/hex_microservice_template/pkg/core/usecase"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// AppVersion is the application version
	AppVersion = "1.0"
	// Version is the git commit version (set by Makefile)
	Version = "none"
	// BuildTime application build time (set by Makefile)
	BuildTime = "none"

	version = flag.Bool("version", false, "print version string")

	appName = "hex_microservice_template"
)

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	fullVersion := fmt.Sprintf("%s-%s", AppVersion, Version)

	if *version {
		fmt.Printf("%s v%s (%s)\n", appName, fullVersion, BuildTime)
		flag.PrintDefaults()

		return
	}

	log.WithFields(log.Fields{
		"app":       appName,
		"version":   Version,
		"buildTime": BuildTime,
	}).Info("Starting up")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	server := driving.InitServer(driving.NewHandler(usecase.NewRedirectService(chooseRepo())))
	server.AppName = appName
	server.Version = Version
	server.BuildTime = BuildTime

	go func() {

		defer func() {
			if r := recover(); r != nil {
				log.WithField("routine", r).Info("Recovered in routine")
				return
			}
		}()

		if err := server.ListenAndServe(os.Getenv("HTTP_ADDRESS")); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start web server")
		}
		termChan <- syscall.SIGTERM
	}()

	select {
	case <-termChan:
		if server != nil {
			server.Close()
		}

		os.Exit(0)
	}
}

func chooseRepo() outbound.RedirectRepository {
	switch os.Getenv("DB") {
	case "REDIS":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := driven.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}
