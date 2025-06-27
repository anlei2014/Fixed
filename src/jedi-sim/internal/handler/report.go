package handler

import (
	"encoding/json"
	"jedi-sim/internal/service"
	"jedi-sim/jediSim" // Import GeneratorStatusErrorCode
	"log"
	"net/http"
	"strconv"
)

// ReportHandler handles error code reporting and CAN message generation
func ReportHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ReportHandler] Request received:", r.Method, r.URL.Path)

	// Only accept POST requests
	if r.Method != "POST" {
		log.Println("[ReportHandler] Invalid request method:", r.Method)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get error code from form
	errorCode := r.FormValue("errorcode")
	log.Printf("[ReportHandler] Error code received: %s", errorCode)

	// Validate error code format
	valid, message := service.ValidateErrorCode(errorCode)
	log.Printf("[ReportHandler] Validation result: valid=%v, message=%s", valid, message)

	if !valid {
		log.Println("[ReportHandler] Error code validation failed")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": message,
		})
		return
	}

	// Convert error code string to int
	codeInt, err := strconv.Atoi(errorCode)
	if err != nil {
		log.Printf("[ReportHandler] Invalid error code format: %s", errorCode)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid error code format",
		})
		return
	}

	// Generate CAN message using GeneratorStatusErrorCode
	// This function already handles message creation AND sending
	msg := [10]int{} // MSG_LENGTH = 10
	success, myMsg := jediSim.GeneratorStatusErrorCode(msg, codeInt)
	if !success {
		log.Println("[ReportHandler] Failed to generate/send CAN message")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to generate/send CAN message",
		})
		return
	}

	// Log only critical fields:
	// - Z0 (Generator Status)
	// - Z4/Z5 (Error Code)
	log.Printf("[ReportHandler] Generated CAN Message: Z0=0x%02X, Z4=0x%02X, Z5=0x%02X",
		myMsg[0], myMsg[6], myMsg[7])

	// Prepare minimal CAN message data for response (first 8 bytes)
	canMsgData := []int{
		myMsg[0], myMsg[1], myMsg[2], myMsg[3],
		myMsg[4], myMsg[5], myMsg[6], myMsg[7],
		myMsg[8], myMsg[9],
	}

	// Return success response
	log.Println("[ReportHandler] CAN message processed successfully")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"code":    errorCode,
		"message": "OK, CAN message has been sent",
		"canMsg":  canMsgData,
	})
}
