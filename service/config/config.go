package config

import "os"

const envLogFile = "LOG_FILE"

var defaultValue = map[string]string{envLogFile: "/home/apenkova/log2/golang_signup.log"}

//ReadLogFileConfig returns path to log file
func ReadLogFileConfig() string {
	if v := os.Getenv(envLogFile); v != "" {
		return v
	}
	return defaultValue[envLogFile]
}
