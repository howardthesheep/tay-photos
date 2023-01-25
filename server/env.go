package main

import (
	"io"
	"log"
	"os"
	"strings"
)

var _env map[string]string

// LoadEnvFile parses our .env file into memory to be used within our webserver
func LoadEnvFile() {
	// Read the .env file storing credentials
	envFile, err := os.Open("../.env")
	if err != nil {
		log.Printf("Error opening file" + err.Error())
		return
	}

	// Read contents into variable
	var contents []byte
	contents, err = io.ReadAll(envFile)
	if err != nil {
		log.Printf("Error reading file" + err.Error())
		return
	}

	// Parse key value pairs from .env file into Env
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		pair := strings.Split(line, "=")
		_env[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
	}
}

// GetEnv returns an instance of our loaded .env file parsed into key value pairs
func GetEnv() map[string]string {
	if _env == nil {
		LoadEnvFile()
	}

	return _env
}
