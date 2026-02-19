package util

import (
	"log"
	"os"
)

// EnsureDirs makes sure directories exist
func EnsureDirs(dirs ...string) {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatalf("failed to create directory: %s, error: %v", dir, err)
			}
			log.Printf("[main] created missing directory: %s", dir)
		}
	}
}
