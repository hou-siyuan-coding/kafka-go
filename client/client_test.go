package client

import (
	"bytes"
	"testing"
)

func TestCutToLastMessage(t *testing.T) {
	res := []byte("100\n101\n10")
	wantTruncated, wantRest := []byte("100\n101\n"), []byte("10")
	truncated, rest, err := cutToLastMessage(res)

	if err != nil {
		t.Errorf("cutToLastMessage(%q): got error %v, want no errors:", string(res), err)
	}

	if !bytes.Equal(truncated, wantTruncated) || !bytes.Equal(rest, wantRest) {
		t.Errorf("cutToLastMessage(%q): got: %q, %q, want: %q, %q", string(res), string(truncated), string(rest),
			string(wantTruncated), string(wantRest))
	}
}
