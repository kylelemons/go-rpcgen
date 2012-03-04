package codec

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"code.google.com/p/goprotobuf/proto"
)

type InvalidRequest struct{}

func TestWriteProto(t *testing.T) {
	tests := []struct {
		Desc string
		Out  interface{}
		In   interface{}
	}{
		{
			Desc: "Test the zero size proto",
			Out:  &InvalidRequest{},
			In:   &InvalidRequest{},
		},
	}

	for _, test := range tests {
		size := make([]byte, binary.MaxVarintLen64)
		data, _ := proto.Marshal(test.In)
		size = size[:binary.PutUvarint(size, uint64(len(data)))]
		data = append(size, data...)

		t.Logf("Marshal(%#v) = %q", test.Out, data)

		b := new(bytes.Buffer)
		// test encode
		if err := WriteProto(b, test.Out); err != nil {
			t.Errorf("WriteProto(%#v) - %s", test.Out, err)
		} else if got, want := b.String(), string(data); got != want {
			t.Errorf("WriteProto(%#v) wrote %q, want %q", test.Out, got, want)
		}

		b.Reset()
		b.Write(data)
		// test decode
		if err := ReadProto(b, test.In); err != nil {
			t.Errorf("ReadProto(%q) - %s", data, err)
		} else if got, want := test.In, test.Out; !reflect.DeepEqual(got, want) {
			t.Errorf("ReadProto(%q) wrote %#v, want %#v", data, got, want)
		}
	}
}
