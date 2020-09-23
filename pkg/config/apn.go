package config

import "github.com/spf13/viper"

func MustAPNKey() string {
	key := viper.GetString("apple.apn_key")
	if key == "" {
		panic("empty receipt verification password")
	}

	return key
}

type APNOption struct {
	AuthKey string `mapstructure:"apn_auth_key"`
	KeyID   string `mapstructure:"apn_key_id"`
	TeamID  string `mapstructure:"apn_team_id"`
}

func MustAPNOption() APNOption {
	var a APNOption

	err := viper.UnmarshalKey("apple", &a)
	if err != nil {
		panic(err)
	}

	return a
}

type APNTopic struct {
	Phone string `mapstructure:"phone"`
	Pad   string `mapstructure:"pad"`
}

func MustAPNTopic() APNTopic {
	var a APNTopic

	err := viper.UnmarshalKey("apple.apn_topics", &a)
	if err != nil {
		panic(err)
	}

	return a
}
