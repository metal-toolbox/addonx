package reconciler

import (
	"go.uber.org/zap"
)

// Reconciler reconciles with downstream system
type Reconciler struct {
	client *ReconcileClient
	logger *zap.Logger
	queue  string
}

// Option is a functional configuration option
type Option func(r *Reconciler)

// WithLogger sets logger
func WithLogger(l *zap.Logger) Option {
	return func(r *Reconciler) {
		r.logger = l
	}
}

// WithClient sets api client
// TODO: Update with your client information
func WithClient(c *ReconcileClient) Option {
	return func(r *Reconciler) {
		r.client = c
	}
}

// WithQueue sets nats queue for events
func WithQueue(q string) Option {
	return func(r *Reconciler) {
		r.queue = q
	}
}

// New returns a new reconciler
func New(opts ...Option) *Reconciler {
	rec := Reconciler{
		logger: zap.NewNop(),
	}

	for _, opt := range opts {
		opt(&rec)
	}

	rec.logger.Debug("creating new reconciler")

	return &rec
}
