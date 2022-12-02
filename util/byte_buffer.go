package util

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type ByteBuffer struct {
	len       int
	pos       int
	data      []byte
	bigending bool
}

func NewByteBuffer(data []byte, ending ...bool) *ByteBuffer {
	flag := true
	if len(ending) > 0 {
		flag = ending[0]
	}
	return &ByteBuffer{data: data, pos: 0, len: len(data), bigending: flag}
}
func (s *ByteBuffer) HasNext() bool {
	return s.pos < s.len
}
func (s *ByteBuffer) Print() {
	fmt.Println(hex.EncodeToString(s.data))
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
func (s *ByteBuffer) ReadEnd() []byte {
	bs := s.data[s.pos:]
	s.pos = s.len
	return bs
}
func (s *ByteBuffer) Read(len int) []byte {
	s.Check(len)
	bs := s.data[s.pos : s.pos+len]
	s.pos += len
	return bs
}
func (s *ByteBuffer) GetString(len int) string {
	return fmt.Sprintf("%s", s.Read(len))
}
func (s *ByteBuffer) GetInt() (i int) {
	if s.bigending {
		ui := binary.BigEndian.Uint32(s.ReadInt32())
		i = int(ui)
	} else {
		ui := binary.LittleEndian.Uint32(s.ReadInt32())
		i = int(ui)
	}
	return
}

func (s *ByteBuffer) GetInt16() (i int16) {
	if s.bigending {
		ui := binary.BigEndian.Uint16(s.ReadInt16())
		i = int16(ui)
	} else {
		ui := binary.LittleEndian.Uint16(s.ReadInt16())
		i = int16(ui)
	}
	return
}

func (s *ByteBuffer) Check(len int) {
	if (s.pos + len) > s.len {
		panic(fmt.Sprintf("ArrayIndexOutOfBoundsException： %d + %d = %d > %d ", s.pos, len, s.pos+len, s.len))
	}
}
func (s *ByteBuffer) Position(pos ...int) int {
	if len(pos) > 0 {
		s.pos = s.pos + pos[0]
	}
	return s.pos
}
func (s *ByteBuffer) Len() int {
	return s.len
}
