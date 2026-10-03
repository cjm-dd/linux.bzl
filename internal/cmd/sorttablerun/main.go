package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	config := flag.String("config", "", "resolved Linux .config")
	nm := flag.String("nm", "", "llvm-nm executable")
	sorttable := flag.String("sorttable", "", "sorttable executable")
	sorttableSource := flag.String("sorttable-source", "", "kernel scripts/sorttable.c")
	in := flag.String("in", "", "input vmlinux")
	nmOut := flag.String("nm_out", "", "output llvm-nm data")
	out := flag.String("out", "", "output vmlinux")
	flag.Parse()

	if *config == "" || *nm == "" || *in == "" || *nmOut == "" || *out == "" {
		fmt.Fprintln(os.Stderr, "-config, -nm, -in, -nm_out, and -out are required")
		os.Exit(2)
	}
	if err := run(*config, *nm, *sorttable, *sorttableSource, *in, *nmOut, *out); err != nil {
		fmt.Fprintf(os.Stderr, "sorttablerun: %v\n", err)
		os.Exit(1)
	}
}

func run(configPath, nmPath, sorttablePath, sorttableSourcePath, inPath, nmOutPath, outPath string) error {
	enabled, err := configEnabled(configPath, "CONFIG_BUILDTIME_TABLE_SORT")
	if err != nil {
		return err
	}
	if err := copyFile(inPath, outPath); err != nil {
		return err
	}
	if !enabled {
		return os.WriteFile(nmOutPath, nil, 0o644)
	}
	if sorttablePath == "" {
		return fmt.Errorf("CONFIG_BUILDTIME_TABLE_SORT=y requires sorttable")
	}
	source, err := os.ReadFile(sorttableSourcePath)
	if err != nil {
		return err
	}
	args, err := sorttableArgs(string(source), nmOutPath, outPath)
	if err != nil {
		return err
	}
	if len(args) > 1 {
		if err := writeCommandOutput(nmOutPath, nmPath, "-S", inPath); err != nil {
			return err
		}
	} else if err := os.WriteFile(nmOutPath, nil, 0o644); err != nil {
		return err
	}
	if err := os.Chmod(outPath, 0o600); err != nil {
		return err
	}
	cmd := exec.Command(sorttablePath, args...)
	cmd.Env = []string{}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func sorttableArgs(source, nmPath, imagePath string) ([]string, error) {
	// The -s option is also backported to stable kernels; inspect the tool's
	// own declared interface rather than inferring it from major/minor.
	if strings.Contains(source, "sorttable [-s nm-file]") {
		return []string{"-s", nmPath, imagePath}, nil
	}
	if strings.Contains(source, "sorttable vmlinux") {
		return []string{imagePath}, nil
	}
	return nil, fmt.Errorf("unrecognized sorttable interface")
}

func configEnabled(path, key string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == key+"=y" {
			return true, nil
		}
	}
	return false, scanner.Err()
}

func writeCommandOutput(path, executable string, args ...string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	cmd := exec.Command(executable, args...)
	cmd.Env = []string{}
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	runErr := cmd.Run()
	closeErr := out.Close()
	if runErr != nil {
		return runErr
	}
	return closeErr
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
