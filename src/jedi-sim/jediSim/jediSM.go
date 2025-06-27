package jediSim

const (
	CAN_BUS_ERROR = "CAN_BUS_ERROR"
	MSG_LENGTH    = 10
)

var (
	subState = map[string][]string{
		"S_POWERED_UP": {
			"S_SYS_ID_RECEIVED",
			"S_REPLYING_JEDI_CAPABILITIES",
			"S_REPLYING_STATIC_ID",
			"S_SET_JEDI_TIME",
		},
	}
	_curState = "S_POWERED_UP"
	err       string
)
