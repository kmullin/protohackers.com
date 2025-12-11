package main

import (
	"errors"
	"io"

	"github.com/rs/zerolog"
)

var errNegativeOffset = errors.New("negative offset")

type sender struct {
	stream
}

func (s *sender) CheckLines() {
	// return iterator ?
}

// stream handles a stream of data
type stream struct {
	pos int
	buf []byte
}

func (s *stream) Pos() int {
	return s.pos
}

func (s *stream) Buf() []byte {
	return s.buf
}

func (s *stream) WriteAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errNegativeOffset
	}

	endPos := int(off) + len(p)

	// Grow buffer if needed
	if endPos > len(s.buf) {
		newBuf := make([]byte, endPos)
		copy(newBuf, s.buf)
		s.buf = newBuf
	}

	// Overwrite
	copy(s.buf[off:endPos], p)
	s.pos += len(p)
	return len(p), nil
}

func (s *stream) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errNegativeOffset
	}

	if int(off) >= len(s.buf) {
		return 0, io.EOF
	}

	n = copy(p, s.buf[off:])
	if n < len(p) {
		return n, io.EOF
	}

	return
}

func (s *stream) MarshalZerologObject(e *zerolog.Event) {
	e.Bytes("buf", s.buf).Int("pos", s.pos).Int("len", len(s.buf))
}
