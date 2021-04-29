package main

import (
	"fmt"
	"os"
)

var env = make(map[string]string)

var envKeys = []string{
"PORT",
"AUTH_API_SERVICE_HOST",
"TASKS_DIRECTORY",
}

// configures the global env variable
func configureEnvVar(keys []string) {
	for _, k := range keys {
		env[k] = os.Getenv(k)
	}
}

// stores the working directory path under env["PWD"]
func configurePwd() error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get pwd: %v", err)
	}
	env["PWD"] = pwd
	return nil
}