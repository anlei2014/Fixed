package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

// LoadErrorInfoFromJSON loads error info from errorCodes.json by code
func LoadErrorInfoFromJSON(code int, path string) (ErrorInfo, bool) {
	log.Printf("[LoadErrorInfoFromJSON] Loading error info for code: %d", code)
	log.Printf("[LoadErrorInfoFromJSON] Attempting to read from path: %s", path)

	// Try to open the specified file path
	file, err := os.Open(path)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Primary path failed: %v", err)
		// Fallback to current directory errorCodes.json
		fallbackPath := "errorCodes.json"
		log.Printf("[LoadErrorInfoFromJSON] Trying fallback path: %s", fallbackPath)
		file, err = os.Open(fallbackPath)
		if err != nil {
			log.Printf("[LoadErrorInfoFromJSON] Fallback path also failed: %v", err)
			return ErrorInfo{}, false
		}
		path = fallbackPath
	}
	defer file.Close()

	log.Printf("[LoadErrorInfoFromJSON] Successfully opened file: %s", path)

	// Read file content
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Failed to read file '%s': %v", path, err)
		return ErrorInfo{}, false
	}

	log.Printf("[LoadErrorInfoFromJSON] Read %d bytes from file", len(data))

	// Parse JSON data
	var errorMap map[string]ErrorInfo
	err = json.Unmarshal(data, &errorMap)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Failed to parse JSON from file '%s': %v", path, err)
		return ErrorInfo{}, false
	}

	// Convert code to string key and lookup
	key := strconv.Itoa(code)
	availableKeys := getMapKeys(errorMap)
	log.Printf("[LoadErrorInfoFromJSON] Available error codes: %v", availableKeys)
	log.Printf("[LoadErrorInfoFromJSON] Looking for key: '%s'", key)

	info, ok := errorMap[key]
	if ok {
		log.Printf("[LoadErrorInfoFromJSON] Found error info for code %d:", code)
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
func getMapKeys(m map[string]ErrorInfo) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// intToString converts integer to string (utility function)
func intToString(i int) string {
	return fmt.Sprintf("%d", i)
}
