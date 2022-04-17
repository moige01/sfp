package server

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"sugud0r.dev/sfp/internal/loggin"
)

type CommandType string

const (
	Suscribe CommandType = "SUSCRIBE"
	Send     CommandType = "SEND"
	Bye      CommandType = "BYE"
)

var (
	ErrMalformedCommand = errors.New("command malformed")
	ErrUnkownCommand    = errors.New("unknown command or no implemented")
	ErrEmptyCommand     = errors.New("empty data to parse")
)

type CommandParser interface {
	Parse(data string) error
	GetArgs() map[string]string
}

type SfpCommandParser struct {
	kind CommandType
	args map[string]string
}

func NewSfpCommandParser() *SfpCommandParser {
	return &SfpCommandParser{}
}

func (c *SfpCommandParser) validateInt(data string) error {
	if _, err := strconv.Atoi(data); err != nil {
		loggin.Debug.Printf("%#v don't look like a number.", data)

		return ErrMalformedCommand
	}

	return nil
}

func (c *SfpCommandParser) Parse(data string) error {
	if data == "" {
		loggin.Error.Print("Trying to parse empty data")

		return ErrEmptyCommand
	}

	c.args = make(map[string]string)

	loggin.Debug.Printf("Parsing: %#v", data)

	fields := strings.FieldsFunc(data, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})

	loggin.Debug.Printf("Resulting fields: %#v", fields)

	switch strings.ToUpper(fields[0]) {
	case string(Suscribe):
		if len(fields) < 2 {
			loggin.Error.Print("No enough arguments for SUSCRIBE command")

			return ErrMalformedCommand
		}

		if err := c.validateInt(fields[1]); err != nil {
			return err
		}

		c.args["channel_id"] = fields[1]
		c.kind = Suscribe
	case string(Send):
		if len(fields) < 4 {
			loggin.Error.Print("No enough arguments for SEND command")

			return ErrMalformedCommand
		}

		if err := c.validateInt(fields[1]); err != nil {
			return err
		}

		if err := c.validateInt(fields[2]); err != nil {
			return err
		}

		c.args["channel_id"] = fields[1]
		c.args["content_size"] = fields[2]
		c.args["filename"] = fields[3]

		c.kind = Suscribe
	case string(Bye):
		c.kind = Bye
	default:
		loggin.Error.Printf("UNKNOWN commmand. %v", fields[0])

		return ErrUnkownCommand
	}

	loggin.Debug.Printf("Result: %#v", *c)

	return nil
}

func (c *SfpCommandParser) GetArgs() map[string]string {
	return c.args
}

func (c *SfpCommandParser) GetKind() CommandType {
	return c.kind
}
