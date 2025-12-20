package main

import (
	"bytes"
	"errors"
	"io"
	"sync"

	"github.com/rs/zerolog"
)

var errNegativeOffset = errors.New("negative offset")

// stream handles a stream of data
type stream struct {
	readPos int // where in our stream we've read from
	buf     []byte
	mu      sync.RWMutex
}

func (s *stream) UnreadLen() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.buf) - s.readPos
}

func (s *stream) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.buf)
}

// UnreadBuf returns a copy of the buffer
func (s *stream) UnreadBuf(b []byte) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return copy(b, s.buf[s.readPos:])
}

func (s *stream) Readline() (line []byte, pos int, fullLine bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.readPos >= len(s.buf) {
		return nil, -1, false
	}

	// Look for the next newline
	idx := bytes.IndexByte(s.buf[s.readPos:], '\n')
	if idx < 0 {
		return nil, -1, false
	}

	startPos := s.readPos
	end := s.readPos + idx + 1
	line = append(line, s.buf[s.readPos:end]...)
	s.readPos = end

	return line, startPos, true
}

func (s *stream) WriteAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errNegativeOffset
	}

	endPos := int(off) + len(p)

	s.mu.Lock()
	defer s.mu.Unlock()
	// Grow buffer if needed
	if endPos > len(s.buf) {
		newBuf := make([]byte, endPos)
		copy(newBuf, s.buf)
		s.buf = newBuf
	}

	// Overwrite
	copy(s.buf[off:endPos], p)
	return len(p), nil
}

func (s *stream) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errNegativeOffset
	}

	if int(off) >= len(s.buf) {
		return 0, io.EOF
	}

	s.mu.RLock()
	n = copy(p, s.buf[off:])
	s.mu.RUnlock()
	s.mu.Lock()
	s.readPos += n
	s.mu.Unlock()
	if n < len(p) {
		return n, io.EOF
	}

	return
}

func (s *stream) MarshalZerologObject(e *zerolog.Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e.Bytes("buf", s.buf).Int("len", len(s.buf)).Int("readPos", s.readPos)
}
