package common

type MessageEncoder interface {
	Encode(messageId, version int, msgType int, msg map[string]interface{}) []byte
}
