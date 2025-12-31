package main

import (
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

	stream *stream // our stream of data unescaped

	tracker *AckTracker
	prevAck *message.Ack

	ttl      *time.Timer // our expire timer
	lastSeen time.Time   // just to get a time stamp

	log zerolog.Logger
}

func NewSession(id int, addr net.Addr, pc net.PacketConn, logger zerolog.Logger) *Session {
	ss := &Session{
		PacketConn: pc,
		ID:         id,
		Addr:       addr,

		stream:  &stream{},
		tracker: newAckTracker(),

		lastSeen: time.Now(),
		log:      logger.With().Int("session", id).Stringer("addr", addr).Logger(),
	}
	ss.ttl = time.AfterFunc(sessionTimeout, func() {
		ss.log.Info().EmbedObject(ss).Msg("session expired")
		if err := ss.Close(); err != nil {
			ss.log.Err(err).EmbedObject(ss).Msg("err sending close")
		}
	})

	ss.log.Info().EmbedObject(ss).Msg("new session")
	return ss
}

// InsertData adds data to the internal buffer, it unescapes any escape characters
func (ss *Session) InsertData(m *message.Data) error {
	defer ss.resetTimer()
	defer ss.checkLines()

	log := ss.log.With().Object("msg", m).Logger()

	if (m.Pos + len(m.Data)) <= ss.stream.Len() {
		log.Debug().EmbedObject(ss).
			Msgf("already read data to POS %v", ss.stream.Len())
		return ss.sendPrevAck()
	}

	// unescape data payload
	data := unescapeData(m.Data)
	if _, err := ss.stream.WriteAt(data, int64(m.Pos)); err != nil {
		log.Error().Err(err).Msg("err writing at")
	}

	log.Debug().
		EmbedObject(ss).
		Dict("escapedData", zerolog.Dict().
			Bytes("data", data).
			Int("len", len(data)),
		).
		Msg("inserted data")

	return ss.Ack(ss.stream.Len())
}

func (ss *Session) Ack(length int) error {
	m := &message.Ack{SessionID: ss.ID, Length: length}
	ss.prevAck = m
	return ss.sendAck(m)
}

func (ss *Session) Close() error {
	m := &message.Close{SessionID: ss.ID}
	_, err := ss.WriteTo(m.Marshal(), ss.Addr)
	ss.log.Info().EmbedObject(ss).Msg("sent close msg")
	return err
}

func (ss *Session) sendAck(m *message.Ack) error {
	_, err := ss.WriteTo(m.Marshal(), ss.Addr)
	ss.log.Info().EmbedObject(ss).Object("msg", m).Msg("sent ack msg")
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
	ss.log.Info().EmbedObject(ss).Object("msg", m).Msg("sent data msg")
	return err
}

// we should check for lines, and chunk out full lines into our stream
func (ss *Session) checkLines() {
	for {
		line, pos, ok := ss.stream.Readline()
		if !ok {
			break
		}

		b := reverseBytes(line)
		ss.log.Debug().
			Bytes("reverse", b).
			Bytes("line", line).
			Msg("check lines")

		if err := ss.sendData(pos, b); err != nil {
			ss.log.Err(err).EmbedObject(ss).Msg("sending data msg")
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
		Object("stream", ss.stream).
		Dur("lastSeen", time.Since(ss.lastSeen).Round(zerolog.DurationFieldUnit))
}
