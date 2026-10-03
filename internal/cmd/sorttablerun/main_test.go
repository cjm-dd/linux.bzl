package main

import (
	"slices"
	"testing"
)

func TestSorttableArgsFollowKernelInterface(t *testing.T) {
	for _, tc := range []struct {
		source string
		want   []string
	}{
		{"usage: sorttable vmlinux...", []string{"vmlinux"}},
		{"usage: sorttable [-s nm-file] vmlinux...", []string{"-s", "symbols", "vmlinux"}},
	} {
		if got, err := sorttableArgs(tc.source, "symbols", "vmlinux"); err != nil || !slices.Equal(got, tc.want) {
			t.Errorf("sorttableArgs(%q) = %v, %v; want %v", tc.source, got, err, tc.want)
		}
	}
	if _, err := sorttableArgs("", "symbols", "vmlinux"); err == nil {
		t.Fatal("unknown interface accepted")
	}
}
