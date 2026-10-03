package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func encode(data []byte, name string) (string, error) {
	if !regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(name) {
		return "", fmt.Errorf("invalid C identifier %q", name)
	}
	var out strings.Builder
	fmt.Fprintf(&out, "const char %s[] =\n", name)
	for len(data) > 0 {
		chunk := data[:min(len(data), 16)]
		out.WriteString("\t\"")
		for _, b := range chunk {
			fmt.Fprintf(&out, "\\x%02x", b)
		}
		out.WriteString("\"\n")
		data = data[len(chunk):]
	}
	fmt.Fprintf(&out, "\t\"\";\n#include <linux/types.h>\nconst size_t %s_size = sizeof(%s) - 1;\n", name, name)
	return out.String(), nil
}

func run(input, output, name string) error {
	data, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	text, err := encode(data, name)
	if err != nil {
		return err
	}
	return os.WriteFile(output, []byte(text), 0o644)
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: bin2c INPUT OUTPUT NAME")
		os.Exit(2)
	}
	if err := run(os.Args[1], os.Args[2], os.Args[3]); err != nil {
		fmt.Fprintf(os.Stderr, "bin2c: %v\n", err)
		os.Exit(1)
	}
}
