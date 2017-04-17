package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
	addr := viper.GetString("address")
	port := viper.GetString("port")

	fmt.Printf("This is the addr %s and port %s\n", addr, port)

}
