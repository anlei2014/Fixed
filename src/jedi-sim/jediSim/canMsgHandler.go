package jediSim

// "jedi-sim/can"

const (
	MSG_DATA_LENGTH = 8
)

// Declare the map at the package level
var zmpClientCmdMap = map[string]int{
	"TubeSpitError_1":  1,
	"TubeSpitError_5":  5,
	"TubeSpitError_15": 15,
}

func SetParam(msg [MSG_LENGTH]int, param string, value int) [MSG_LENGTH]int {
	var myMsg [MSG_LENGTH]int
	id := msg[0]
	ln := msg[1]
	//fmt.Printf("[SetParam] Input ID: 0x%02X, Length: %d\n", id, ln)
	for _, key := range MSGs {
		if id == msgId[key] && ln == msgLn[key] {
			myMsg = setParamFuncMap[key](msg, param, value)
			break
		}
	}
	return myMsg
}

func CreateMsg(name string, msg [MSG_LENGTH]int) [MSG_LENGTH]int {
	var myMsg [MSG_LENGTH]int
	myMsg[0] = msgId[name]
	myMsg[1] = msgLn[name]
	for i := 2; i < MSG_LENGTH; i++ {
		myMsg[i] = msg[i]
	}
	return myMsg
}
