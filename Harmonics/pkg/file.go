package pkg

import (
	"os"
	"strings"
)

func AppendToFile(fileName string, response []string) bool {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return false
	}

	defer f.Close()

	var text strings.Builder
	for _, value := range response {
		text.WriteString(value)
	}

	if _, err = f.WriteString(text.String()); err != nil {
		return false
	}
	return true
}
