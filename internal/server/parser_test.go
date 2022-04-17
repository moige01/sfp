package server_test

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	. "sugud0r.dev/sfp/internal/server"
)

func checkError(t *testing.T, got error, want error) {
	t.Helper()

	if got == nil || !errors.Is(got, want) {
		t.Errorf("Error: got %#v; want %#v", got, want)
	}
}

func TestParser(t *testing.T) {
	var parser *SfpCommandParser = NewSfpCommandParser()

	t.Run("Errors", func(t *testing.T) {
		t.Run("Should return ErrUnknownCommand if command does not exists", func(t *testing.T) {
			fooComand := "FOOCOMMAND"

			err := parser.Parse(fooComand)

			checkError(t, err, ErrUnkownCommand)
		})

		t.Run("Should return ErrEmptyCommand if command is empty", func(t *testing.T) {
			emptyCommand := ""

			err := parser.Parse(emptyCommand)

			checkError(t, err, ErrEmptyCommand)
		})

		t.Run("Should return ErrMalformedCommand if command does not have enough arguments", func(t *testing.T) {
			t.Run("SUSCRIBE", func(t *testing.T) {
				malformedCommand := "SUSCRIBE"

				err := parser.Parse(malformedCommand)

				checkError(t, err, ErrMalformedCommand)
			})

			t.Run("SEND", func(t *testing.T) {
				malformedCmommand := "SEND 1 23"

				err := parser.Parse(malformedCmommand)

				checkError(t, err, ErrMalformedCommand)
			})
		})
	})

	t.Run("SfpCommandParser implement CommandParser", func(t *testing.T) {
		commandParserI := reflect.TypeOf((*CommandParser)(nil)).Elem()

		sfpC := reflect.TypeOf((*SfpCommandParser)(nil))

		if !sfpC.Implements(commandParserI) {
			t.Error("SfpCommandParser should implement CommandParser interface")
		}
	})

	t.Run("Should parse given command", func(t *testing.T) {
		t.Run("Parsing BYE", func(t *testing.T) {
			byeCommand := "BYE"

			if err := parser.Parse(byeCommand); err != nil {
				t.Errorf("Should not return an error. Given %#v", err)
			}

			if parser.GetKind() != Bye {
				t.Errorf("Command should be type Bye. Given %#v", parser.GetKind())
			}
		})

		t.Run("Parsing SUSCRIBE", func(t *testing.T) {
			suscribeCommand := "SUSCRIBE 1"

			if err := parser.Parse(suscribeCommand); err != nil {
				t.Errorf("Should not return an error. Given %#v", err)
			}

			if parser.GetKind() != Suscribe {
				t.Errorf("Command should be type Suscribe. Given %#v", parser.GetKind())
			}

			args := parser.GetArgs()

			if argsLen := len(args); argsLen != 1 {
				t.Errorf("Erroneous amount of arguments. Given %#v", argsLen)
			}

			v, ok := args["channel_id"]

			if !ok {
				t.Errorf("channel_id should exists as argument")
			}

			if _, err := strconv.Atoi(v); err != nil {
				t.Errorf("Value of channel_id shoueld be a number. Given %v", v)
			}
		})
	})
}
