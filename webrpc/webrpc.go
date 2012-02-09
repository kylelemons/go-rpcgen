package webrpc

import (
	"fmt"
	"net/url"
	"net/http"
)

type Call struct {
	*http.Request
}

func (c *Call) ReadProto(pb interface{}) error { return nil }
func (c *Call) WriteProto(pb interface{}) error { return nil }

type Handler func(*Call) error

type ServeMux map[string]Handler

func (m ServeMux) Handle(path string, handler Handler) error {
	if _, exist := m[path]; exist {
		return fmt.Errorf("webrpc: handler already registered for %q", path)
	}
	m[path] = handler
	return nil
}

var DefaultServeMux = ServeMux{}

func Post(base *url.URL, path string, in, out interface{}) error {
	return nil
}
