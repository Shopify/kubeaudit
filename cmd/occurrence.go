package cmd

type Occurrence struct {
	kind    int    // or int? just needs to represent  {debug, log, warn, error}
	id      int    // KubeAuditInfo, ErrorImageTagMissing ...
	message string // the message that currently is in the printResultX function, which would go away with the introduction of this
}
