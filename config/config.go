package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var c = loadConfiguration()

func loadConfiguration() *viper.Viper {
	cfg := viper.New()
	cfg.SetConfigName("config")
	cfg.SetConfigType("toml")
	cfg.AddConfigPath("/etc/tabula/")
	cfg.AddConfigPath("$HOME/.tabula")
	cfg.AddConfigPath(".")
	err := cfg.ReadInConfig()
	if err != nil {
		log.Fatalln(fmt.Errorf("fatal error config file: %w", err))
	}
	return cfg
}

func AppEnv() string {
	return c.GetString("app.env")
}

func DatabaseUser() string {
	return c.GetString("database.user")
}

func DatabasePassword() string {
	return c.GetString("database.password")
}

func DatabaseHost() string {
	return c.GetString("database.host")
}

func DatabasePort() int {
	return c.GetInt("database.port")
}

func DatabaseName() string {
	return c.GetString("database.dbname")
}

func DatabaseSSLMode() string {
	return c.GetString("database.ssl")
}

func ServerPort() int {
	return c.GetInt("server.port")
}
