package service

import (
	"fmt"
	"jedi-sim/internal/model"
	"log"
	"strings"
)

func ValidateErrorCode(code string) (bool, string) {
	code = strings.TrimSpace(code) // Remove leading/trailing spaces
	log.Printf("[ValidateErrorCode] Validating code: %s", code)
	for _, c := range model.ErrorCodes {
		if c == code {
			log.Printf("[ValidateErrorCode] Code %s is valid", code)
			return true, fmt.Sprintf("Error code %s is valid", code)
		}
	}
	log.Printf("[ValidateErrorCode] Code %s is not valid", code)
	return false, fmt.Sprintf("Error code %s is not valid", code)
}