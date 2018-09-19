package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

// DebugHook is a log hook intended to be used for debug logging.
type DebugHook struct{}

// NewDebugHook creates a new DebugHook.
func NewDebugHook() *DebugHook {
	hook := &DebugHook{}
	return hook
}

// Fire is called when a log event is triggered having a log level
// specified by the Levels method.
func (hook *DebugHook) Fire(entry *logrus.Entry) error {
	_, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.ErrorLevel:
		debugPrint()
		return nil
	default:
		return nil
	}
}

// Levels returns the log levels for which DebugHook.Fire
// should be called. This method is called when the hook is fist added
// to a logger instance.
func (hook *DebugHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
