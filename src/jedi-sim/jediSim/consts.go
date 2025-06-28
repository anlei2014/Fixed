package jediSim

var (
	MSGs = []string{
		"GET_JEDI_STATUS",
		"NOTIFY_JEDI_STATUS",
	}
	msgId = map[string]int{
		"GET_JEDI_STATUS":    520,
		"NOTIFY_JEDI_STATUS": 520,
	}
	msgLn = map[string]int{
		"GET_JEDI_STATUS":    16, // Length of the GET_JEDI_STATUS
		"NOTIFY_JEDI_STATUS": 8,  // Length of the NOTIFY_JEDI_STATUS
	}
)

// MSG_LENGTH is the length of the CAN message
