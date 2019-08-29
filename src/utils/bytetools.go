package utils

import (
	"bytes"
	"encoding/binary"
)

type ByteTools struct {
}

func (bt *ByteTools) BytesToString(b []byte, len int) string {
	ret := string(b[0:len])
	return ret
}

func (bt *ByteTools) StringToBytes(n string) []byte {
	return []byte(n)
}

func (bt *ByteTools) WriteString(buffer *bytes.Buffer, v string) {
	l := len(v)
	(*buffer).Write(bt.ShortToBytes(l))
	if l > 0 {
		(*buffer).Write([]byte(v))
	}
}

func (bt *ByteTools) ReadString(data []byte, pos *int) string {
	len := bt.BytesToShort(data, pos)
	ret := string(data[*pos : *pos+len])
	*pos += len
	return ret
}

func (bt *ByteTools) BytesToInt(b []byte, pos *int) int {
	bytesBuffer := bytes.NewBuffer(b[*pos : *pos+4])
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	*pos += 4
	return int(x)
}

func (bt *ByteTools) IntToBytes(n int32) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func (bt *ByteTools) BytesToShort(b []byte, pos *int) int {
	bytesBuffer := bytes.NewBuffer(b[*pos : *pos+2])
	var x int16
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	*pos += 2
	return int(x)
}

func (bt *ByteTools) ShortToBytes(n int) []byte {
	x := int16(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func (bt *ByteTools) BytesToBool(b []byte, pos *int) bool {
	bytesBuffer := bytes.NewBuffer(b[*pos : *pos+1])
	var x bool
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	*pos += 1
	return x
}

func (bt *ByteTools) BoolToBytes(n bool) []byte {
	x := n
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func (bt *ByteTools) BytesToFloat(b []byte, pos *int) float32 {
	bytesBuffer := bytes.NewBuffer(b[*pos : *pos+4])
	var x float32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	*pos += 4
	return x
}

func (bt *ByteTools) FloatToBytes(n float32) []byte {
	x := float32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func (bt *ByteTools) BytesToLong(b []byte, pos *int) int64 {
	bytesBuffer := bytes.NewBuffer(b[*pos : *pos+8])
	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	*pos += 4
	return x
}

func (bt *ByteTools) LongToBytes(n int64) []byte {
	x := int64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
