package jediSim

import (
	"fmt"
	"jedi-sim/internal/model"
)

// const (
// 	MSG_DATA_LENGTH = 8
// )

var (
	jediMode        int
	jediOption      = "MULTISEQ"
	jediOptionValue = make(map[string]float64)
	jediError       = "NO_ERROR"
	gkv             float64
	gmA             float64
	gTime           float64
	gExpNum         int = 70
	frameRate       int
	saveMsgForDE    [MSG_LENGTH]int
)

func init() {
	//setJediOption("MULTISEQ")
}

func sendReactMSG(msg [MSG_LENGTH]int) {
	//GetWriteCanMSGHandler().SendMessage(msg[:])
}

func GeneratorStatusErrorCode(msg [MSG_LENGTH]int, errorCode int) (bool, [MSG_LENGTH]int) {
	myMsg := CreateMsg("NOTIFY_JEDI_STATUS", msg)

	// Validate message validity
	if myMsg[0] == 0 && myMsg[1] == 0 {
		fmt.Printf("[ERROR] Failed to create message: Invalid ID or length\n")
		return false, [MSG_LENGTH]int{}
	}

	// Look up error code information
	errInfo, ok := model.ErrorMap[errorCode]
	if !ok {
		fmt.Printf("[ERROR] Unknown error code: %d\n", errorCode)
		return false, [MSG_LENGTH]int{}
	}

	// Set Z0: Generator Status
	myMsg = SetParam(myMsg, "JEDI_STATUS_PHASE", errInfo.Z0)

	// Set Z1: Simplified Error Code
	myMsg = SetParam(myMsg, "SIMPLIFIED_ERROR_CODE", errInfo.Z1)

	// Set Z2: Display Bitmap (bit 0-7)
	myMsg = SetParam(myMsg, "DISPLAY_BITMAP", int(errInfo.Z2))

	// Set Z3: Phase (high 4 bits) + Error Class (low 4 bits)
	myMsg = SetParam(myMsg, "PHASE_OF_ERROR_OCCURED", errInfo.Z3_Phase)
	myMsg = SetParam(myMsg, "ERROR_CLASS", errInfo.Z3_Class)

	// Set Z4/Z5: Generator Error Code (little-endian format)
	myMsg = SetParam(myMsg, "JEDI_ERROR_CODE", errInfo.Z4Z5_ErrorCode)

	// Set Z6/Z7: Error Data (little-endian format)
	myMsg = SetParam(myMsg, "DATA_RELATED_TO_ERROR", errInfo.Z6Z7_ErrorData)

	// Log generated CAN message fields
	fmt.Printf("[DEBUG] Generated CAN Message: Z0=0x%02X, Z1=0x%02X, Z2=0x%02X, Z3=0x%02X, ErrorCode=0x%04X, ErrorData=0x%04X\n",
		myMsg[2], myMsg[3], myMsg[4], myMsg[5], errInfo.Z4Z5_ErrorCode, errInfo.Z6Z7_ErrorData)

	// Send the message via CAN bus
	sendReactMSG(myMsg)

	return true, myMsg
}

// func GeneratorStatusErrorCode(msg [MSG_LENGTH]int, errorCode int) (bool, [MSG_LENGTH]int) {
// 	// Create CAN message based on "NOTIFY_JEDI_STATUS" configuration
// 	myMsg := CreateMsg("NOTIFY_JEDI_STATUS", msg)

// 	// Basic validation: Check if message ID and length are valid
// 	if myMsg[0] == 0 && myMsg[1] == 0 {
// 		fmt.Printf("[ERROR] Failed to create message: Invalid ID or length\n")
// 		return false, [MSG_LENGTH]int{} // Return failure and empty message
// 	}

// 	// Set generator status (written into Z0 field)
// 	myMsg = SetParam(myMsg, "JEDI_STATUS_PHASE", 6)
// 	// Set generator status (written into Z0 field)
// 	myMsg = SetParam(myMsg, "SIMPLIFIED_ERROR_CODE", 40)
// 	// Set error code (written into Z4/Z5 fields)
// 	myMsg = SetParam(myMsg, "JEDI_ERROR_CODE", errorCode)

// 	// Log critical fields for debugging:
// 	// - Z0 (Generator Status)
// 	// - Z4 (Error Code LSB)
// 	// - Z5 (Error Code MSB)
// 	fmt.Printf("[DEBUG] Generated CAN Message: msgID=0x%02X, Z0=0x%02X, Z4=0x%02X, Z5=0x%02X\n",
// 		myMsg[0], myMsg[2], myMsg[6], myMsg[7])

// 	// Send the message via CAN bus
// 	sendReactMSG(myMsg)

// 	// Return success status and generated message
// 	return true, myMsg
// }
