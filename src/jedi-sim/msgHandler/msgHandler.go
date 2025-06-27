package msgHandler

type MSGHandler interface {
	// function to send data to canbus or other bus
	SendMessage(msg []int)
	// function to read data from  canbus or other bus
	ReadMessage() []int
	// function to call JediSim/Jedi actions when receiving messages
	UponMessage(msg []int)
	// function to be called from JediSim/Jedi when react messages ready
	UponReactedMessage(msg []int)
}
