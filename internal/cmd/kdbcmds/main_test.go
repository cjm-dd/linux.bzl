package main

import (
	"strings"
	"testing"
)

func TestGenerateCommands(t *testing.T) {
	got, err := generate(strings.NewReader("# ignored\n \t\ndefcmd dump \"\" \"a dump\"\n  -bt\nendefcmd\n"))
	if err != nil {
		t.Fatal(err)
	}
	want := `#include <linux/stddef.h>
#include <linux/init.h>
static __initdata char kdb_cmd0[] = "defcmd dump \"\" \"a dump\"\n";
static __initdata char kdb_cmd1[] = "  -bt\n";
static __initdata char kdb_cmd2[] = "endefcmd\n";
extern char *kdb_cmds[]; char __initdata *kdb_cmds[] = {
  kdb_cmd0,
  kdb_cmd1,
  kdb_cmd2,
  NULL
};
`
	if got != want {
		t.Fatalf("generated commands = %q, want %q", got, want)
	}
}

func TestGenerateEmptyCommands(t *testing.T) {
	got, err := generate(strings.NewReader("# no commands\n"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "static __initdata") || !strings.HasSuffix(got, "{\n  NULL\n};\n") {
		t.Fatalf("empty command table = %q", got)
	}
}
