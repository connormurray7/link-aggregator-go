package main

import (
	"fmt"
	"log"

	"net/http"

	"github.com/connormurray7/link-aggregator-go/linkagg"
	"github.com/spf13/viper"
)

func main() {
	viper := viper.New()
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}

	l := linkagg.NewServer(viper)
	http.HandleFunc("/search", l.Handle)

	addr := viper.GetString("address")
	port := ":" + viper.GetString("port")

	log.Printf("This is the addr %s and port %s\n", addr, port)

	fs := http.FileServer(http.Dir("run-locally"))
	http.Handle("/", fs)

	log.Println("Starting server...")
	http.ListenAndServe(port, nil)
}
