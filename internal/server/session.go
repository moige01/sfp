package server

import (
	"bufio"
	"io"
	"net"
	"sync/atomic"

	"sugud0r.dev/sfp/internal/loggin"
)

var atmid int64 = 0

type Sessions []*Session

// Represent a session in the current active connection.
// id is an atomic varaible to handle uniques sessions.
// conn is a TCP connection.
type Session struct {
	id   int64
	conn net.Conn
}

// Create a new session. An unique id would be assigned
func NewSession(conn net.Conn) *Session {
	atomic.AddInt64(&atmid, 1)

	loggin.Info.Printf("Creating new sesion with ID %v for address %v", atmid, conn.RemoteAddr())

	return &Session{id: atmid, conn: conn}
}

// Close connection for this session.
func (s *Session) Close() {
	loggin.Info.Printf("Closing connection #%v...\n", s.id)

	s.conn.Write([]byte("BYE!"))

	if err := s.conn.Close(); err != nil {
		loggin.Error.Print(err)
	}
}

func (s *Session) Read() {
	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(s.conn)
	)

READ_LOOP:
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			break READ_LOOP
		case nil:
			loggin.Debug.Printf("Receive: %#v", data)
		default:
			loggin.Error.Printf("Receive data failed: %#v", err)
		}
	}
}

func (s *Session) Write(b []byte) {
	if _, err := s.conn.Write(b); err != nil {
		loggin.Error.Print(err)
	}
}

func (s *Sessions) Append(session *Session) {
	*s = append(*s, session)
}

func (s *Sessions) Empty() {
	*s = nil
}

func (s *Sessions) CloseAll() {
	loggin.Info.Println("Closing all sessions")

	for _, s := range *s {
		s.Close()
	}

	s.Empty()
}
