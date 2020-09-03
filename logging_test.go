package logging

import (
	"errors"
	"os"
	"testing"
)

func logstuff(l *LevelLogger) {
	err := errors.New("test log error")
	l.Error(err)
	l.Errorf("a formatted error %+v", err)

	l.Mandatory("a mandatory log, next line is mandatory without format string %s", err)
	l.Mandatoryp(err)

	l.Warn("some message formatted %s", err)
	l.Warnp(err)

	l.Info("some message formatted %s", err)
	l.Infop(err)

	l.Debug("some message formatted %s", err)
	l.Debugp(err)

	l.Trace("some message formatted %s", err)
	l.Tracep(err)
	// test stdlib
	l.Printf("test printf stdlib bypass levels %d", 1)
}

func TestLogging(t *testing.T) {
	for _, logger := range []*LevelLogger{
		DefaultLogger,
		New(os.Stderr, "custom pre-prefix ", DefaultFlags),
	} {
		for _, l := range []string{
			levelNone,
			levelError,
			levelWarn,
			levelInfo,
			levelDebug,
			levelTrace,
			"INVALIDLOGLEVEL",
		} {
			logger.Printf("\n----- TESTING LOG LEVEL %s -----", l)
			logger.SetLogLevel(l)
			logstuff(logger)
		}

		logger.Printf("\n----- TESTING CONSTANT SET LEVEL -----")
		logger.SetLevel(LError)
		for _, l := range []LogLevel{
			LNone,
			LError,
			LWarn,
			LInfo,
			LDebug,
			LTrace,
			100,
		} {
			logger.Printf("\n----- TESTING LOG LEVEL %s -----", l)
			logger.SetLevel(l)
			logstuff(logger)
		}

		if logger.CurrentLevel() != LError {
			t.Errorf("expected log level %s, got level %s", LError, logger.CurrentLevel())
		}

		err := logger.SetLevel(254)
		if !errors.Is(err, ErrInvalidLevel) {
			t.Errorf("expected error for invalid level %v, got %v", ErrInvalidLevel, err)
		}

		if logger.CurrentLevel() != LError {
			t.Errorf("expected log level %s, got level %s", LError, logger.CurrentLevel())
		}
	}
}
