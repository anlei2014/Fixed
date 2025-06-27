package jediSim

import (
	"math"
)

var setParamFuncMap = make(map[string]func([MSG_LENGTH]int, string, int) [MSG_LENGTH]int)

func init() {

	setParamFuncMap["NOTIFY_JEDI_STATUS"] = func(msg [MSG_LENGTH]int, param string, value int) [MSG_LENGTH]int {
		var myMsg [MSG_LENGTH]int
		switch param {
		case "JEDI_STATUS_PHASE":
			myMsg = write1Byte(msg, 0, value)
		case "SIMPLIFIED_ERROR_CODE":
			myMsg = write1Byte(msg, 1, value)
		case "DISPLAY_BITMAP":
			myMsg = write1Byte(msg, 2, value)
		case "PHASE_OF_ERROR_OCCURED":
			myMsg = writeNBits(msg, 3, 0, 4, value)
		case "ERROR_CLASS":
			myMsg = writeNBits(msg, 3, 4, 4, value)
		case "JEDI_ERROR_CODE":
			myMsg = write2Byte(msg, 4, value)
		case "DATA_RELATED_TO_ERROR":
			myMsg = write2Byte(msg, 6, value)
		default:
			myMsg = msg
		}
		return myMsg
	}

}

func writeFloat(msg [MSG_LENGTH]int, pos int, val int) [MSG_LENGTH]int {
	myVal := floatToULong(float32(val))
	msg0 := write4Byte(msg, pos, int(myVal))
	return msg0
}

func floatToULong(val float32) uint32 {
	res := math.Float32bits(val)
	return res
}

func writeNothing(msg [MSG_LENGTH]int, param string, value int) [MSG_LENGTH]int {
	return msg
}

func write1Byte(msg [MSG_LENGTH]int, pos int, val int) [MSG_LENGTH]int {
	myPos := pos + 2
	myMsg := msg
	myMsg[myPos] = val
	return myMsg
}

func write2Byte(msg [MSG_LENGTH]int, pos int, val int) [MSG_LENGTH]int {
	val0 := val % 0x100
	msg0 := write1Byte(msg, pos, val0)
	val0 = val / 0x100
	pos0 := pos + 1
	msg0 = write1Byte(msg0, pos0, val0)
	return msg0
}

func write3Byte(msg [MSG_LENGTH]int, pos int, val int) [MSG_LENGTH]int {
	val0 := val % 0x100
	msg0 := write1Byte(msg, pos, val0)

	val0 = val / 0x100
	pos0 := pos + 1
	msg0 = write2Byte(msg0, pos0, val0)
	return msg0
}

func write4Byte(msg [MSG_LENGTH]int, pos int, val int) [MSG_LENGTH]int {
	val0 := val % 0x100
	msg0 := write1Byte(msg, pos, val0)

	val0 = val / 0x100
	pos0 := pos + 1
	msg0 = write3Byte(msg0, pos0, val0)

	return msg0
}

func writeNBits(msg [MSG_LENGTH]int, bytePos int, bitPos int, len int, val int) [MSG_LENGTH]int {
	myBytePos := bytePos + 2
	myVal := msg[myBytePos]

	myMaskVal := val << uint(bitPos)

	mask := 0xFF
	mask >>= uint(8 - len)
	mask <<= uint(bitPos)
	mask = ^mask
	mask = 0xFF & mask

	myVal &= mask
	myVal |= myMaskVal
	myMsg := write1Byte(msg, bytePos, myVal)
	return myMsg

}
