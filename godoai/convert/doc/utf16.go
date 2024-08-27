package doc

import (
	"encoding/binary"
)

type utf16Buffer struct {
	haveReadLowerByte bool
	char              [2]byte
	data              []uint16
}

func (buf *utf16Buffer) Write(p []byte) (n int, err error) {
	for i := range p {
		buf.WriteByte(p[i])
	}
	return len(p), nil
}

func (buf *utf16Buffer) WriteByte(b byte) error {
	if buf.haveReadLowerByte {
		buf.char[1] = b
		buf.data = append(buf.data, binary.LittleEndian.Uint16(buf.char[:]))
	} else {
		buf.char[0] = b
	}
	buf.haveReadLowerByte = !buf.haveReadLowerByte
	return nil
}

func (buf *utf16Buffer) Chars() []uint16 {
	if buf.haveReadLowerByte {
		return append(buf.data, uint16(buf.char[0]))
	}
	return buf.data
}
