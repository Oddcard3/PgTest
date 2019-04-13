package main

import (
	"pgtest/config"

	"github.com/spf13/viper"
)

func main() {
	config.Init()
	// TestDB()
	CreateTenants(viper.GetInt("tenants"))
}
