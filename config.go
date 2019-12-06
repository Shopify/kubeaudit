package kubeaudit

// Config includes options set by the user to modify the behaviour of auditors
type Config interface {
	Namespace() string
}

// implements Config
type config struct {
	// The namespace scope to audit
	namespace string
}

func (c *config) Namespace() string {
	return c.namespace
}
