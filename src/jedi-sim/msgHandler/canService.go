package msgHandler

import "log"

// SendCANMessage organizes and sends a CAN message, returns success and error info
func SendCANMessage(errorCode string) (bool, string) {
    log.Printf("[SendCANMessage] Organizing CAN message for error code: %s", errorCode)
    // Here you can organize CAN message content based on errorCode
    // In a real project, call the underlying CAN bus send interface
    // Example logic: errorCode "110" is considered success, others fail
    if errorCode == "110" {
        log.Println("[SendCANMessage] CAN message send simulated as success")
        return true, ""
    }
    log.Println("[SendCANMessage] CAN message send simulated as failure")
    return false, "CAN bus error"
}