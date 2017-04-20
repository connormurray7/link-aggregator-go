package main

import (
	"fmt"
	"log"

	"net/http"

	"github.com/spf13/viper"
)

func main() {
	viper := viper.New()
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
	addr := viper.GetString("address")
	port := ":" + viper.GetString("port")

	log.Printf("This is the addr %s and port %s\n", addr, port)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Println("Starting server...")
	http.ListenAndServe(port, nil)
}
