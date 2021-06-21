package encoder

import (
	"bytes"
	"github.com/golang/snappy"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func TestSszNetworkEncoder_BufferedReader(t *testing.T) {
	r := make([]byte, 10)
	bufR := snappy.NewReader(bytes.NewBuffer(r))
	ptr := reflect.ValueOf(bufR).Pointer()
	bufReaderPool.Put(bufR)

	r2 := make([]byte, 10)
	rdr := newBufferedReader(bytes.NewBuffer(r2))

	nPtr := reflect.ValueOf(rdr).Pointer()
	assert.Equal(t, ptr, nPtr, "invalid pointer value")
}

func TestSszNetworkEncoder_BufferedWriter(t *testing.T) {
	r := make([]byte, 10)
	bufR := snappy.NewBufferedWriter(bytes.NewBuffer(r))
	ptr := reflect.ValueOf(bufR).Pointer()
	bufWriterPool.Put(bufR)

	r2 := make([]byte, 10)
	rdr := newBufferedWriter(bytes.NewBuffer(r2))

	nPtr := reflect.ValueOf(rdr).Pointer()
	assert.Equal(t, ptr, nPtr, "invalid pointer value")
}

