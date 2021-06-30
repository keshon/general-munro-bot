// Package config provides methods for accesing config file in TOML format
package config

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	CORS struct {
		AllowOrigins string
		AllowMethods string
		AllowHeaders string
	}
	Kitsu struct {
		Debug          bool
		Hostname       string
		Email          string
		Password       string
		ListenHostname string
	}
	Bot struct {
		Debug          bool
		Token          string
		StateTimeout   int
		Webhook        bool
		Hostname       string
		ListenHostname string
		Language       string
	}
	Notification struct {
		IsEnabled              bool
		Threads                int
		PollDuration           int
		CommentTruncateAt      int
		SilentUpdate           bool
		SuppressUndefinedRoles bool
		AdminChatID            string
		ChatIDByRoles          []string
		DoneStatuses           []string
	}
	Backup struct {
		IsEnabled       bool
		Threads         int
		PollDuration    int
		LocalStorage    string
		IgnoreExtension []string
		S3              struct {
			AccessKey        string
			SecretKey        string
			BucketName       string
			Endpoint         string
			Region           string
			S3ForcePathStyle bool
			RootFolderName   string
		}
	}
}

func Read() Config {
	path := "conf.toml"
	if os.Getenv("TEST") == "true" {
		path = "c:/Users/SokolovIA/Dropbox/Projects/GitHub/general-munro-bot/conf.toml"
	}

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var config Config
	if err := toml.NewDecoder(f).Decode(&config); err != nil {
		panic(err)
	}

	return config
}
