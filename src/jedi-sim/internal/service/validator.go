package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jedi-sim/internal/model"
	"log"
	"os"
	"strconv"
	"strings"
)

// ValidateErrorCode validates whether the input error code exists in the errorCodes.json
func ValidateErrorCode(code string) (bool, string) {
	// Trim whitespace from input
	originalCode := code
	code = strings.TrimSpace(code)

	if originalCode != code {
		log.Printf("[ValidateErrorCode] Trimmed whitespace from '%s' to '%s'", originalCode, code)
	}

	log.Printf("[ValidateErrorCode] Validating error code: '%s'", code)

	// Convert error code to integer
	errorCodeInt, err := strconv.Atoi(code)
	if err != nil {
		log.Printf("[ValidateErrorCode] Conversion failed: '%s' is not a valid integer - %v", code, err)
		return false, fmt.Sprintf("Invalid error code format: %s", code)
	}

	log.Printf("[ValidateErrorCode] Successfully converted '%s' to integer: %d", code, errorCodeInt)

	// Check if error code exists in errorCodes.json
	log.Printf("[ValidateErrorCode] Checking existence in error codes database...")
	_, ok := model.LoadErrorInfoFromJSON(errorCodeInt, "src/jedi-sim/errorCodes.json")
	if ok {
		log.Printf("[ValidateErrorCode] Validation successful: error code %d exists in database", errorCodeInt)
		return true, fmt.Sprintf("Error code %s is valid", code)
	}

	log.Printf("[ValidateErrorCode] Validation failed: error code %d not found in database", errorCodeInt)
	return false, fmt.Sprintf("Error code %s is not valid", code)
}

// LoadErrorInfoFromJSON loads error info from JSON file with fallback paths
// This function provides a more robust path resolution for error codes
func LoadErrorInfoFromJSON(code int, path string) (model.ErrorInfo, bool) {
	log.Printf("[LoadErrorInfoFromJSON] Loading error info for code: %d", code)
	log.Printf("[LoadErrorInfoFromJSON] Primary path: %s", path)

	// Check if specified path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("[LoadErrorInfoFromJSON] Primary path does not exist: %s", path)
		// Fallback to current working directory
		cwd, _ := os.Getwd()
		fallbackPath := cwd + "/src/jedi-sim/errorCodes.json"
		log.Printf("[LoadErrorInfoFromJSON] Using fallback path: %s", fallbackPath)
		path = fallbackPath
	} else {
		log.Printf("[LoadErrorInfoFromJSON] Primary path exists: %s", path)
	}

	// Open file
	file, err := os.Open(path)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Failed to open file '%s': %v", path, err)
		return model.ErrorInfo{}, false
	}
	defer file.Close()

	log.Printf("[LoadErrorInfoFromJSON] Successfully opened file: %s", path)

	// Read file content
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Failed to read file '%s': %v", path, err)
		return model.ErrorInfo{}, false
	}

	log.Printf("[LoadErrorInfoFromJSON] Read %d bytes from file", len(data))

	// Log raw JSON for debugging (truncated if too long)
	jsonPreview := string(data)
	if len(jsonPreview) > 200 {
		jsonPreview = jsonPreview[:200] + "..."
	}
	log.Printf("[LoadErrorInfoFromJSON] JSON preview: %s", jsonPreview)

	// Parse JSON data
	var errorMap map[string]model.ErrorInfo
	err = json.Unmarshal(data, &errorMap)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Failed to parse JSON from file '%s': %v", path, err)
		return model.ErrorInfo{}, false
	}

	// Log all available keys
	availableKeys := getMapKeys(errorMap)
	log.Printf("[LoadErrorInfoFromJSON] Found %d error codes in database: %v", len(availableKeys), availableKeys)

	// Convert code to string key and lookup
	key := strconv.Itoa(code)
	log.Printf("[LoadErrorInfoFromJSON] Looking for key: '%s'", key)

	info, ok := errorMap[key]
	if ok {
		log.Printf("[LoadErrorInfoFromJSON] Successfully found error info for code %d", code)
		log.Printf("[LoadErrorInfoFromJSON] Error details:")
		log.Printf("[LoadErrorInfoFromJSON]   - Z0 (Generator Status): %d", info.Z0)
		log.Printf("[LoadErrorInfoFromJSON]   - Z1 (Simplified Error Code): %d", info.Z1)
		log.Printf("[LoadErrorInfoFromJSON]   - Z2 (Display Bitmap): %d", info.Z2)
		log.Printf("[LoadErrorInfoFromJSON]   - Z3_Phase: %d", info.Z3_Phase)
		log.Printf("[LoadErrorInfoFromJSON]   - Z3_Class: %d", info.Z3_Class)
		log.Printf("[LoadErrorInfoFromJSON]   - Z4Z5_ErrorCode: %d", info.Z4Z5_ErrorCode)
		log.Printf("[LoadErrorInfoFromJSON]   - Z6Z7_ErrorData: %d", info.Z6Z7_ErrorData)
		log.Printf("[LoadErrorInfoFromJSON]   - Description: '%s'", info.Description)
	} else {
		log.Printf("[LoadErrorInfoFromJSON] Error code %d not found in database", code)
	}

	return info, ok
}

// getMapKeys returns all keys from a map as a slice
func getMapKeys(m map[string]model.ErrorInfo) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
