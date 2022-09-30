/*
Package broadcast provides pubsub of messages over channels.
A provider has a Broadcaster into which it Submits messages and into
which subscribers Register to pick up those messages.
*/
package broadcaster

// The Broadcaster interface describes the main entry points to
// broadcasters.
type Broadcaster[Data any] interface {
	// Register a new channel to receive broadcasts
	Register(chan<- Data)
	// Unregister a channel so that it no longer receives broadcasts.
	Unregister(chan<- Data)
	// Shut this broadcaster down.
	Close() error
	// Submit a new object to all subscribers
	Submit(Data)
	// Try Submit a new object to all subscribers return false if input chan is fill
	TrySubmit(Data) bool
}

// NewBroadcaster creates a new broadcaster with the given input
// channel buffer length.
func NewBroadcaster[Data any](buflen int) Broadcaster[Data] {
	b := &broadcaster[Data]{
		input:   make(chan Data, buflen),
		reg:     make(chan chan<- Data),
		unreg:   make(chan chan<- Data),
		outputs: make(map[chan<- Data]bool),
	}

	go b.run()

	return b
}
