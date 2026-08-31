package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// run takes its I/O as arguments, so these tests drive the whole command
// in-process: no subprocess, no redirection of os.Stdout.

// exercise calls run with the given standard input (which may be nil)
// and arguments, capturing what it writes.
func exercise(stdin io.Reader, args ...string) (stdout, stderr string, code int) {
	var out, errs bytes.Buffer
	if stdin == nil {
		stdin = strings.NewReader("")
	}
	code = run(args, stdin, &out, &errs)
	return out.String(), errs.String(), code
}

// testdata reads a file from the testdata directory.
func testdata(t *testing.T, name string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func TestDeduplicates(t *testing.T) {
	want := testdata(t, "expected.txt")

	tests := []struct {
		name  string
		args  []string
		stdin io.Reader
	}{
		{name: "file argument", args: []string{filepath.Join("testdata", "input.txt")}},
		{name: "standard input", stdin: strings.NewReader(testdata(t, "input.txt"))},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, code := exercise(tc.stdin, tc.args...)
			if code != 0 {
				t.Fatalf("exit code = %d, stderr = %q", code, stderr)
			}
			if stdout != want {
				t.Errorf("stdout =\n%s\nwant:\n%s", stdout, want)
			}
		})
	}
}

func TestEmptyInput(t *testing.T) {
	stdout, stderr, code := exercise(nil)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr)
	}
	if stdout != "" {
		t.Errorf("stdout = %q, want it empty", stdout)
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantStderr string
	}{
		{"too many arguments", []string{"a", "b"}, "usage:"},
		{"missing file", []string{"no-such-file.txt"}, "no-such-file.txt"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, code := exercise(nil, tc.args...)
			if code != 1 {
				t.Errorf("exit code = %d, want 1", code)
			}
			if stdout != "" {
				t.Errorf("stdout = %q, want it empty", stdout)
			}
			if !strings.Contains(stderr, tc.wantStderr) {
				t.Errorf("stderr = %q, want it to contain %q", stderr, tc.wantStderr)
			}
		})
	}
}
