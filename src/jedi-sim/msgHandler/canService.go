package msgHandler

import (
    "fmt"
    "log"
)

// CAN208Message represents the structure of CAN message with ID 208
type CAN208Message struct {
    Z0 byte
    Z1 byte
    Z2 byte
    Z3 byte
    Z4 byte
    Z5 byte
    Z6 byte
    Z7 byte
}

// SendCANMessage organizes and sends a CAN message, returns success and error info
func SendCANMessage(errorCode string) (bool, string) {
    log.Printf("[SendCANMessage] Organizing CAN message for error code: %s", errorCode)

    // Example: create a CAN208Message with errorCode in Z4/Z5
    var msg CAN208Message
    if errorCode != "" {
        // You can add your own logic to fill other fields as needed
        // Here just fill Z4/Z5 with errorCode if it is a valid int
        // (If errorCode is not a valid int, Z4/Z5 will be 0)
        var codeInt int
        _, err := fmt.Sscanf(errorCode, "%d", &codeInt)
        if err == nil {
            msg.Z4 = byte((codeInt >> 8) & 0xFF)
            msg.Z5 = byte(codeInt & 0xFF)
        } else {
            log.Printf("[SendCANMessage] Error parsing errorCode to int: %v", err)
        }
    }

    log.Printf("[SendCANMessage] CAN208Message to send: %+v", msg)

    // Simulate CAN message send: treat any non-empty errorCode as success
    if errorCode != "" {
        log.Println("[SendCANMessage] CAN message send simulated as success")
        return true, ""
    }
    log.Println("[SendCANMessage] CAN message send simulated as failure")
    return false, "CAN bus error"
}