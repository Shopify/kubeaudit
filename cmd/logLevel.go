package cmd

// Log levels
const (
	_ = iota
	Error
	Warn
	Info
	Debug
)

// KubeauditLogLevel is the default log level to be used by the logger. All log events with this log level and above
// will be logged.
var KubeauditLogLevel = Info

// KubeauditLogLevels represents an enum for the supported log levels.
var KubeauditLogLevels = map[string]int{"ERROR": Error, "WARN": Warn, "INFO": Info}
