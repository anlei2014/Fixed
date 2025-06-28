package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jedi-sim/internal/model"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// AddErrorCodeHandler handles adding new error code definitions to the JSON database
func AddErrorCodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[AddErrorCodeHandler] Processing %s request to %s", r.Method, r.URL.Path)

	// Validate HTTP method
	if err := validateHTTPMethod(r, "POST"); err != nil {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	errorCodeData, err := parseRequestBody(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("[AddErrorCodeHandler] Parsed error code data: Z4Z5_ErrorCode=%d, Description='%s'",
		errorCodeData.Z4Z5_ErrorCode, errorCodeData.Description)

	// Save error code data
	if err := saveErrorCodeData(errorCodeData); err != nil {
		// Check if it's a duplicate error code error
		if httpErr, ok := err.(*httpError); ok && httpErr.message != "" {
			// For duplicate error codes, return 400 Bad Request
			if strings.Contains(httpErr.message, "already exists") {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
		}
		// For other errors, return 500 Internal Server Error
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return success response
	respondWithAddSuccess(w, "Error code data saved successfully")
}

// parseRequestBody parses the JSON request body into ErrorInfo struct
func parseRequestBody(r *http.Request) (model.ErrorInfo, error) {
	var errorCodeData model.ErrorInfo
	err := json.NewDecoder(r.Body).Decode(&errorCodeData)
	if err != nil {
		log.Printf("[AddErrorCodeHandler] JSON parsing failed: %v", err)
		return model.ErrorInfo{}, &httpError{message: "Invalid JSON data"}
	}

	log.Printf("[AddErrorCodeHandler] JSON request body parsed successfully")
	return errorCodeData, nil
}

// saveErrorCodeData reads existing data, adds new entry, and writes back to file
func saveErrorCodeData(errorCodeData model.ErrorInfo) error {
	filePath := "errorCodes.json"
	log.Printf("[AddErrorCodeHandler] Saving error code data to file: %s", filePath)

	// Read existing data
	existingData, err := readExistingData(filePath)
	if err != nil {
		return err
	}

	// Add new error code using Z4Z5_ErrorCode as key
	key := strconv.Itoa(errorCodeData.Z4Z5_ErrorCode)
	log.Printf("[AddErrorCodeHandler] Adding error code with key: '%s'", key)

	// Check if key already exists
	if existingInfo, exists := existingData[key]; exists {
		log.Printf("[AddErrorCodeHandler] Error: Error code key '%s' already exists", key)
		log.Printf("[AddErrorCodeHandler] Existing error code details:")
		log.Printf("[AddErrorCodeHandler]   - Z0 (Generator Status): %d", existingInfo.Z0)
		log.Printf("[AddErrorCodeHandler]   - Z1 (Simplified Error Code): %d", existingInfo.Z1)
		log.Printf("[AddErrorCodeHandler]   - Z2 (Display Bitmap): %d", existingInfo.Z2)
		log.Printf("[AddErrorCodeHandler]   - Z3_Phase: %d", existingInfo.Z3_Phase)
		log.Printf("[AddErrorCodeHandler]   - Z3_Class: %d", existingInfo.Z3_Class)
		log.Printf("[AddErrorCodeHandler]   - Z4Z5_ErrorCode: %d", existingInfo.Z4Z5_ErrorCode)
		log.Printf("[AddErrorCodeHandler]   - Z6Z7_ErrorData: %d", existingInfo.Z6Z7_ErrorData)
		log.Printf("[AddErrorCodeHandler]   - Description: '%s'", existingInfo.Description)
		return &httpError{message: fmt.Sprintf("Error code %s already exists in the database", key)}
	}

	existingData[key] = errorCodeData
	log.Printf("[AddErrorCodeHandler] Error code data added to memory, total entries: %d", len(existingData))

	// Write updated data back to file
	return writeDataToFile(filePath, existingData)
}

// readExistingData reads and parses existing JSON data from file
func readExistingData(filePath string) (map[string]model.ErrorInfo, error) {
	var existingData map[string]model.ErrorInfo

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		// If file doesn't exist, create a new empty map
		log.Printf("[AddErrorCodeHandler] File '%s' not found, creating new error codes database", filePath)
		return make(map[string]model.ErrorInfo), nil
	}

	log.Printf("[AddErrorCodeHandler] Reading existing data from file: %s (%d bytes)", filePath, len(fileData))

	// Parse existing JSON data
	err = json.Unmarshal(fileData, &existingData)
	if err != nil {
		log.Printf("[AddErrorCodeHandler] Failed to parse JSON from file '%s': %v", filePath, err)
		return nil, &httpError{message: "Error processing existing data"}
	}

	log.Printf("[AddErrorCodeHandler] Successfully loaded %d existing error codes", len(existingData))
	return existingData, nil
}

// writeDataToFile marshals data to JSON and writes to file
func writeDataToFile(filePath string, data map[string]model.ErrorInfo) error {
	log.Printf("[AddErrorCodeHandler] Writing %d error codes to file: %s", len(data), filePath)

	// Marshal data to JSON with indentation
	updatedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("[AddErrorCodeHandler] Failed to marshal data to JSON: %v", err)
		return &httpError{message: "Error processing data"}
	}

	log.Printf("[AddErrorCodeHandler] JSON data prepared (%d bytes)", len(updatedData))

	// Write data to file
	err = ioutil.WriteFile(filePath, updatedData, os.ModePerm)
	if err != nil {
		log.Printf("[AddErrorCodeHandler] Failed to write file '%s': %v", filePath, err)
		return &httpError{message: "Error saving data"}
	}

	log.Printf("[AddErrorCodeHandler] Successfully wrote %d bytes to file: %s", len(updatedData), filePath)
	return nil
}

// respondWithAddSuccess sends a success response for add error code operation
func respondWithAddSuccess(w http.ResponseWriter, message string) {
	log.Printf("[AddErrorCodeHandler] Sending success response: %s", message)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": message,
	})
}
