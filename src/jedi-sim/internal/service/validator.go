package service

import (
	"fmt"
	"jedi-sim/internal/model"
	"strconv"
	"strings"
)

// ValidateErrorCode validates whether the input error code exists in the errorCodes.json
func ValidateErrorCode(code string) (bool, string) {
	// Trim whitespace from input
	code = strings.TrimSpace(code)

	// Convert error code to integer
	errorCodeInt, err := strconv.Atoi(code)
	if err != nil {
		return false, fmt.Sprintf("Invalid error code format: %s", code)
	}

	// Check if error code exists in errorCodes.json
	_, ok := model.LoadErrorInfoFromJSON(errorCodeInt, "")
	if ok {
		return true, fmt.Sprintf("Error code %s is valid", code)
	}

	return false, fmt.Sprintf("Error code %s is not valid", code)
}
