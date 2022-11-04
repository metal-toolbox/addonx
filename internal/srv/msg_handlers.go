package srv

import "github.com/nats-io/nats.go"

// AddonMessageHandler handles messages from Nats event bus
type AddonMessageHandler interface {
	// Handle messages from nats
	Handle(m *nats.Msg)
}
