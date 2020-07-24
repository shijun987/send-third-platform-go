package utils

import (
	"strconv"
)

// String2float 字符串转浮点数
func String2float(value string) float64 {
	ret, _ := strconv.ParseFloat(value, 32)
	return ret
}

// Crc16 Crc16
func Crc16(buf []byte, len int) uint16 {
	var crc uint16 = 0xFFFF
	var polynomial uint16 = 0xA001

	if len == 0 {
		return 0
	}

	for i := 0; i < len; i++ {
		crc ^= uint16(buf[i]) & 0x00FF
		for j := 0; j < 8; j++ {
			if (crc & 0x0001) != 0 {
				crc >>= 1
				crc ^= polynomial
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}
