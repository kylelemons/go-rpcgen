package webrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"net/http"
	"path"

	"code.google.com/p/goprotobuf/proto"
)

const (
	DefaultRPCPath = "/_webRPC_"
)

const (
	JSON     Protocol = "application/json"
	ProtoBuf Protocol = "application/protobuf"
)

type Protocol string

func (p Protocol) String() string { return string(p) }

type Handler func(*Call) error

type Call struct {
	http.ResponseWriter
	*http.Request

	ContentType string
}

func (c *Call) ReadRequest(pb interface{}) error {
	ctype := Protocol(c.ContentType)
	switch ctype {
	case JSON:
		return json.NewDecoder(c.Request.Body).Decode(pb)
	case ProtoBuf:
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return fmt.Errorf("webrpc: readall: %s", err)
		}
		return proto.Unmarshal(body, pb)
	}
	return fmt.Errorf("webrpc: read: %s: bad content type", ctype)
}

func (c *Call) WriteResponse(pb interface{}) error {
	ctype := c.ContentType
	switch ctype {
	case "application/json":
		return json.NewEncoder(c.ResponseWriter).Encode(pb)
	case "application/protobuf":
		body, err := proto.Marshal(pb)
		if err != nil {
			return fmt.Errorf("webrpc: readall: %s", err)
		}
		if _, err := c.ResponseWriter.Write(body); err != nil {
			return fmt.Errorf("webrpc: write: %s", err)
		}
		return nil
	}
	return fmt.Errorf("webrpc: write: %s: bad content type", ctype)
}

type ServeMux map[string]Handler

func (m ServeMux) Handle(path string, handler Handler) error {
	path = DefaultRPCPath + path
	if _, exist := m[path]; exist {
		return fmt.Errorf("webrpc: handler already registered for %q", path)
	}
	m[path] = handler
	return nil
}

func (m ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	handler, found := m[r.URL.Path]
	if !found {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, r.URL.Path+" is not a registered RPC", http.StatusNotFound)
		return
	}

	ctype := r.Header.Get("Content-Type")
	switch ctype {
	case "application/json":
	case "application/protobuf":
	default:
		http.Error(w, ctype+": invalid content type (must be application/json or application/protobuf)", http.StatusUnsupportedMediaType)
		return
	}

	c := &Call{
		ResponseWriter: w,
		Request:        r,
		ContentType:    ctype,
	}
	if err := handler(c); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, r.URL.Path+": "+err.Error(), http.StatusInternalServerError)
		return
	}
}

var DefaultServeMux = ServeMux{}

func Post(protocol Protocol, base *url.URL, method string, in, out interface{}) error {
	var b []byte
	var err error

	url := *base
	url.Path = path.Join(url.Path, DefaultRPCPath, method)

	switch protocol {
	case JSON:
		b, err = json.Marshal(in)
	case ProtoBuf:
		b, err = proto.Marshal(in)
	default:
		return fmt.Errorf("webrpc.post: %s: bad content type", protocol)
	}
	if err != nil {
		return fmt.Errorf("webrpc.post: %s: marshal: %s", protocol, err)
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("webrpc.post: newrequest: %s", err)
	}
	req.Header.Set("Content-Type", protocol.String())
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(b)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("webrpc.post: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("webrpc.post: readall: %s", err)
		}
		return fmt.Errorf("webrc.post: %s: %s", resp.Status, bytes.TrimSpace(b))
	}

	switch protocol {
	case JSON:
		err = json.NewDecoder(resp.Body).Decode(out)
	case ProtoBuf:
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("webrpc.post: readall: %s", err)
		}
		err = proto.Unmarshal(b, out)
	}
	if err != nil {
		return fmt.Errorf("webrpc.post: %s: unmarshal: %s in %q", protocol, err, b)
	}

	return nil
}

func ListenAndServe(addr string, mux ServeMux) error {
	if mux == nil {
		mux = DefaultServeMux
	}
	return http.ListenAndServe(addr, mux)
}

func init() {
	http.Handle(DefaultRPCPath+"/", DefaultServeMux)
}
