package counters

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCountingWriter(t *testing.T) {
	var buf bytes.Buffer
	ncw, pn := CountingWriter(&buf)
	wn, err := fmt.Fprintf(ncw, "hello world")
	if err != nil {
		t.Fatalf("buf.WriteString error:%s", err)
	}
	if wn != int(*pn) {
		t.Errorf("expect %d, got %d", wn, *pn)
	}
	t.Logf("write %d bytes", *pn)
}
