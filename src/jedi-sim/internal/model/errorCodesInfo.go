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
	Z2             int    // Z2: Display Bitmap (bit 0-7) - 必须为 int
	Z3_Phase       int    // Z3: Phase (high 4 bits)
	Z3_Class       int    // Z3: Error Class (low 4 bits)
	Z4Z5_ErrorCode int    // Z4/Z5: Generator Error Code (little-endian format)
	Z6Z7_ErrorData int    // Z6/Z7: Error Data (little-endian format)
	Description    string // Description of the error
}

// LoadErrorInfoFromJSON loads error info from errorCodes.json by code
func LoadErrorInfoFromJSON(code int, path string) (ErrorInfo, bool) {
	file, err := os.Open(path)
	if err != nil {
		// 尝试用当前目录下的 errorCodes.json
		file, err = os.Open("errorCodes.json")
		if err != nil {
			log.Printf("[LoadErrorInfoFromJSON] Open file error: %v", err)
			return ErrorInfo{}, false
		}
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Read file error: %v", err)
		return ErrorInfo{}, false
	}
	var m map[string]ErrorInfo
	err = json.Unmarshal(data, &m)
	if err != nil {
		log.Printf("[LoadErrorInfoFromJSON] Unmarshal error: %v", err)
		return ErrorInfo{}, false
	}
	key := strconv.Itoa(code)
	log.Printf("[LoadErrorInfoFromJSON] All keys: %v, lookup key: %s", keysOfMap(m), key)
	info, ok := m[key]
	log.Printf("[LoadErrorInfoFromJSON] Lookup result: %v, found: %v", info, ok)
	return info, ok
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

func keysOfMap(m map[string]ErrorInfo) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return
}
