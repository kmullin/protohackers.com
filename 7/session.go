package main

import (
	"bytes"
	"net"
	"time"

	"github.com/kmullin/protohackers.com/7/message"
	"github.com/rs/zerolog"
)

const (
	sessionTimeout = 60 * time.Second
	reXmitTimeout  = time.Duration(3 * time.Second)
)

type Session struct {
	net.PacketConn

	ID   int
	Addr net.Addr

	recvBuf []byte // the input buffer
	sendBuf []byte // the output buffer
	Pos     int    // the position in the buffer we last transmitted

	prevAck *message.Ack

	ttl      *time.Timer // our expire timer
	lastSeen time.Time   // just to get a time stamp

	log zerolog.Logger
}

func NewSession(id int, addr net.Addr, pc net.PacketConn, logger zerolog.Logger) *Session {
	ss := &Session{
		ID:         id,
		Addr:       addr,
		PacketConn: pc,
		log:        logger.With().Int("session", id).Stringer("addr", addr).Logger(),
	}
	ss.ttl = time.AfterFunc(sessionTimeout, func() {
		err := ss.Close()
		if err != nil {
			ss.log.Error().Err(err).Msg("err sending close")
		}
	})

	ss.log.Debug().Msg("new session")
	return ss
}

// AddData adds data to the internal buffer, it unescapes any escape characters
func (ss *Session) AddData(m *message.Data) error {
	ss.resetTimer()

	// unescape data payload
	data := unescapeData(m.Data)

	endPos := m.Pos + len(data)

	// Grow buffer if needed
	if endPos > len(ss.recvBuf) {
		newBuf := make([]byte, endPos)
		copy(newBuf, ss.recvBuf)
		ss.recvBuf = newBuf
	}

	// Overwrite
	copy(ss.recvBuf[m.Pos:endPos], data)

	ss.log.Debug().
		Bytes("buf", ss.recvBuf).
		Int("len", len(data)).
		Int("pos", m.Pos).
		Msg("inserted data")
	err := ss.Ack(len(data))

	// check for lines
	ss.checkLines()
	return err
}

func (ss *Session) Ack(length int) error {
	m := &message.Ack{SessionID: ss.ID, Length: length}
	ss.prevAck = m
	return ss.sendAck(m)
}

func (ss *Session) Close() error {
	m := &message.Close{SessionID: ss.ID}
	_, err := ss.WriteTo(m.Marshal(), ss.Addr)
	ss.log.Debug().Msg("sent close msg")
	return err
}

func (ss *Session) sendAck(m *message.Ack) error {
	_, err := ss.WriteTo(m.Marshal(), ss.Addr)
	ss.log.Debug().EmbedObject(m).Msg("sent ack msg")
	return err
}

func (ss *Session) sendPrevAck() error {
	if ss.prevAck == nil {
		ss.prevAck = &message.Ack{SessionID: ss.ID, Length: 0}
	}
	return ss.sendAck(ss.prevAck)
}

func (ss *Session) sendData(pos int, data []byte) error {
	m := &message.Data{SessionID: ss.ID, Pos: pos, Data: escapeData(data)}
	_, err := ss.WriteTo(m.Marshal(), ss.Addr)
	ss.log.Debug().EmbedObject(m).Msg("sent data msg")
	return err
}

func (ss *Session) checkLines() {
	for line := range bytes.Lines(ss.recvBuf) {
		b := reverseBytes(line)
		ss.log.Debug().Bytes("reverse", b).Bytes("line", line).Msg("check lines")
		err := ss.sendData(0, b)
		if err != nil {
			ss.log.Error().Err(err).Msg("sending data msg")
		}
	}
}

func (ss *Session) resetTimer() {
	ss.ttl.Reset(sessionTimeout)
	ss.lastSeen = time.Now()
}

func (ss *Session) MarshalZerologObject(e *zerolog.Event) {
	e.Int("session", ss.ID).
		Stringer("addr", ss.Addr).
		Int("len", len(ss.recvBuf)).
		Str("buf", string(ss.recvBuf)).
		Stringer("last", time.Since(ss.lastSeen))
}
