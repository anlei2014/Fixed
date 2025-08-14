package handler

import (
	"encoding/json"
	"jedi-sim/internal/service"
	"jedi-sim/jediSim"
	"log"
	"net/http"
	"strconv"
)

// ReportHandler handles error code reporting and CAN message generation
func ReportHandler(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method
	if err := validateHTTPMethod(r, "POST"); err != nil {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Extract and validate error code
	errorCode, err := extractAndValidateErrorCode(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generate CAN message
	canMsg, err := generateCANMessage(errorCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return success response
	respondWithSuccess(w, errorCode, canMsg)
}

// validateHTTPMethod checks if the request method matches the expected method
func validateHTTPMethod(r *http.Request, expectedMethod string) error {
	if r.Method != expectedMethod {
		return &httpError{message: "Method not allowed"}
	}
	return nil
}

// extractAndValidateErrorCode extracts error code from request and validates it
func extractAndValidateErrorCode(r *http.Request) (string, error) {
	errorCode := r.FormValue("errorcode")

	// Validate error code format and existence
	valid, message := service.ValidateErrorCode(errorCode)
	if !valid {
		return "", &httpError{message: message}
	}

	return errorCode, nil
}

// generateCANMessage converts error code to integer and generates CAN message
func generateCANMessage(errorCode string) ([]int, error) {
	// Convert error code string to integer
	codeInt, err := strconv.Atoi(errorCode)
	if err != nil {
		return nil, &httpError{message: "Invalid error code format"}
	}

	// Generate CAN message using GeneratorStatusErrorCode
	msg := [10]int{} // MSG_LENGTH = 10
	success, myMsg := jediSim.GeneratorStatusErrorCode(msg, codeInt)
	if !success {
		return nil, &httpError{message: "Failed to generate/send CAN message"}
	}

	// Prepare CAN message data for response (all 10 bytes)
	canMsgData := []int{
		myMsg[0], myMsg[1], myMsg[2], myMsg[3],
		myMsg[4], myMsg[5], myMsg[6], myMsg[7],
		myMsg[8], myMsg[9],
	}

	return canMsgData, nil
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	log.Printf("[ReportHandler] Error: HTTP %d - %s", statusCode, message)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"message": message,
	})
}

// respondWithSuccess sends a success response with CAN message data
func respondWithSuccess(w http.ResponseWriter, errorCode string, canMsg []int) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"code":    errorCode,
		"message": "OK, CAN message has been sent",
		"canMsg":  canMsg,
	})
}

// httpError represents an HTTP error
type httpError struct {
	message string
}

func (e *httpError) Error() string {
	return e.message
}
