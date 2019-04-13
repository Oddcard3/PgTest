package config

import (
	"flag"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Init initialize config
func Init() {
	const defaultURL string = "postgres://postgres:postgres@localhost:5432/gotest?sslmode=disable"
	// default values
	viper.SetDefault("database_url", defaultURL)
	viper.SetDefault("config", "")

	//aliases
	//viper.RegisterAlias("database_url", "db.url")
	//viper.RegisterAlias("database_url", "dburl")

	// flags
	flag.String("dburl", defaultURL, "DB URL connection")
	flag.String("config", "", "path to config file")
	flag.Int("tenants", 1, "number of tenants")
	flag.Int("tables", 1, "number of tables")
	flag.Int("records", 1, "number of records")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// env
	viper.SetEnvPrefix("pgtest")
	viper.BindEnv("dburl")

	// config file
	if cfgFilePath := viper.GetString("config"); cfgFilePath != "" {
		viper.SetConfigFile(cfgFilePath)
	} else {
		viper.AddConfigPath("./")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Config isn't loaded: %s\n", err)
	}
}
