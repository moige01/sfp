package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"sugud0r.dev/sfp/internal/loggin"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var sessions Sessions = nil
var l net.Listener = nil

func init() {
	signalChanel := make(chan os.Signal, 1)
	exit := make(chan int)

	signal.Notify(signalChanel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	closeAll := func(s os.Signal) {
		loggin.Info.Println("Handling signal", s)

		if sessions != nil {
			sessions.CloseAll()
		}

		if l != nil {
			if err := l.Close(); err != nil {
				loggin.Error.Fatalln(err)
			}
		}

		os.Exit(1)
	}

	go func() {
		for {
			switch s := <-signalChanel; s {
			case syscall.SIGINT:
				closeAll(s)
			case syscall.SIGTERM:
				closeAll(s)
			case syscall.SIGQUIT:
				closeAll(s)

			default:
				loggin.Info.Println("Unknown signal.")
				exit <- 1
			}
		}
	}()
}

func CreateServer() {
	l, err := net.Listen(CONN_TYPE, fmt.Sprintf("%v:%v", CONN_HOST, CONN_PORT))

	if err != nil {
		loggin.Error.Fatalln(err.Error())
	}

	defer l.Close()

	loggin.Info.Printf("Listening on %v:%v", CONN_HOST, CONN_PORT)

	for {
		loggin.Info.Print("Waiting for new connection...")

		conn, err := l.Accept()

		if err != nil {
			loggin.Error.Fatalln(err.Error())
		}

		loggin.Info.Print("New connection received from ", conn.RemoteAddr())

		session := NewSession(conn)

		sessions.Append(session)

		go session.Read()
	}
}
