// Generate the built-in policies described by security/tomoyo/Makefile.
package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

var policies = []string{"profile", "exception_policy", "domain_policy", "manager", "stat"}

func generate(inputs []string) (string, error) {
	data := map[string][]byte{}
	for _, input := range inputs {
		name, path, ok := strings.Cut(input, "=")
		if !ok || !slices.Contains(policies, name) {
			return "", fmt.Errorf("invalid policy input %q", input)
		}
		if _, exists := data[name]; exists {
			return "", fmt.Errorf("duplicate policy %q", name)
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		data[name] = contents
	}
	var out strings.Builder
	for _, name := range policies {
		fmt.Fprintf(&out, "static char tomoyo_builtin_%s[] __initdata =\n", name)
		contents := data[name]
		// scripts/bin2c always emits a final line, even for empty input.
		for {
			length := min(len(contents), 16)
			out.WriteString("\t\"")
			for _, b := range contents[:length] {
				fmt.Fprintf(&out, "\\x%02x", b)
			}
			out.WriteString("\"\n")
			contents = contents[length:]
			if length < 16 {
				break
			}
		}
		out.WriteString(";\n")
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: tomoyopolicy OUTPUT [POLICY=INPUT ...]")
		os.Exit(2)
	}
	content, err := generate(os.Args[2:])
	if err == nil {
		err = os.WriteFile(os.Args[1], []byte(content), 0o644)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
