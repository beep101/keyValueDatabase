package utils

import (
	"encoding/binary"
)

func IntToByte(x int) []byte {
	bts := make([]byte, 4)
	binary.LittleEndian.PutUint32(bts, uint32(x))
	return bts
}

func ByteToInt(bts []byte) int {
	x := int(binary.LittleEndian.Uint32(bts))
	return x
}

func Int64ToByte(x int64) []byte {
	bts := make([]byte, 8)
	binary.LittleEndian.PutUint64(bts, uint64(x))
	return bts
}

func ByteToInt64(bts []byte) int64 {
	x := int64(binary.LittleEndian.Uint64(bts))
	return x
}
