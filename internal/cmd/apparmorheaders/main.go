package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type entry struct {
	index string
	name  string
}

func numericEntries(input, pattern string, excluded ...string) []entry {
	var entries []entry
	matcher := regexp.MustCompile(pattern)
	for _, line := range strings.Split(input, "\n") {
		skip := false
		for _, name := range excluded {
			skip = skip || strings.Contains(line, name)
		}
		if match := matcher.FindStringSubmatch(line); !skip && match != nil {
			entries = append(entries, entry{index: match[2], name: strings.ToLower(match[1])})
		}
	}
	return entries
}

func writeTable(out *strings.Builder, declaration, mask string, entries []entry) error {
	if len(entries) == 0 {
		return fmt.Errorf("no entries found for %s", declaration)
	}
	fmt.Fprintf(out, "%s = {\n", declaration)
	var names []string
	for _, entry := range entries {
		fmt.Fprintf(out, "[%s] = %q,\n", entry.index, entry.name)
		names = append(names, entry.name)
	}
	out.WriteString("};\n")
	if mask != "" {
		fmt.Fprintf(out, "#define %s %q\n", mask, strings.Join(names, " "))
	}
	return nil
}

// generate implements the name-table transformations in security/apparmor/Makefile.
func generate(kind string, inputs []string) (string, error) {
	want := 1
	if kind == "net" {
		want = 2
	}
	if len(inputs) != want {
		return "", fmt.Errorf("%s requires %d inputs", kind, want)
	}
	var out strings.Builder
	switch kind {
	case "capability":
		entries := numericEntries(inputs[0], `^#define[ \t]+CAP_([A-Z0-9_]+)[ \t]+([0-9]+)`, "CAP_FS_MASK")
		if err := writeTable(&out, "static const char *const capability_names[]", "AA_SFS_CAPS_MASK", entries); err != nil {
			return "", err
		}
	case "net":
		families := numericEntries(inputs[0], `^#define[ \t]+AF_([A-Z0-9_]+)[ \t]+([0-9]+)`, "AF_MAX", "AF_LOCAL", "AF_ROUTE")
		if err := writeTable(&out, "static const char *address_family_names[]", "AA_SFS_AF_MASK", families); err != nil {
			return "", err
		}
		sockets := numericEntries(inputs[1], `^\tSOCK_([A-Z0-9_]+)[\t]+=[ \t]+([0-9]+)`)
		if err := writeTable(&out, "static const char *sock_type_names[]", "", sockets); err != nil {
			return "", err
		}
	case "rlim":
		var entries []entry
		for _, match := range regexp.MustCompile(`(?m)^# ?define[ \t]+(RLIMIT_([A-Z0-9_]+))`).FindAllStringSubmatch(inputs[0], -1) {
			entries = append(entries, entry{index: match[1], name: strings.ToLower(match[2])})
		}
		if err := writeTable(&out, "static const char *const rlim_names[RLIM_NLIMITS]", "", entries); err != nil {
			return "", err
		}
		out.WriteString("static const int rlim_map[RLIM_NLIMITS] = {\n")
		var names []string
		for _, entry := range entries {
			fmt.Fprintf(&out, "%s,\n", entry.index)
			names = append(names, entry.name)
		}
		fmt.Fprintf(&out, "};\n#define AA_SFS_RLIMIT_MASK %q\n", strings.Join(names, " "))
	default:
		return "", fmt.Errorf("unknown AppArmor header kind %q", kind)
	}
	return out.String(), nil
}

func run(kind, output string, paths []string) error {
	var inputs []string
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		inputs = append(inputs, string(data))
	}
	header, err := generate(kind, inputs)
	if err != nil {
		return err
	}
	return os.WriteFile(output, []byte(header), 0o644)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "usage: apparmorheaders KIND OUTPUT INPUT...")
		os.Exit(2)
	}
	if err := run(os.Args[1], os.Args[2], os.Args[3:]); err != nil {
		fmt.Fprintf(os.Stderr, "apparmorheaders: %v\n", err)
		os.Exit(1)
	}
}
