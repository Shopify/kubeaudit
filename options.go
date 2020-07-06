package kubeaudit

import (
	log "github.com/sirupsen/logrus"
)

// Option is used to specify the behaviour of Kubeaudit Auditor
type Option func(*Kubeaudit) error

// WithLogger specifies the log formatter to use
func WithLogger(formatter log.Formatter) Option {
	return func(_ *Kubeaudit) error {
		log.SetFormatter(formatter)
		return nil
	}
}

func (a *Kubeaudit) parseOptions(opts []Option) error {
	for _, opt := range opts {
		if err := opt(a); err != nil {
			return err
		}
	}
	return nil
}
