package babble

import (
	"crypto/ecdsa"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/andrecronje/babble/src/node"
	"github.com/andrecronje/babble/src/proxy"
	"github.com/sirupsen/logrus"
)

const DefaultKeyfile = "priv_key"

type BabbleConfig struct {
	NodeConfig node.Config `mapstructure:",squash"`

	DataDir     string `mapstructure:"datadir"`
	BindAddr    string `mapstructure:"listen"`
	ServiceAddr string `mapstructure:"service-listen"`
	MaxPool     int    `mapstructure:"max-pool"`
	Store       bool   `mapstructure:"store"`
	LogLevel    string `mapstructure:"log"`
	Moniker     string `mapstructure:"moniker"`

	LoadPeers bool
	Proxy     proxy.AppProxy
	Key       *ecdsa.PrivateKey
	Logger    *logrus.Logger
}

func NewDefaultConfig() *BabbleConfig {

	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	nodeConfig := *node.DefaultConfig()
	nodeConfig.Logger = logger

	config := &BabbleConfig{
		NodeConfig: nodeConfig,
		DataDir:    DefaultDataDir(),
		BindAddr:   ":1337",
		MaxPool:    2,
		Store:      false,
		LoadPeers:  true,
		Proxy:      nil,
		Key:        nil,
		Logger:     logger,
	}

	return config
}

func (c *BabbleConfig) BadgerDir() string {
	return filepath.Join(c.DataDir, "badger_db")
}

func (c *BabbleConfig) Keyfile() string {
	return filepath.Join(c.DataDir, DefaultKeyfile)
}

func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := HomeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, ".babble")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "BABBLE")
		} else {
			return filepath.Join(home, ".babble")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func HomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

func LogLevel(l string) logrus.Level {
	switch l {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}
