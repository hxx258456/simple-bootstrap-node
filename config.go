package main

import "github.com/spf13/viper"

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	prvKey = viper.GetString("privKey")
}

func WriteConfig() {
	viper.Set("privKey", prvKey)
	if err := viper.WriteConfig(); err != nil {
		panic(err)
	}
}
