package service

import (
	"encoding/json"
	"fmt"
	"jedi-sim/internal/model"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// ValidateErrorCode validates whether the input error code exists in the errorCodes.json
func ValidateErrorCode(code string) (bool, string) {
	code = strings.TrimSpace(code)
	log.Printf("[ValidateErrorCode] Validating code: %s", code)

	// Attempt to convert the error code to an integer
	errorCodeInt, err := strconv.Atoi(code)
	if err != nil {
		log.Printf("[ValidateErrorCode] Invalid format: %s is not a valid integer", code)
		return false, fmt.Sprintf("Invalid error code format: %s", code)
	}

	// Check if the error code exists in errorCodes.json
	_, ok := model.LoadErrorInfoFromJSON(errorCodeInt, "src/jedi-sim/errorCodes.json")
	if ok {
		log.Printf("[ValidateErrorCode] Code %s is valid", code)
		return true, fmt.Sprintf("Error code %s is valid", code)
	}

	log.Printf("[ValidateErrorCode] Code %s is not valid", code)
	return false, fmt.Sprintf("Error code %s is not valid", code)
}

// 只用 model.ErrorInfo，不要重复定义
func LoadErrorInfoFromJSON(code int, path string) (model.ErrorInfo, bool) {
	log.Printf("[LoadErrorInfoFromJSON] Try path: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cwd, _ := os.Getwd()
		path = cwd + "/src/jedi-sim/errorCodes.json"
		log.Printf("[LoadErrorInfoFromJSON] Fallback path: %s", path)
	}
	file, err := os.Open(path)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Open file error: %v", err)
		return model.ErrorInfo{}, false
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Read file error: %v", err)
		return model.ErrorInfo{}, false
	}
	log.Printf("[LoadErrorInfoFromJSON] Raw JSON: %s", string(data))
	var m map[string]model.ErrorInfo
	err = json.Unmarshal(data, &m)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Unmarshal error: %v", err)
		return model.ErrorInfo{}, false
	}
	log.Printf("[LoadErrorInfoFromJSON] All keys: %v", keysOfMap(m))
	key := strconv.Itoa(code)
	log.Printf("[LoadErrorInfoFromJSON] Lookup key: %s", key)
	info, ok := m[key]
	log.Printf("[LoadErrorInfoFromJSON] Lookup result: %v, found: %v", info, ok)
	return info, ok
}

func keysOfMap(m map[string]model.ErrorInfo) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
