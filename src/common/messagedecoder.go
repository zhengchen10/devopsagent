package common

type MessageDecoder interface {
	Decode(messageId int, version int, msgType int, data []byte) map[string]interface{}
}
