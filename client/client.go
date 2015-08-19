package client

import (
	"errors"
	"github.com/bradhe/go-rpcgen/codec"
	"math"
	"net"
	"net/rpc"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrConnectionFailure   = errors.New("failed to connect")
	ErrClosed              = errors.New("closed")
	ErrInvalidPoolObject   = errors.New("invalid pool object")
	ErrPermanentlyShutdown = errors.New("permenantly shutdown")
)

const (
	DefaultRetryCount = 6
)

type Client struct {
	wg sync.WaitGroup

	addr     string
	pool     sync.Pool
	shutdown int32
}

func backoff(i int) time.Duration {
	if i < 1 {
		i = 1
	}

	ms := int(math.Exp2(float64(i)))
	return time.Duration(ms) * time.Millisecond
}

func (c *Client) Close() {
	// If someone else called this, we'll just wait a bit for it all to close down.
	if atomic.LoadInt32(&c.shutdown) > 0 {
		c.wg.Wait()
		return
	}

	atomic.SwapInt32(&c.shutdown, 1)
	c.wg.Wait()

	// By adding a second number here that means that we're completely shut down.
	atomic.AddInt32(&c.shutdown, 1)

	// Now let's close down all of the connections in the poo in the pool.
	for {
		client, ok := c.pool.Get().(*rpc.Client)

		if !ok {
			// we closed them all
			break
		}

		client.Close()
	}
}

func (c *Client) create() interface{} {
	// If the service is shutdown, let's wait to kill it all.
	if atomic.LoadInt32(&c.shutdown) > 1 {
		return nil
	}

	conn, err := net.Dial("tcp", c.addr)

	// If the connection failed, there's nothin' we can really do about it.
	if err != nil {
		return nil
	}

	return rpc.NewClientWithCodec(codec.NewClientCodec(conn))
}

func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	// If we're shut down, let's tell the user.
	if c.shutdown > 0 {
		return ErrClosed
	}

	// Signal that something is going on.
	c.wg.Add(1)
	defer c.wg.Done()

	// Number of times we've retried
	var retry int

	for {
		client, ok := c.pool.Get().(*rpc.Client)

		if !ok {
			return ErrInvalidPoolObject
		}

		if client == nil {
			return ErrConnectionFailure
		}

		err := client.Call(serviceMethod, args, reply)

		// No error, so let's relinquish this back to the pool and get outta here.
		if err == nil {
			c.pool.Put(client)
			break
		}

		// If we got here, let's see what type of error it is.
		if err == rpc.ErrShutdown {
			retry += 1

			if retry > DefaultRetryCount {
				return ErrPermanentlyShutdown
			}

			// Let's try again!
			time.Sleep(backoff(retry))
		} else {
			// This means err != nil, so we just report the error.
			return err
		}
	}

	// We win the day!
	return nil
}

func (c *Client) doCall(call *rpc.Call, serviceMethod string, args interface{}, reply interface{}, done chan *rpc.Call) {
	err := c.Call(serviceMethod, args, reply)
	call.Error = err
	done <- call
}

func (c *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *rpc.Call) *rpc.Call {
	call := new(rpc.Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.Done = done

	// If we're shut down, let's tell the user.
	if c.shutdown > 0 {
		call.Error = ErrClosed
		return call
	}

	// If we made it here, we're good.
	go c.doCall(call, serviceMethod, args, reply, done)
	return call
}

func NewClient(addr string) *Client {
	c := new(Client)
	c.addr = addr
	c.pool.New = c.create
	return c
}
