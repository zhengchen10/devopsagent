package utils

import (
	"bytes"
	"encoding/binary"
)

type ByteTools struct {
}

func (bt *ByteTools) BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

func (bt *ByteTools) IntToBytes(n int32) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func (bt *ByteTools) BytesToShort(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int16
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

func (bt *ByteTools) ShortToBytes(n int) []byte {
	x := int16(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func (bt *ByteTools) BytesToBool(b []byte) bool {
	bytesBuffer := bytes.NewBuffer(b)
	var x bool
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

func (bt *ByteTools) BoolToBytes(n bool) []byte {
	x := n
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func (bt *ByteTools) BytesToFloat(b []byte) float32 {
	bytesBuffer := bytes.NewBuffer(b)
	var x float32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

func (bt *ByteTools) FloatToBytes(n float32) []byte {
	x := float32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func (bt *ByteTools) BytesToLong(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

func (bt *ByteTools) LongToBytes(n int64) []byte {
	x := int64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
