package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// ErrorInfo defines all error-related fields in a CAN message
// Z2 字段类型改为 int，保证与 JSON 一致
type ErrorInfo struct {
	Z0             int    // Z0: Generator Status (e.g., 0x06 = Error condition)
	Z1             int    // Z1: Simplified Error Code (e.g., 30 = Simplified error code 30)
	Z2             int    // Z2: Display Bitmap (bit 0-7) - must be int
	Z3_Phase       int    // Z3: Phase (high 4 bits)
	Z3_Class       int    // Z3: Error Class (low 4 bits)
	Z4Z5_ErrorCode int    // Z4/Z5: Generator Error Code (little-endian format)
	Z6Z7_ErrorData int    // Z6/Z7: Error Data (little-endian format)
	Description    string // Description of the error
}

// getConfigPath returns the correct path to the errorCodes.json file
func getConfigPath() string {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("[getConfigPath] Failed to get current working directory: %v", err)
		return "config/errorCodes.json"
	}

	// Try multiple possible paths relative to current working directory
	possiblePaths := []string{
		filepath.Join(cwd, "config", "errorCodes.json"),
		filepath.Join(cwd, "src", "jedi-sim", "config", "errorCodes.json"),
		filepath.Join(cwd, "..", "config", "errorCodes.json"),
		filepath.Join(cwd, "..", "..", "config", "errorCodes.json"),
		"config/errorCodes.json",
		"src/jedi-sim/config/errorCodes.json",
	}

	for _, path := range possiblePaths {
		log.Printf("[getConfigPath] Trying path: %s", path)
		if _, err := os.Stat(path); err == nil {
			log.Printf("[getConfigPath] Found config file at: %s", path)
			return path
		}
	}

	log.Printf("[getConfigPath] No config file found, using default: config/errorCodes.json")
	return "config/errorCodes.json"
}

// LoadErrorInfoFromJSON loads error info from errorCodes.json by code
func LoadErrorInfoFromJSON(code int, path string) (ErrorInfo, bool) {
	// If path is empty, try to find the config file automatically
	if path == "" {
		path = getConfigPath()
	}

	log.Printf("[LoadErrorInfoFromJSON] Loading error info for code %d from path: %s", code, path)

	// Try to open the specified file path
	file, err := os.Open(path)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Failed to open file '%s': %v", path, err)
		return ErrorInfo{}, false
	}
	defer file.Close()

	// Read file content
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Failed to read file '%s': %v", path, err)
		return ErrorInfo{}, false
	}

	log.Printf("[LoadErrorInfoFromJSON] Successfully read %d bytes from file", len(data))

	// Parse JSON data
	var errorMap map[string]ErrorInfo
	err = json.Unmarshal(data, &errorMap)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Failed to parse JSON from file '%s': %v", path, err)
		return ErrorInfo{}, false
	}

	// Convert code to string key and lookup
	key := strconv.Itoa(code)
	log.Printf("[LoadErrorInfoFromJSON] Looking for key: '%s'", key)

	// Log all available keys for debugging
	availableKeys := make([]string, 0, len(errorMap))
	for k := range errorMap {
		availableKeys = append(availableKeys, k)
	}
	log.Printf("[LoadErrorInfoFromJSON] Available keys: %v", availableKeys)

	info, ok := errorMap[key]
	if ok {
		log.Printf("[LoadErrorInfoFromJSON] Found error info for code %d", code)
	} else {
		log.Printf("[LoadErrorInfoFromJSON] Error code %d not found in database", code)
	}

	return info, ok
}

// intToString converts integer to string (utility function)
func intToString(i int) string {
	return fmt.Sprintf("%d", i)
}
