package main

import (
	"encoding/json"
	"github.com/swarmbit/spacemesh-state-api/config"
	"log"
	"os"
)

func main() {
	StartServer(readConfig())
}

func readConfig() *config.Config {
	if len(os.Args) < 2 {
		log.Fatal("Usage: server <path to config>")
	}

	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configValues := config.Config{}
	err = decoder.Decode(&configValues)
	if err != nil {
		log.Fatal(err)
	}
	return &configValues
}
