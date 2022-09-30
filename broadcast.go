package broadcaster

type broadcaster[Data any] struct {
	input chan Data
	reg   chan chan<- Data
	unreg chan chan<- Data

	outputs map[chan<- Data]bool
}

func (b *broadcaster[Data]) broadcast(m Data) {
	for ch := range b.outputs {
		ch <- m
	}
}

func (b *broadcaster[Data]) run() {
	for {
		select {
		case m := <-b.input:
			b.broadcast(m)
		case ch, ok := <-b.reg:
			if !ok {
				return
			}

			b.outputs[ch] = true
		case ch := <-b.unreg:
			delete(b.outputs, ch)
		}
	}
}

func (b *broadcaster[Data]) Register(newch chan<- Data) chan<- Data {
	b.reg <- newch

	return newch
}

func (b *broadcaster[Data]) Unregister(newch chan<- Data) {
	b.unreg <- newch
}

func (b *broadcaster[Data]) Close() error {
	close(b.reg)
	close(b.unreg)
	return nil
}

// Submit an item to be broadcast to all listeners.
func (b *broadcaster[Data]) Submit(m Data) {
	if b != nil {
		b.input <- m
	}
}

// TrySubmit attempts to submit an item to be broadcast, returning
// true iff it the item was broadcast, else false.
func (b *broadcaster[Data]) TrySubmit(m Data) bool {
	if b == nil {
		return false
	}
	select {
	case b.input <- m:
		return true
	default:
		return false
	}
}
