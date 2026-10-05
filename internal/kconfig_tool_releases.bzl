"""Release metadata for prebuilt kconfig repository-rule tools.

Each archive must extract the requested host executables at its root.
"""

visibility("//...")

KCONFIG_TOOL_VERSION = "2ffc5cdffb2ec082505acb6a24f075e3af3f33f6"

_RELEASE_BASE_URL = "https://ddartifacts.jfrog.io/artifactory/ddlabs-cityhopper/bazelbuild/linux.bzl/kconfig-{version}".format(
    version = KCONFIG_TOOL_VERSION,
)

KCONFIG_TOOL_RELEASES = {
    "darwin_amd64": struct(
        integrity = "sha256-Ww6Tee85lm0FnmVlhkkUxj0exk09W0NEAmWYjLrGQLA=",
        urls = ["{}/kconfig-darwin-amd64.tar.zst".format(_RELEASE_BASE_URL)],
    ),
    "darwin_arm64": struct(
        integrity = "sha256-WWW47SCoNDRbTMXrc7YLCKEIXbrg8NVZ3sWumPGcgj8=",
        urls = ["{}/kconfig-darwin-arm64.tar.zst".format(_RELEASE_BASE_URL)],
    ),
    "linux_amd64": struct(
        integrity = "sha256-LVkp/zFUY6IFHwoag1trUrsKRECPUO1y4zK9k79OLaA=",
        urls = ["{}/kconfig-linux-amd64.tar.zst".format(_RELEASE_BASE_URL)],
    ),
    "linux_arm64": struct(
        integrity = "sha256-eUuACUYlJymrUHBXn/uPzvITRydspYlyHVY/q6DIqNo=",
        urls = ["{}/kconfig-linux-arm64.tar.zst".format(_RELEASE_BASE_URL)],
    ),
    "windows_amd64": struct(
        integrity = "sha256-Q4mFhZcg9yD+O8JSweC8ZFfD/Bk4eGCWi+Cjc29LONE=",
        urls = ["{}/kconfig-windows-amd64.tar.zst".format(_RELEASE_BASE_URL)],
    ),
    "windows_arm64": struct(
        integrity = "sha256-TCsIjb5I93uktycsV0QlYhSd/kj8HZ9+mF1JaMiGvEU=",
        urls = ["{}/kconfig-windows-arm64.tar.zst".format(_RELEASE_BASE_URL)],
    ),
}
