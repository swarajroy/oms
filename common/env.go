package common

import (
	"log"
	"syscall"
)

func EnvString(key, fallback string) string {
	log.Printf("key = %s, fallback = %s\n", key, fallback)
	if val, ok := syscall.Getenv(key); ok {
		return val
	}
	return fallback
}
