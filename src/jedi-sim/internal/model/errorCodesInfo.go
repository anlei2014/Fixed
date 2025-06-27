package model

// ErrorInfo defines all error-related fields in a CAN message
type ErrorInfo struct {
	Z0             int    // Z0: Generator Status (e.g., 0x06 = Error condition)
	Z1             int    // Z1: Simplified Error Code (e.g., 30 = Simplified error code 30)
	Z2             byte   // Z2: Display Bitmap (bit 0-7)
	Z3_Phase       int    // Z3: Phase (high 4 bits)
	Z3_Class       int    // Z3: Error Class (low 4 bits)
	Z4Z5_ErrorCode int    // Z4/Z5: Generator Error Code (little-endian format)
	Z6Z7_ErrorData int    // Z6/Z7: Error Data (little-endian format)
	Description    string // Description of the error
}

// ErrorMap maps error codes to their corresponding ErrorInfo
var ErrorMap = map[int]ErrorInfo{
	804: {
		Z0:             6,      // Z0: Generator Status (Error condition)
		Z1:             30,     // Z1: Simplified error code 30
		Z2:             0x00,   // Z2: Display Bitmap (All status OK)
		Z3_Phase:       2,      // Z3: Phase=2 (Preparation in Progress)
		Z3_Class:       2,      // Z3: Error Class=2 (DEBUG)
		Z4Z5_ErrorCode: 804,    // Z4/Z5: Generator Error Code (0x0324)
		Z6Z7_ErrorData: 0x0005, // Z6/Z7: Error Data (e.g., temperature 5°C)
		Description:    "Tube spit (all kV drop/regul errors)",
	},
	134: {
		Z0:             6,      // Z0: Generator Status (Error condition)
		Z1:             90,     // Z1: Simplified error code 30
		Z2:             0x00,   // Z2: Display Bitmap (All status OK)
		Z3_Phase:       2,      // Z3: Phase=2 (Preparation in Progress)
		Z3_Class:       5,      // Z3: Error Class=2 (DEBUG)
		Z4Z5_ErrorCode: 134,    // Z4/Z5: Generator Error Code (0x86)
		Z6Z7_ErrorData: 0x0005, // Z6/Z7: Error Data (e.g., temperature 5°C)
		Description:    "This error is usually caused by a HW issue.",
	},
	// Add other error codes here...
}
