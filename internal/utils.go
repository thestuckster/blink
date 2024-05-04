package internal

import (
	"fmt"
	"os"
)

func FindGameInstallationDirectory() (string, error) {
	paths := []string{
		"C:\\Program Files (x86)\\World of Warcraft",
		"C:\\Program Files\\World of Warcraft",
		"D:\\Games\\World of Warcraft",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("World of Warcraft installation directory not found")
}
