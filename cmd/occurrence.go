package cmd

// An Occurrence represents a potential security issue. There may be multiple Occurrences per resource and audit.
type Occurrence struct {
	kind      int    // represent  {debug, log, warn, error}
	id        int    // KubeAuditInfo, ErrorImageTagMissing ...
	message   string // just the message
	container string // name of the container
	metadata  Metadata
}
