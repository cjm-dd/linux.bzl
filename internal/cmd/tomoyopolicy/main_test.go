package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "exception_policy.conf.default")
	for _, input := range []string{"", "\x00\xff\"\\\n", strings.Repeat("a", 16), strings.Repeat("b", 17)} {
		if err := os.WriteFile(path, []byte(input), 0o644); err != nil {
			t.Fatal(err)
		}
		got, err := generate([]string{"exception_policy=" + path})
		if err != nil {
			t.Fatal(err)
		}
		if strings.Count(got, `\x`) != len(input) {
			t.Fatalf("policy bytes lost or added: %s", got)
		}
		position := -1
		for _, name := range policies {
			declaration := "static char tomoyo_builtin_" + name + "[] __initdata =\n"
			next := strings.Index(got, declaration)
			if next <= position {
				t.Fatalf("policy order or declaration differs: %s", got)
			}
			position = next
			if name != "exception_policy" && !strings.Contains(got, declaration+"\t\"\"\n;\n") {
				t.Fatalf("missing policy must be an empty string: %s", got)
			}
		}
		if len(input) == 16 && !strings.Contains(got, "\t\"\"\n;\nstatic char tomoyo_builtin_domain_policy") {
			t.Fatalf("exact chunks need a final empty string: %s", got)
		}
	}
	for _, inputs := range [][]string{{"unknown=" + path}, {"profile"}, {"profile=" + path, "profile=" + path}, {"profile=" + path + ".missing"}} {
		if _, err := generate(inputs); err == nil {
			t.Fatalf("accepted invalid inputs %v", inputs)
		}
	}
}
