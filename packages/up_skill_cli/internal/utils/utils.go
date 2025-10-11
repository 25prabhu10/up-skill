package utils

import (
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

func ValidateName(topicName string) bool {
	strLen := utf8.RuneCountInString(topicName)
	return 0 < strLen && strLen < 255
}

func ValidateNameWithTrim(topicName string) bool {
	return ValidateName(strings.TrimSpace(topicName))
}

func DoesDirectoryExistAndIsNotEmpty(name string) (bool, error) {
	if _, err := os.Stat(name); err == nil {
		dirEntries, err := os.ReadDir(name)
		if err != nil {
			log.Printf("could not read directory: %v", err)
			return false, err
		}
		if len(dirEntries) > 0 {
			return true, nil
		}
	}
	return false, nil
}

func NormalizeName(name string) string {
	r := regexp.MustCompile(`\W`)
	return strings.ToLower(r.ReplaceAllString(name, "_"))
}
