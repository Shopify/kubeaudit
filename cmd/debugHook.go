package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type DebugHook struct{}

func NewDebugHook() *DebugHook {
	hook := &DebugHook{}
	return hook
}

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

func (hook *DebugHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
