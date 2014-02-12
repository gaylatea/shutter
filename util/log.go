package util

import (
	"fmt"
	"github.com/mgutz/ansi"
	"io"
	"time"
)

// Implementations of **Logger** provide an interface to emit data
// to various sources.
type Logger interface {
	Info(s string) error
	Warning(s string) error
	Error(s string) error
	Debug(s string) error
	Success(s string) error

	Infof(format string, a ...interface{}) error
	Warningf(format string, a ...interface{}) error
	Errorf(format string, a ...interface{}) error
	Debugf(format string, a ...interface{}) error
	Successf(format string, a ...interface{}) error

	SetDebugOutput(debugState bool) error
}

// **ColourizedOutputLogger** is an implementation of **Logger** that
// outputs to the console with various colour levels to assist in
// debugging.
type ColourizedOutputLogger struct {
	levelColours map[string]func(string) string
	levels       []string
	outputTarget io.Writer
	debugEnabled bool
}

// Create a new **ColourizedOutputLogger** and cache colours.
func NewColourizedOutputLogger(w io.Writer) (*ColourizedOutputLogger, error) {
	// Create a cache of the colour functions for speed.
	colours := map[string]func(string) string{
		"info":    ansi.ColorFunc("white"),
		"debug":   ansi.ColorFunc("cyan"),
		"warning": ansi.ColorFunc("yellow"),
		"error":   ansi.ColorFunc("red"),
		"succeed": ansi.ColorFunc("green"),
	}

	return &ColourizedOutputLogger{
		levelColours: colours,
		outputTarget: w,
		debugEnabled: false,
	}, nil
}

// Enable debug output from the logger, which is usually off.
func (csl *ColourizedOutputLogger) SetDebugOutput(debugState bool) error {
	csl.debugEnabled = debugState
	return nil
}

// Delegate the logging statements to the right handler.
func (csl *ColourizedOutputLogger) Info(s string) error {
	return csl.emit("INFO   ", s, "info")
}

func (csl *ColourizedOutputLogger) Infof(format string, a ...interface{}) error {
	output := fmt.Sprintf(format, a...)
	return csl.emit("INFO   ", output, "info")
}

func (csl *ColourizedOutputLogger) Warning(s string) error {
	return csl.emit("WARNING", s, "warning")
}

func (csl *ColourizedOutputLogger) Warningf(format string, a ...interface{}) error {
	output := fmt.Sprintf(format, a...)
	return csl.emit("WARNING", output, "warning")
}

func (csl *ColourizedOutputLogger) Error(s string) error {
	return csl.emit("FAILURE", s, "error")
}

func (csl *ColourizedOutputLogger) Errorf(format string, a ...interface{}) error {
	output := fmt.Sprintf(format, a...)
	return csl.emit("FAILURE", output, "error")
}

func (csl *ColourizedOutputLogger) Debug(s string) error {
	if !csl.debugEnabled {
		return nil
	}

	return csl.emit("DEBUG  ", s, "debug")
}

func (csl *ColourizedOutputLogger) Debugf(format string, a ...interface{}) error {
	if !csl.debugEnabled {
		return nil
	}

	output := fmt.Sprintf(format, a...)
	return csl.emit("DEBUG  ", output, "debug")
}

func (csl *ColourizedOutputLogger) Success(s string) error {
	return csl.emit("SUCCESS", s, "succeed")
}

func (csl *ColourizedOutputLogger) Successf(format string, a ...interface{}) error {
	output := fmt.Sprintf(format, a...)
	return csl.emit("SUCCESS", output, "succeed")
}

// **ColourizedOutputLogger.emit** outputs the standard log format
// with the selected colours and prefix.
func (csl *ColourizedOutputLogger) emit(prefix string, msg string, level string) error {
	// Log output always needs a datetime.
	now := time.Now().Format(time.RFC822Z)
	output := fmt.Sprintf("[%s][%s] %s\n", prefix, now, msg)

	colourizedByteOutput := []byte(csl.levelColours[level](output))
	_, err := csl.outputTarget.Write(colourizedByteOutput)
	return err
}
