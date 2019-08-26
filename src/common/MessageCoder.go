package common

import (
	"strconv"
)

type MessageCoder struct {
	decoders map[string]MessageDecoder
	encoders map[string]MessageEncoder
}

func (s *MessageCoder) InitMessageCoder() {
	s.decoders = make(map[string]MessageDecoder)
	s.encoders = make(map[string]MessageEncoder)
}
func (s *MessageCoder) RegisterDecoder(messageId, version int, decoder MessageDecoder) {
	key := strconv.Itoa(messageId) + ":" + strconv.Itoa(version)
	s.decoders[key] = decoder
}
func (s *MessageCoder) RegisterEncoder(messageId, version int, encoder MessageEncoder) {
	key := strconv.Itoa(messageId) + ":" + strconv.Itoa(version)
	s.encoders[key] = encoder
}

func (s *MessageCoder) GetDecoder(messageId, version int) MessageDecoder {
	key := strconv.Itoa(messageId) + ":" + strconv.Itoa(version)
	return s.decoders[key]
}

func (s *MessageCoder) GetEncoder(messageId, version int) MessageEncoder {
	key := strconv.Itoa(messageId) + ":" + strconv.Itoa(version)
	return s.encoders[key]
}
