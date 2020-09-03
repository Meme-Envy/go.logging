# LOGGING

### v2 changes:

Color is gone - depended on CGo. Use a colorizer, e.g. ccze:
https://linux.die.net/man/1/ccze

For historical reasons, Warn, Info, Debug, Trace all accept the Printf syntax,
without the "f" suffix. Errorf is the exception. Logging v1 used this syntax
and it is kept for an easier upgrade to v2 (without changing all the logging
calls in each application). Warnp, Infop, Debugp, Tracep are added to log
without a format string at the specified log level.

The default logger is exposed and directly embeds the standard lib Logger
struct. This means all functions on the standard lib logger are available. See
the go std lib documentation on "log". Access it with `logging.DefaultLogger`.
This means that `logging.DefaultLogger.Printf` will work - and yes it will
bypass the level check and behave like "MANDATORY".

Package level convenience functions are added to use the DefaultLogger object.
This maintains API compatibility with logging v1 package functions.

The "New()" function no longer accepts a second writer (was for stdout, logs do
not belong in stdout).

Logging is imported with semantic import versioning:

github.com/erock530/logging/v2

# To use:

To set the debug level at runtime, export the "LOG_LEVEL" environment variable set.
Ex: `export LOG_LEVEL=WARN`

Levels (low to high):
* NONE: must be explicitly set - turns off EVERYTHING
* MANDATORY: logging.MANDATORY [Default]
* ERROR: logging.Error, logging.Errorf [Default]
* WARN: logging.Warn
* INFO: logging.Info
* DEBUG: logging.Debug
* TRACE: logging.Trace

MANDATORY and Error warnings are active by default, so do not require any LOG_LEVEL to be set.

The logging.Error function ONLY accepts an error object.
The logging.Errorf function behaves like a string formatter.  It will accept any number of strings or objects as parameters

If a warning level is set, then that level and all lower levels will be emitted.  If the INFO level is set, then Info, Error, Errorf and MANDATORY messages will all be emitted.

# Other Functions

## Fatal, Fatalf, Panic, Panicf

Inherited from https://golang.org/pkg/log/#Fatal

