package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// generate follows cmd_gen-kdb in kernel/debug/kdb/Makefile.
func generate(input io.Reader) (string, error) {
	var out strings.Builder
	out.WriteString("#include <linux/stddef.h>\n#include <linux/init.h>\n")
	scanner := bufio.NewScanner(input)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.Trim(line, " \t") == "" {
			continue
		}
		fmt.Fprintf(&out, "static __initdata char kdb_cmd%d[] = \"%s\\n\";\n", count, strings.ReplaceAll(line, "\"", "\\\""))
		count++
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	out.WriteString("extern char *kdb_cmds[]; char __initdata *kdb_cmds[] = {\n")
	for i := 0; i < count; i++ {
		fmt.Fprintf(&out, "  kdb_cmd%d,\n", i)
	}
	out.WriteString("  NULL\n};\n")
	return out.String(), nil
}

func run(input, output string) error {
	file, err := os.Open(input)
	if err != nil {
		return err
	}
	defer file.Close()
	source, err := generate(file)
	if err != nil {
		return err
	}
	return os.WriteFile(output, []byte(source), 0o644)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: kdbcmds INPUT OUTPUT")
		os.Exit(2)
	}
	if err := run(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintf(os.Stderr, "kdbcmds: %v\n", err)
		os.Exit(1)
	}
}
