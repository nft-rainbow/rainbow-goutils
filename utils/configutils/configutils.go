package configutils

import (
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func MustLoad[T any](fs ...mapstructure.DecodeHookFunc) *T {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.AddConfigPath("..")     // optionally look for config in the working directory
	viper.AddConfigPath("../..")  // optionally look for config in the working directory
	return mustLoadViper[T](fs...)
}

func MustLoadByFile[T any](configPath string, fs ...mapstructure.DecodeHookFunc) *T {
	viper.SetConfigFile(configPath)
	return mustLoadViper[T]()
}

func mustLoadViper[T any](fs ...mapstructure.DecodeHookFunc) *T {
	var _config T
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalln(fmt.Errorf("fatal error config file: %w", err))
	}
	fmt.Printf("viper user config file: %v\n", viper.ConfigFileUsed())
	if err := viper.Unmarshal(&_config, func(dc *mapstructure.DecoderConfig) {
		dc.ErrorUnset = true
		dc.ErrorUnused = true
		dc.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			append(
				fs,
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.TextUnmarshallerHookFunc(),
			)...,
		)
	}); err != nil {
		panic(err)
	}
	return &_config
}
