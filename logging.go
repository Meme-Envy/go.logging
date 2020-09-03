package logging

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	logLevelEnv = "LOG_LEVEL"

	levelNone      = "NONE"
	levelMandatory = "NOTICE"
	levelError     = "ERROR"
	levelWarn      = "WARN"
	levelInfo      = "INFO"
	levelDebug     = "DEBUG"
	levelTrace     = "TRACE"

	// DefaultFlags are the default go.logging flags
	DefaultFlags = log.LstdFlags | log.Lshortfile
)

var (
	//LogLevelConfig is the level set by the env
	logLevelConfig = os.Getenv(logLevelEnv)

	// DefaultLogger is the package level logger - exposed on purpose for library user to have control over native log functions (SetOutput, etc)
	DefaultLogger = New(os.Stderr, "", DefaultFlags)

	// ErrInvalidLevel is returned by SetLevel when an invalid log level is set
	ErrInvalidLevel = errors.New("invalid log level")
)

//LogLevel is a logging level
type LogLevel uint8

//Prefix returns the prefix string for the log level
func (ll LogLevel) Prefix() string {
	return fmt.Sprintf("%s - ", ll)
}

//Logging levels
const (
	LNone = LogLevel(iota)
	LMandatory
	LError
	LWarn
	LInfo
	LDebug
	LTrace
)

//LevelLogger logs according to indicated log level
type LevelLogger struct {
	// unexported to ensure log level is always valid
	currentLevel LogLevel
	*log.Logger
}

//New gets a new LevelLogger
func New(dest io.Writer, prefix string, flag int) *LevelLogger {
	logger := &LevelLogger{
		Logger: log.New(dest, prefix, flag),
	}
	logger.SetLogLevel(logLevelConfig)
	return logger
}

//CurrentLevel returns the current level
func (l LevelLogger) CurrentLevel() LogLevel {
	return l.currentLevel
}

// SetLevel is the same as SetLogLevel but accepts a LogLevel constant instead of string name
func (l *LevelLogger) SetLevel(ll LogLevel) (err error) {
	if ll > LTrace {
		// invalid log level set to Error to match SetLogLevel
		ll = LError
		err = ErrInvalidLevel
	}
	l.currentLevel = ll
	return
}

//String implements String interface
func (ll LogLevel) String() string {
	switch ll {
	case LNone:
		return levelNone
	case LMandatory:
		return levelMandatory
	case LWarn:
		return levelWarn
	case LInfo:
		return levelInfo
	case LDebug:
		return levelDebug
	case LTrace:
		return levelTrace
	}
	// Error is default log level
	return levelError
}

//SetLogLevel allows applications to change the log level with a reload instead of restart
func (l *LevelLogger) SetLogLevel(level string) {
	switch strings.ToUpper(level) {
	case levelNone:
		l.currentLevel = LNone
	case levelError:
		l.currentLevel = LError
	case levelWarn:
		l.currentLevel = LWarn
	case levelInfo:
		l.currentLevel = LInfo
	case levelDebug:
		l.currentLevel = LDebug
	case levelTrace:
		l.currentLevel = LTrace
	default: // Error is the default log level and includes mandatory
		l.currentLevel = LError
	}
}

//log sends the format and the params to the underlying logger
func (l LevelLogger) logf(level LogLevel, formattedString string, params ...interface{}) {
	if level <= l.currentLevel {
		l.Output(3, level.Prefix()+fmt.Sprintf(formattedString, params...))
	}
}

func (l LevelLogger) log(level LogLevel, params ...interface{}) {
	if level <= l.currentLevel {
		l.Output(3, level.Prefix()+fmt.Sprint(params...))
	}
}

//Error log
func (l LevelLogger) Error(err error) {
	l.log(LError, err)
}

//Errorf log
func (l LevelLogger) Errorf(formattedString string, params ...interface{}) {
	l.logf(LError, formattedString, params...)
}

//Mandatoryp always logs regardless of logging level
func (l LevelLogger) Mandatoryp(params ...interface{}) {
	l.log(LMandatory, params...)
}

//Mandatory always logs regardless of logging level
func (l LevelLogger) Mandatory(formattedString string, params ...interface{}) {
	l.logf(LMandatory, formattedString, params...)
}

//Warnp log
func (l LevelLogger) Warnp(params ...interface{}) {
	l.log(LWarn, params...)
}

//Warn log
func (l LevelLogger) Warn(formattedString string, params ...interface{}) {
	l.logf(LWarn, formattedString, params...)
}

//Infop log
func (l LevelLogger) Infop(params ...interface{}) {
	l.log(LInfo, params...)
}

//Info log
func (l LevelLogger) Info(formattedString string, params ...interface{}) {
	l.logf(LInfo, formattedString, params...)
}

//Debugp log
func (l LevelLogger) Debugp(params ...interface{}) {
	l.log(LDebug, params...)
}

//Debug log
func (l LevelLogger) Debug(formattedString string, params ...interface{}) {
	l.logf(LDebug, formattedString, params...)
}

//Tracep log
func (l LevelLogger) Tracep(params ...interface{}) {
	l.log(LTrace, params...)
}

//Trace log
func (l LevelLogger) Trace(formattedString string, params ...interface{}) {
	l.logf(LTrace, formattedString, params...)
}

/*
package level convenience functions for the default logger:
*/

//SetLogLevel sets the log level
func SetLogLevel(level string) {
	DefaultLogger.SetLogLevel(level)
}

//CurrentLevel returns the current log level as constant
func CurrentLevel() LogLevel {
	return DefaultLogger.CurrentLevel()
}

//SetLevel sets the log level by constant
func SetLevel(level LogLevel) error {
	return DefaultLogger.SetLevel(level)
}

//Fatal prints and exits 1
// must copy paste from go src.
func Fatal(v ...interface{}) {
	DefaultLogger.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

//Fatalf is equivalent to calling Errorf followed by os.Exit(1)
// copied from go src
func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

//Panic is equivalent to calling Errorf followed by panic(params)
// copy pasted from go src. must do this for line numbers to work.
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	DefaultLogger.Output(2, s)
	panic(s)
}

//Panicf is equivalent to calling Errorf followed by panic(params)
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	DefaultLogger.Output(2, s)
	panic(s)
}

//Error log
func Error(err error) {
	DefaultLogger.log(LError, err)
}

//Errorf log
func Errorf(formattedString string, params ...interface{}) {
	DefaultLogger.logf(LError, formattedString, params...)
}

//Mandatoryp always logs regardless of logging level (unless NONE)
func Mandatoryp(params ...interface{}) {
	DefaultLogger.log(LMandatory, params...)
}

//Mandatory always logs regardless of logging level (unless NONE)
func Mandatory(formattedString string, params ...interface{}) {
	DefaultLogger.logf(LMandatory, formattedString, params...)
}

//Warnp log
func Warnp(params ...interface{}) {
	DefaultLogger.log(LWarn, params...)
}

//Warn log
func Warn(formattedString string, params ...interface{}) {
	DefaultLogger.logf(LWarn, formattedString, params...)
}

//Infop log
func Infop(params ...interface{}) {
	DefaultLogger.log(LInfo, params...)
}

//Info log
func Info(formattedString string, params ...interface{}) {
	DefaultLogger.logf(LInfo, formattedString, params...)
}

//Debugp log
func Debugp(params ...interface{}) {
	DefaultLogger.log(LDebug, params...)
}

//Debug log
func Debug(formattedString string, params ...interface{}) {
	DefaultLogger.logf(LDebug, formattedString, params...)
}

//Tracep log
func Tracep(params ...interface{}) {
	DefaultLogger.log(LTrace, params...)
}

//Trace log
func Trace(formattedString string, params ...interface{}) {
	DefaultLogger.logf(LTrace, formattedString, params...)
}
