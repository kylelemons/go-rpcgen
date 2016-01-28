package client

import (
	"net/rpc"
)

type ConnectionPool struct {
	conns chan *rpc.Client

	New func() (*rpc.Client, error)
}

func (p *ConnectionPool) open() *rpc.Client {
	for {
		c, err := p.New()

		if IsTimeoutError(err) {
			continue
		}

		if err != nil {
			logMessage("[go-rpcgen/connection_pool] Error opening connection. %v", err)
			return nil
		}

		return c
	}

	panic("unreachable")
}

func (p *ConnectionPool) Get() *rpc.Client {
	select {
	case c := <-p.conns:
		return c
	default:
		return p.open()
	}
}

func (p *ConnectionPool) Put(c *rpc.Client) {
	select {
	case p.conns <- c:
		// Do nothing.
		return
	default:
		c.Close()
	}
}

func (p *ConnectionPool) Close() {
	close(p.conns)

	for c := range p.conns {
		c.Close()
	}
}

func NewConnectionPool(capacity int) *ConnectionPool {
	p := new(ConnectionPool)
	p.conns = make(chan *rpc.Client, capacity)
	return p
}
