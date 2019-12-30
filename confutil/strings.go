package confutil

import (
	"errors"
	"strconv"
	"unicode/utf8"
)

type RuneString string

func (s RuneString) Postprocess() (err error) {
	_, err = s.Rune()
	return
}

func (s RuneString) Rune() (r rune, err error) {
	if s == "" {
		return
	}
	r, n := utf8.DecodeRuneInString(string(s))
	if n != len(s) {
		err = errors.New("more than one character provided")
		return
	}
	if r == utf8.RuneError {
		err = errors.New("invalid character")
	}
	return
}

type EscapedRuneString string

func (e *EscapedRuneString) Postprocess() error {
	q := `"` + string(*e) + `"`
	s, err := strconv.Unquote(q)
	if err != nil {
		return err
	}
	rs := RuneString(s)
	_, err = rs.Rune()
	if err != nil {
		return err
	}
	*e = EscapedRuneString(s)
	return nil
}

func (e EscapedRuneString) Rune() rune {
	r, _ := RuneString(e).Rune()
	return r
}
