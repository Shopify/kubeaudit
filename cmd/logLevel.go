package cmd

const (
	_ = iota
	Error
	Warn
	Info
	Debug
)

var KubeauditLogLevel = Info
var KubeauditLogLevels = map[string]int{"ERROR": Error, "WARN": Warn, "INFO": Info, "DEBUG": Debug}
