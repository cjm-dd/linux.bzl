package main

import (
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	for _, data := range [][]byte{nil, {0, 0xff, '"', '\\'}, make([]byte, 17)} {
		got, err := encode(data, "kexec_purgatory")
		if err != nil {
			t.Fatal(err)
		}
		if count := strings.Count(got, `\x`); count != len(data) {
			t.Errorf("encoded %d bytes, want %d", count, len(data))
		}
		if !strings.Contains(got, "const size_t kexec_purgatory_size = sizeof(kexec_purgatory) - 1;") {
			t.Errorf("size must exclude the string terminator: %s", got)
		}
	}
	got, _ := encode([]byte{0, 0xff, '"', '\\'}, "image")
	if !strings.Contains(got, `"\x00\xff\x22\x5c"`) {
		t.Errorf("unexpected byte encoding: %s", got)
	}
	if _, err := encode(nil, "bad;name"); err == nil {
		t.Fatal("invalid C identifier accepted")
	}
}
