package server

import "github.com/spf13/viper"

// here for backwards compatibility
func init() {
	viper.SetDefault("addr", ":8080")
	viper.MustBindEnv("addr", "ADDRESS")
}
