package utils

func Uint16ToUint32(a, b uint16) uint32 {
	return uint32(a)<<16 | uint32(b)
}

func Uint32ToUint16(a uint32) (uint16, uint16) {
	return uint16(a >> 16), uint16(a)
}
