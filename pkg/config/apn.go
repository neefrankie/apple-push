package config

import "github.com/spf13/viper"

func MustAPNKey() string {
	key := viper.GetString("apple.apn_key")
	if key == "" {
		panic("empty receipt verification password")
	}

	return key
}