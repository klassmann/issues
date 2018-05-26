package main

import (
	"fmt"
)

type ForegroundColor int

const (
	FgBlack ForegroundColor = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
	FgDefault
)

type Attribute int

const (
	AttrReset Attribute = iota
	AttrBold
	AttrDim
)

type TermFormat struct {
	Color ForegroundColor
	Attr  Attribute
}

func (t TermFormat) Quote(s string) string {
	f := "\033"
	f += "[%dm"
	var args []interface{}
	args = append(args, t.Color)

	if t.Attr != AttrReset {
		f += "\033[%dm"
		args = append(args, t.Attr)
	}

	f += "%s"
	args = append(args, s)

	f += "\033[%dm" // Reset
	args = append(args, AttrReset)

	return fmt.Sprintf(f, args...)
}

func (t TermFormat) Printf(format string, v ...interface{}) (int, error) {
	return fmt.Printf(t.Quote(fmt.Sprintf(format, v...)))
}
