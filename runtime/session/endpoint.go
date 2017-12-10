package session

import (
	"sync"

	"github.com/nickng/scribble-go-runtime/runtime/transport"
)

type Endpoint struct {
	Id       int
	NumRoles int

	connWg sync.WaitGroup // Counts initiated connections from this Endpoint.

	// guards Conn
	ConnMu sync.RWMutex
	Conn   map[string][]transport.Channel
}

func NewEndpoint(roleId, numRoles int, conn map[string][]transport.Channel) *Endpoint {
	return &Endpoint{
		Id:       roleId,
		NumRoles: numRoles,
		Conn:     conn,
	}
}

// CheckConnection ensures connections initiated (by Accept)
// in Endpoint e are fully established.
func (e *Endpoint) CheckConnection() {
	e.connWg.Wait()
}
