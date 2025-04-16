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
	case "neo4j_url":
		return getEnvVarValue("NEO4J_URL")
	case "neo4j_username":
		return getEnvVarValue("NEO4J_USERNAME")
	case "neo4j_password":
		return getEnvVarValue("NEO4J_PASSWORD")

	}
	return ""
}

func getEnvVarValue(key string) string {
	val, found := os.LookupEnv(key)
	if !found {
		log.Println("Environment variable " + key + " not found")
	}
	return val
}
