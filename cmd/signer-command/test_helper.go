package signercommand

import (
	"bytes"
	"os"
)

// CaptureOutput captures stdout output from a function
func CaptureOutput(f func()) []byte {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	f()

	os.Stdout = stdout
	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.Bytes()
}

// CaptureOutputWithErr captures both stdout and stderr
func CaptureOutputWithErr(f func()) (stdout, stderr []byte) {
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	os.Stdout = wOut
	os.Stderr = wErr

	f()

	os.Stdout = oldStdout
	os.Stderr = oldStderr
	wOut.Close()
	wErr.Close()

	var bufOut, bufErr bytes.Buffer
	bufOut.ReadFrom(rOut)
	bufErr.ReadFrom(rErr)

	return bufOut.Bytes(), bufErr.Bytes()
}

// GetProjectRoot returns the project root directory
func GetProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}
