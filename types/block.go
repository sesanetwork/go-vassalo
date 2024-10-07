package types

import "github.com/sesanetwork/go-vassalo/hash"

// Block is a part of an ordered chain of batches of events.
type Block struct {
	Event    hash.Event
	Cheaters Cheaters
}
