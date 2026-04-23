package core

import "os"

const (
	MaxBackups    = 10
	HighLatencyMs = 1500
)

func homeDir() string {
	return os.Getenv("HOME")
}
