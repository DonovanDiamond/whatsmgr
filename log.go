package whatsmgr

import (
	"github.com/rs/zerolog"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Logger struct {
	zerolog.Logger
	Module string
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.Warn().Str("_module", l.Module).Msgf(msg, args...)
}
func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.Error().Str("_module", l.Module).Msgf(msg, args...)
}
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.Debug().Str("_module", l.Module).Msgf(msg, args...)
}
func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.Trace().Str("_module", l.Module).Msgf(msg, args...)
}
func (l *Logger) Sub(module string) waLog.Logger {
	var subLogger Logger
	subLogger.Logger = l.Logger
	subLogger.Module = module
	return &subLogger
}
