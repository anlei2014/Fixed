package handler

import (
	"encoding/json"
	"jedi-sim/internal/service"
	"jedi-sim/msgHandler"
	"log"
	"net/http"
)

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

	// Call and print CAN message sending result
	sendOK, sendMsg := msgHandler.SendCANMessage(errorCode)
	log.Printf("[ReportHandler] SendCANMessage result: ok=%v, msg=%s", sendOK, sendMsg)

	if sendOK {
		log.Println("[ReportHandler] CAN message sent successfully")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  true,
			"code":    errorCode,
			"message": "OK, CAN message has been sent",
		})
	} else {
		log.Println("[ReportHandler] CAN message send failed:", sendMsg)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "CAN message send failed: " + sendMsg,
		})
	}
}