package api

import (
	"io"
	"log"

	"github.com/go-resty/resty/v2"
)

var _ resty.Logger = (*logger)(nil)

type logger struct {
	l *log.Logger
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.output("ERROR - RESTY - "+format, v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.output("WARN - RESTY -"+format, v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.output("DEBUG - RESTY"+format, v...)
}

func (l *logger) output(format string, v ...interface{}) {
	l.l.Printf(format, v...)
}

func createRestyLogger(out io.Writer) *logger {
	return &logger{l: log.New(out, "", log.LstdFlags)}
}
