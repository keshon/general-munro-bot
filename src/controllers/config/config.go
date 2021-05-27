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
		Hostname string
		Email    string
		Password string
		Debug    bool
	}
	Bot struct {
		Token          string
		StateTimeout   int
		Debug          bool
		Webhook        bool
		Hostname       string
		ListenHostname string
		Language       string
	}
	Credentials struct {
		AdminChatID   string
		ChatIDByRoles []string
	}
	Messaging struct {
		PollDuration      int
		CommentTruncateAt int
		SilentUpdate      bool
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
