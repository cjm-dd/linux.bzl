"""Analysis tests for platform-first multi-architecture kernel facades."""

load("@bazel_skylib//lib:unittest.bzl", "analysistest", "asserts")
load("//internal:providers.bzl", "LinuxKernelInfo", "LinuxModuleSdkInfo")

visibility("private")

def _graph_fixture_impl(ctx):
    out = ctx.actions.declare_file(ctx.label.name + ".txt")
    ctx.actions.write(out, ctx.attr.arch + "\n")
    providers = [
        DefaultInfo(files = depset([out])),
        LinuxModuleSdkInfo(
            arch = ctx.attr.arch,
            module_symvers = out,
            modules = depset([out]),
            modules_builtin = out,
            modules_builtin_modinfo = out,
            modules_order = out,
        ),
        OutputGroupInfo(selected_profile = depset([out])),
    ]
    if not ctx.attr.sdk_only:
        providers.append(LinuxKernelInfo(
            arch = ctx.attr.arch,
            version = "test",
            kernel_release = out,
            image = out,
            vmlinux = out,
            config = out,
            system_map = out,
        ))
    return providers

multiarch_graph_fixture = rule(
    implementation = _graph_fixture_impl,
    attrs = {
        "arch": attr.string(mandatory = True),
        "sdk_only": attr.bool(),
    },
)

def _facade_provider_test_impl(ctx):
    env = analysistest.begin(ctx)
    target = analysistest.target_under_test(env)

    asserts.true(env, LinuxKernelInfo in target)
    asserts.true(env, LinuxModuleSdkInfo in target)
    asserts.true(env, OutputGroupInfo in target)
    asserts.equals(env, "armv7", target[LinuxKernelInfo].arch)
    asserts.equals(env, "armv7", target[LinuxModuleSdkInfo].arch)
    asserts.equals(env, 1, len(target[DefaultInfo].files.to_list()))
    asserts.equals(env, 1, len(target[OutputGroupInfo].selected_profile.to_list()))
    asserts.equals(
        env,
        target[DefaultInfo].files.to_list()[0],
        target[OutputGroupInfo].selected_profile.to_list()[0],
    )
    return analysistest.end(env)

facade_provider_test = analysistest.make(_facade_provider_test_impl)

def _projection_test_impl(ctx):
    env = analysistest.begin(ctx)
    target = analysistest.target_under_test(env)
    files = target[DefaultInfo].files.to_list()

    asserts.equals(env, 1, len(files))
    asserts.equals(env, ctx.attr.expected_file, files[0].basename)
    return analysistest.end(env)

projection_test = analysistest.make(
    _projection_test_impl,
    attrs = {"expected_file": attr.string(default = "armv7_graph.txt")},
)

def _sdk_provider_test_impl(ctx):
    env = analysistest.begin(ctx)
    target = analysistest.target_under_test(env)
    asserts.true(env, LinuxKernelInfo not in target)
    asserts.equals(env, "armv7", target[LinuxModuleSdkInfo].arch)
    asserts.equals(env, "armv7_sdk.txt", target[DefaultInfo].files.to_list()[0].basename)
    return analysistest.end(env)

sdk_provider_test = analysistest.make(_sdk_provider_test_impl)

def _selected_file_test_impl(ctx):
    env = analysistest.begin(ctx)
    files = analysistest.target_under_test(env)[DefaultInfo].files.to_list()

    asserts.equals(env, 1, len(files))
    asserts.equals(env, "multiarch_facade_test.bzl", files[0].basename)
    return analysistest.end(env)

selected_file_test = analysistest.make(_selected_file_test_impl)
