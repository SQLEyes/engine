package util

import (
	"bytes"
	"encoding/binary"
)

type ByteBuffer struct {
	len  int64
	pos  int64
	data []byte
}

func NewByteBuffer(data []byte) *ByteBuffer {
	return &ByteBuffer{data: data, pos: 0, len: int64(len(data))}
}
func (s *ByteBuffer) HasNext() bool {
	return s.pos < s.len
}

func (s *ByteBuffer) ReadShort() byte {
	s.Check(1)
	bs := s.data[s.pos : s.pos+1]
	s.pos++
	return bs[0]
}
func (s *ByteBuffer) ReadInt32() []byte {
	s.Check(4)
	bs := s.data[s.pos : s.pos+4]
	s.pos += 4
	return bs
}
func (s *ByteBuffer) ReadInt16() []byte {
	s.Check(2)
	bs := s.data[s.pos : s.pos+2]
	s.pos += 2
	return bs
}

func (s *ByteBuffer) Read(len int64) []byte {
	s.Check(len)
	bs := s.data[s.pos : s.pos+len]
	s.pos += len
	return bs
}
func (s *ByteBuffer) GetString(len int64) string {
	return string(bytes.TrimRight(s.Read(len), "\x00"))
}
func (s *ByteBuffer) GetInt32() int32 {
	ui := binary.BigEndian.Uint32(s.ReadInt32())
	return int32(ui)
}
func (s *ByteBuffer) GetInt16() int16 {
	ui := binary.BigEndian.Uint16(s.ReadInt16())
	return int16(ui)
}

func (s *ByteBuffer) Check(len int64) {
	if (s.pos + len) > s.len {
		panic("index of array")
	}
}
