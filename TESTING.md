# CLI Testing Guide for Signer

There are several approaches to test CLI commands in Go. Here's a comprehensive guide for testing your Cobra-based CLI.

## 1. **Manual Testing** (Easiest for quick verification)

```bash
# Test algorithms command
./bin/signer algorithms

# Test signing
./bin/signer sign -a ECDSA-P256-SHA256 -k my-key -i "test data"

# Test verification
./bin/signer verify -a ECDSA-P256-SHA256 -k my-key -i "test data" -s "SIGNATURE_HERE"

# Test help
./bin/signer --help
./bin/signer sign --help
./bin/signer verify --help
```

## 2. **Unit Testing with Cobra**

### Basic Structure

```go
package signercommand

import (
	"bytes"
	"testing"
)

func TestAlgorithmsCommand(t *testing.T) {
	output := &bytes.Buffer{}
	
	cmd := AlgorithmsCmd
	cmd.SetOut(output)
	cmd.SetErr(output)
	cmd.SetArgs([]string{})
	
	err := cmd.Execute()
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Check output contains expected algorithms
	if !bytes.Contains(output.Bytes(), []byte("ECDSA-P256-SHA256")) {
		t.Error("expected ECDSA-P256-SHA256 in output")
	}
}
```

### Key Testing Patterns

#### Pattern 1: Testing with SetArgs()
```go
cmd := SignCmd
cmd.SetArgs([]string{"-a", "ECDSA-P256-SHA256", "-k", "my-key", "-i", "data"})
err := cmd.Execute()
```

#### Pattern 2: Testing with Buffered Output
```go
output := &bytes.Buffer{}
errOutput := &bytes.Buffer{}

cmd := SignCmd
cmd.SetOut(output)
cmd.SetErr(errOutput)
cmd.SetArgs([]string{"-a", "ECDSA-P256-SHA256", "-k", "my-key", "-i", "data"})
err := cmd.Execute()

// Verify output
if !bytes.Contains(output.Bytes(), []byte("Signature")) {
	t.Error("expected Signature in output")
}
```

#### Pattern 3: Testing Flag Validation
```go
func TestMissingFlags(t *testing.T) {
	cmd := SignCmd
	cmd.SetArgs([]string{"-a", "ECDSA-P256-SHA256"})  // Missing -k and -i
	
	err := cmd.Execute()
	
	if err == nil {
		t.Fatal("expected error for missing flags")
	}
}
```

#### Pattern 4: Testing with Files
```go
func TestWithFile(t *testing.T) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	// Write test data
	tmpFile.WriteString("test content")
	tmpFile.Close()
	
	// Test with file
	cmd := SignCmd
	cmd.SetArgs([]string{"-a", "ECDSA-P256-SHA256", "-k", "key", "-i", tmpFile.Name()})
	err = cmd.Execute()
	
	// Verify result
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
```

## 3. **Integration Testing**

Test the entire CLI flow from root command:

```go
package cmd

import (
	"bytes"
	"testing"
)

func TestRootCommandWithSign(t *testing.T) {
	output := &bytes.Buffer{}
	
	cmd := rootCmd
	cmd.SetOut(output)
	cmd.SetErr(output)
	cmd.SetArgs([]string{"sign", "-a", "ECDSA-P256-SHA256", "-k", "test", "-i", "data"})
	
	err := cmd.Execute()
	
	if err != nil {
		t.Errorf("root command failed: %v", err)
	}
}

func TestRootCommandHelp(t *testing.T) {
	output := &bytes.Buffer{}
	
	cmd := rootCmd
	cmd.SetOut(output)
	cmd.SetArgs([]string{"--help"})
	
	err := cmd.Execute()
	
	if err != nil {
		t.Errorf("help command failed: %v", err)
	}
	
	if !bytes.Contains(output.Bytes(), []byte("sign")) {
		t.Error("expected 'sign' command in help")
	}
}
```

## 4. **Shell/Bash Testing** (For end-to-end tests)

Create a test script `test_cli.sh`:

```bash
#!/bin/bash

set -e

BINARY="./bin/signer"
TESTS_PASSED=0
TESTS_FAILED=0

echo "Running CLI Tests..."

# Test 1: Algorithms command
echo -n "Test: algorithms command... "
if $BINARY algorithms | grep -q "ECDSA-P256-SHA256"; then
    echo "✓ PASS"
    ((TESTS_PASSED++))
else
    echo "✗ FAIL"
    ((TESTS_FAILED++))
fi

# Test 2: Help command
echo -n "Test: help command... "
if $BINARY --help | grep -q "Sign data"; then
    echo "✓ PASS"
    ((TESTS_PASSED++))
else
    echo "✗ FAIL"
    ((TESTS_FAILED++))
fi

# Test 3: Sign missing flags
echo -n "Test: sign with missing flags... "
if ! $BINARY sign -a ECDSA-P256-SHA256 2>&1 | grep -q "required"; then
    echo "✗ FAIL"
    ((TESTS_FAILED++))
else
    echo "✓ PASS"
    ((TESTS_PASSED++))
fi

# Test 4: Sign with valid input
echo -n "Test: sign with valid input... "
OUTPUT=$($BINARY sign -a ECDSA-P256-SHA256 -k test-key -i "Hello World" 2>&1)
if echo "$OUTPUT" | grep -q "Signature"; then
    echo "✓ PASS"
    ((TESTS_PASSED++))
else
    echo "✗ FAIL"
    ((TESTS_FAILED++))
fi

echo ""
echo "Tests passed: $TESTS_PASSED"
echo "Tests failed: $TESTS_FAILED"

if [ $TESTS_FAILED -eq 0 ]; then
    exit 0
else
    exit 1
fi
```

Run with:
```bash
bash test_cli.sh
```

## 5. **Table-Driven Tests** (Best Practice)

```go
func TestSignCommandVariations(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		wantInOutput string
	}{
		{
			name:        "valid ECDSA signature",
			args:        []string{"-a", "ECDSA-P256-SHA256", "-k", "key", "-i", "data"},
			wantErr:     false,
			wantInOutput: "Signature",
		},
		{
			name:        "missing algorithm",
			args:        []string{"-k", "key", "-i", "data"},
			wantErr:     true,
			wantInOutput: "",
		},
		{
			name:        "invalid algorithm",
			args:        []string{"-a", "INVALID", "-k", "key", "-i", "data"},
			wantErr:     true,
			wantInOutput: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := SignCmd
			cmd.SetArgs(tt.args)
			
			err := cmd.Execute()
			
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %v, want %v", err != nil, tt.wantErr)
			}
		})
	}
}
```

## 6. **Run Tests**

```bash
# Run all tests
make test

# Run specific package tests
go test ./cmd/signer-command -v

# Run with coverage
go test ./cmd/signer-command -v -cover

# Run specific test
go test ./cmd/signer-command -run TestAlgorithmsCommand -v
```

## Best Practices for CLI Testing

✅ **DO:**
- Test both success and failure paths
- Use table-driven tests for multiple scenarios
- Test flag validation
- Test with temporary files for file I/O
- Capture and verify output
- Test help text

❌ **DON'T:**
- Rely on external state (use temp files instead)
- Test implementation details (test behavior)
- Skip error cases
- Ignore flag parsing

## Troubleshooting

### Issue: Flags not being set properly
**Solution:** Use `cmd.SetArgs()` before `cmd.Execute()`, not direct flag assignment

### Issue: Output not captured
**Solution:** Use `cmd.SetOut()` and `cmd.SetErr()` before executing

### Issue: Tests pass locally but fail in CI
**Solution:** Ensure no hardcoded paths; use temp files and relative paths

### Issue: Signing fails in tests
**Solution:** This is normal if keys aren't available in test environment. Focus on flag validation and error handling instead.
