package handler

import (
	"encoding/json"
	"jedi-sim/internal/service"
	"jedi-sim/msgHandler"
	"log"
	"net/http"
	"strconv"
)

// CAN208Message represents the structure of CAN message with ID 208
type CAN208Message struct {
	Z0 byte // Generator Status (Phase)
	Z1 byte // Simplified error or warning code
	Z2 byte // Display bitmap
	Z3 byte // Phase/Error class
	Z4 byte // Generator error/warning code (high byte)
	Z5 byte // Generator error/warning code (low byte)
	Z6 byte // Data related to error/warning
	Z7 byte // Data related to error/warning
}

// CreateCANMessage creates a CAN208Message with errorCode filled in Z4/Z5
func CreateCANMessage(errorCode string) (CAN208Message, error) {
	codeInt, err := strconv.Atoi(errorCode)
	if err != nil {
		log.Printf("[CreateCANMessage] Invalid errorCode: %s", errorCode)
		return CAN208Message{}, err
	}
	z4 := byte((codeInt >> 8) & 0xFF)
	z5 := byte(codeInt & 0xFF)
	msg := CAN208Message{
		Z0: 0x00, // Example: Idle
		Z1: 0x00, // No error/warning
		Z2: 0x00, // All OK
		Z3: 0x00, // Example: Phase 0, Error class 0
		Z4: z4,
		Z5: z5,
		Z6: 0x00,
		Z7: 0x00,
	}
	log.Printf("[CreateCANMessage] Created CAN208Message: %+v", msg)
	return msg, nil
}

func ReportHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ReportHandler] Received request:", r.Method, r.URL.Path)

	if r.Method != "POST" {
		log.Println("[ReportHandler] Invalid request method:", r.Method)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	errorCode := r.FormValue("errorcode")
	log.Printf("[ReportHandler] Received error code: %s", errorCode)
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

	// Create CAN message and print
	canMsg, err := CreateCANMessage(errorCode)
	if err != nil {
		log.Println("[ReportHandler] Failed to create CAN message:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to create CAN message: " + err.Error(),
		})
		return
	}
	log.Printf("[ReportHandler] CAN208Message to send: %+v", canMsg)

	// Call and print CAN message sending result
	sendOK, sendMsg := msgHandler.SendCANMessage(errorCode)
	log.Printf("[ReportHandler] SendCANMessage result: ok=%v, msg=%s", sendOK, sendMsg)

	if sendOK {
		log.Println("[ReportHandler] CAN message sent successfully")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  true,
			"code":    errorCode,
			"message": "OK, CAN message has been sent",
			"canMsg":  canMsg,
		})
	} else {
		log.Println("[ReportHandler] CAN message send failed:", sendMsg)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "CAN message send failed: " + sendMsg,
			"canMsg":  canMsg,
		})
	}
}