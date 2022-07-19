package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	HTTP        HTTPConfig
	Log         LogConfig
	RedisConfig RedisConfig
}

type HTTPConfig struct {
	Address string
	Domain  string
}

type LogConfig struct {
	Level      string
	Format     string
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
}

type RedisConfig struct {
	DB         int
	MaxRetries int
	Expiration time.Duration
}

func InitLogging(lc *LogConfig) {
	switch lc.Format {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		fallthrough
	case "text":
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	switch lc.Level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		fallthrough
	case "info":
		log.SetLevel(log.InfoLevel)
	}
}
