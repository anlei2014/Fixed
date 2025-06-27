package service

import (
	"fmt"
	"jedi-sim/internal/model"
	"log"
	"strconv"
	"strings"
)

// ValidateErrorCode validates whether the input error code exists in the ErrorMap
func ValidateErrorCode(code string) (bool, string) {
	code = strings.TrimSpace(code)
	log.Printf("[ValidateErrorCode] Validating code: %s", code)

	// Attempt to convert the error code to an integer
	errorCodeInt, err := strconv.Atoi(code)
	if err != nil {
		log.Printf("[ValidateErrorCode] Invalid format: %s is not a valid integer", code)
		return false, fmt.Sprintf("Invalid error code format: %s", code)
	}

	// Check if the error code exists in ErrorMap
	if _, ok := model.ErrorMap[errorCodeInt]; ok {
		log.Printf("[ValidateErrorCode] Code %s is valid", code)
		return true, fmt.Sprintf("Error code %s is valid", code)
	}

	log.Printf("[ValidateErrorCode] Code %s is not valid", code)
	return false, fmt.Sprintf("Error code %s is not valid", code)
}
