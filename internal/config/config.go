package config

import (
	"log"
	"os"
	"strings"
)

var address string = ":8080"

func GetConfigValue(key string) string {
	switch strings.ToLower(key) {
	case "address":
		return address
	default:
		return getEnvVarValue(strings.ToUpper(key))
	}
}

func getEnvVarValue(key string) string {
	val, found := os.LookupEnv(key)
	if !found {
		log.Println("Environment variable " + key + " not found")
	}
	return val
}
