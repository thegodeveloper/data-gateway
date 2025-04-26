// Package common
// pkg/common/log.go
package common

import "log"

func Info(msg string, args ...any) {
	log.Printf("[INFO] "+msg, args...)
}

func Error(msg string, args ...any) {
	log.Printf("[ERROR] "+msg, args...)
}
