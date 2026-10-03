package main

import (
	"strings"
	"testing"
)

func TestGenerateHeaders(t *testing.T) {
	for _, tc := range []struct {
		kind   string
		inputs []string
		want   string
	}{
		{"capability", []string{"#define CAP_CHOWN 0\n#define CAP_NET_ADMIN 12\n#define CAP_LAST_CAP CAP_NET_ADMIN\n#define CAP_FS_MASK 99\n"}, `static const char *const capability_names[] = {
[0] = "chown",
[12] = "net_admin",
};
#define AA_SFS_CAPS_MASK "chown net_admin"
`},
		{"net", []string{"#define AF_UNIX 1\n#define AF_LOCAL 1\n#define AF_ROUTE AF_NETLINK\n#define AF_INET 2\n#define AF_MAX 3\n", "\tSOCK_STREAM\t= 1,\n\tSOCK_DGRAM\t= 2,\n"}, `static const char *address_family_names[] = {
[1] = "unix",
[2] = "inet",
};
#define AA_SFS_AF_MASK "unix inet"
static const char *sock_type_names[] = {
[1] = "stream",
[2] = "dgram",
};
`},
		{"rlim", []string{"#define RLIMIT_CPU 0\n# define RLIMIT_STACK 3 /* comment */\n#define RLIM_NLIMITS 4\n"}, `static const char *const rlim_names[RLIM_NLIMITS] = {
[RLIMIT_CPU] = "cpu",
[RLIMIT_STACK] = "stack",
};
static const int rlim_map[RLIM_NLIMITS] = {
RLIMIT_CPU,
RLIMIT_STACK,
};
#define AA_SFS_RLIMIT_MASK "cpu stack"
`},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			got, err := generate(tc.kind, tc.inputs)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("header = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRejectMissingDefinitions(t *testing.T) {
	if _, err := generate("capability", []string{"/* unknown format */"}); err == nil || !strings.Contains(err.Error(), "no entries") {
		t.Fatalf("error = %v, want missing entries", err)
	}
	if _, err := generate("net", []string{"#define AF_INET 2\n"}); err == nil {
		t.Fatal("missing socket header accepted")
	}
}
