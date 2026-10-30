package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	env "github.com/caarlos0/env/v11"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/zbot"
)

// environmentVariables are the process-level settings not covered by the configuration file.
type environmentVariables struct {
	ConfigurationFile string `env:"ZBOT_CONFIG_FILE" envDefault:"./zbot.conf"`
	LogLevel          string `env:"ZBOT_LOG_LEVEL" envDefault:"info"`
}

// // logLevel holds the current log level and can be inspected in tests
// var logLevel = new(slog.LevelVar)

// logFileName is the fixed name of the file that receives a copy of every
// log line in addition to the standard output.
const logFileName = "zbot.log"

// setupLogLevel configures the default slog logger for the given level name
// ("debug", "info" or "error"). Unknown values fall back to info.
func setupLogLevel(level string) *slog.LevelVar {

	logLevel := new(slog.LevelVar)
	switch level {
	case "debug":
		logLevel.Set(slog.LevelDebug)
	case "error":
		logLevel.Set(slog.LevelError)
	default:
		logLevel.Set(slog.LevelInfo)
	}
	return logLevel
}

// setupLog configures the default slog logger to write to both stdout and a
// file, with the given log level.
func setupLog(level string) {
	logLevel := setupLogLevel(level)

	var w io.Writer = os.Stdout
	f, err := os.OpenFile(logFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open %s for logging: %v (logging to stdout only)\n", logFileName, err)
	} else {
		w = io.MultiWriter(os.Stdout, f)
	}

	h := slog.NewTextHandler(w, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(h))
}

// setup reads the configuration file and builds a ready-to-run bot from it.
func setup(configurationFile string) (*zbot.Bot, error) {
	conf, err := zbot.LoadConfiguration(configurationFile)
	if err != nil {
		return nil, fmt.Errorf("read configuration: %w", err)
	}

	cfg, err := zbot.LoadConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	cfg.LogFileName = logFileName

	database, err := db.New(conf.Db)
	if err != nil {
		return nil, fmt.Errorf("setup database: %w", err)
	}

	return zbot.New(cfg, database)
}

// run is the real entry point. It returns an error instead of exiting so that
// main is the only place deciding the process exit code.
func run(ctx context.Context) error {
	envVars := environmentVariables{}
	if err := env.Parse(&envVars); err != nil {
		return fmt.Errorf("parse environment: %w", err)
	}

	setupLog(envVars.LogLevel)
	slog.Debug("environment", "cfg", fmt.Sprintf("%+v", envVars))

	bot, err := setup(envVars.ConfigurationFile)
	if err != nil {
		return err
	}

	return bot.Run(ctx)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Error("zbot stopped", "err", err)
		os.Exit(1)
	}
}
