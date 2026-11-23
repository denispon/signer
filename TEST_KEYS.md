# Test Keys for Signer CLI

## Available Test Keys

The signer project includes pre-generated test keys in PKCS#12 format (.p12):

### **Key Locations**

```
backend/local/keys/
├── ecdsa/
│   └── signing-key.p12    (ECDSA P-256 key)
└── rsa/
    └── signing-key.p12    (RSA key)
```

### **Key Details**

| Key Type | Algorithm | File | Password | Format |
|----------|-----------|------|----------|--------|
| ECDSA P-256 | ECDSA-P256-SHA256 | `backend/local/keys/ecdsa/signing-key.p12` | `password1!` | PKCS#12 |
| RSA | RSA-PSS-SHA256 | `backend/local/keys/rsa/signing-key.p12` | `password1!` | PKCS#12 |

### **Key Password**

All test keys are protected with the password: **`password1!`**

This password is hardcoded in `backend/local/localsigner.go` (line 35):
```go
privKey, _, err := loadP12Key(path+"/"+"signing-key.p12", "password1!")
```

---

## Testing the CLI with These Keys

### **Test 1: Sign with ECDSA Key**

```bash
./bin/signer sign -a ECDSA-P256-SHA256 -k ecdsa-test-key -i "Hello World"
```

**Expected Output:**
```
Signature (base64): MEUCIQCoQ3vL2ysmoMuFMrJSfzy17mQO6oR/8rz7GsXzVZGS3AIgBU04TehXJsjAXTm2ZWly/A6FSeWtd+O8yrhcyIpQvsU=
Algorithm: ECDSA-P256-SHA256
Key ID: ecdsa-test-key
```

### **Test 2: Sign with RSA Key**

```bash
./bin/signer sign -a RSA-PSS-SHA256 -k rsa-test-key -i "Hello World"
```

**Expected Output:**
```
Signature (base64): [Base64 encoded RSA signature]
Algorithm: RSA-PSS-SHA256
Key ID: rsa-test-key
```

### **Test 3: Verify ECDSA Signature**

```bash
SIG=$(./bin/signer sign -a ECDSA-P256-SHA256 -k ecdsa-test-key -i "Hello World" | grep "Signature" | cut -d' ' -f3)

./bin/signer verify -a ECDSA-P256-SHA256 -k ecdsa-test-key -i "Hello World" -s "$SIG"
```

**Expected Output:**
```
✓ Signature verification successful
Algorithm: ECDSA-P256-SHA256
Key ID: ecdsa-test-key
```

### **Test 4: Verify RSA Signature**

```bash
SIG=$(./bin/signer sign -a RSA-PSS-SHA256 -k rsa-test-key -i "Hello World" | grep "Signature" | cut -d' ' -f3)

./bin/signer verify -a RSA-PSS-SHA256 -k rsa-test-key -i "Hello World" -s "$SIG"
```

---

## Understanding the Key ID Parameter

The `-k` or `--key-id` parameter is **not** used to select which key file to load. Instead:

1. The **algorithm** (`-a`) determines which key file is loaded
2. The **key ID** is just metadata included in the response

### How It Works:

```
-a ECDSA-P256-SHA256  → Loads backend/local/keys/ecdsa/signing-key.p12
-a RSA-PSS-SHA256     → Loads backend/local/keys/rsa/signing-key.p12
-a Ed25519            → (Not yet implemented in backend)

-k my-key-id          → Just metadata, included in output
```

**Example:**
```bash
./bin/signer sign -a ECDSA-P256-SHA256 -k "my-custom-id" -i "data"

# Output includes:
# Key ID: my-custom-id
```

---

## Full Testing Workflow

### **1. Build the Binary**
```bash
make build-cli
```

### **2. Test All Algorithms**
```bash
# Test ECDSA
echo "=== Testing ECDSA ==="
./bin/signer sign -a ECDSA-P256-SHA256 -k ecdsa-test -i "test message"

# Test RSA
echo "=== Testing RSA ==="
./bin/signer sign -a RSA-PSS-SHA256 -k rsa-test -i "test message"

# Test Ed25519 (should fail - not implemented in backend yet)
echo "=== Testing Ed25519 ==="
./bin/signer sign -a Ed25519 -k ed25519-test -i "test message"
```

### **3. Test with Files**
```bash
# Create test file
echo "File content to sign" > /tmp/test.txt

# Sign it
./bin/signer sign -a ECDSA-P256-SHA256 -k file-key -i /tmp/test.txt -o /tmp/signature.b64

# Verify
./bin/signer verify -a ECDSA-P256-SHA256 -k file-key -i /tmp/test.txt -s /tmp/signature.b64
```

### **4. Test Error Cases**
```bash
# Missing algorithm
./bin/signer sign -k key -i "data"
# Expected: Error - --algorithm is required

# Invalid algorithm
./bin/signer sign -a INVALID -k key -i "data"
# Expected: Error - invalid algorithm

# Missing input
./bin/signer sign -a ECDSA-P256-SHA256 -k key
# Expected: Error - --input is required

# Wrong data for verification
./bin/signer verify -a ECDSA-P256-SHA256 -k key -i "wrong data" -s "signature_here"
# Expected: Verification failed
```

---

## Test Key Generation (Optional)

If you need to regenerate the test keys:

### **Generate ECDSA Key**
```bash
# Generate EC private key
openssl ecparam -name prime256v1 -genkey -noout -out ec_key.pem

# Create PKCS#12 file
openssl pkcs12 -export -in ec_key.pem -out signing-key.p12 \
  -passout pass:password1! -nokeys
```

### **Generate RSA Key**
```bash
# Generate RSA private key
openssl genrsa -out rsa_key.pem 2048

# Create PKCS#12 file
openssl pkcs12 -export -in rsa_key.pem -out signing-key.p12 \
  -passout pass:password1! -nokeys
```

---

## Troubleshooting

### **Issue: "signing failed" error**
- Ensure you're in the correct directory (project root)
- Check that key files exist: `ls backend/local/keys/*/signing-key.p12`
- Verify file permissions are readable

### **Issue: "invalid algorithm" error**
- Use exact algorithm names (case-sensitive):
  - `ECDSA-P256-SHA256` ✅
  - `ecdsa-p256-sha256` ❌
  - `RSA-PSS-SHA256` ✅
  - `rsa-pss-sha256` ❌

### **Issue: Verification fails even with correct data**
- Verify you're using the same input data
- Check the signature is base64 encoded
- Ensure algorithm matches

### **Issue: Ed25519 not working**
- Ed25519 is defined but not fully implemented in `backend/local/`
- Currently only ECDSA-P256-SHA256 and RSA-PSS-SHA256 work

---

## Unit Test Key Usage

In unit tests, you can:

1. **Use these keys** - If running tests from project root
2. **Mock keys** - Create in-memory keys for testing
3. **Skip signing tests** - Focus on CLI flag validation

### Example: Testing Flag Validation Only
```go
func TestSignFlagValidation(t *testing.T) {
	cmd := SignCmd
	cmd.SetArgs([]string{"-a", "INVALID", "-k", "key", "-i", "data"})
	
	err := cmd.Execute()
	
	// Error expected (invalid algorithm)
	if err == nil {
		t.Fatal("expected error for invalid algorithm")
	}
}
```

This tests the CLI logic without needing actual keys!

---

## Summary

✅ **Available Test Keys:**
- ECDSA: `backend/local/keys/ecdsa/signing-key.p12`
- RSA: `backend/local/keys/rsa/signing-key.p12`
- Password: `password1!`

✅ **Algorithm Selection:**
- Use `-a ECDSA-P256-SHA256` for ECDSA key
- Use `-a RSA-PSS-SHA256` for RSA key
- Key ID (`-k`) is just metadata, not for key selection

✅ **Testing Strategy:**
- For integration tests: Use real keys + real data
- For unit tests: Mock or skip key-dependent tests
- For CLI validation: Test flags, not cryptography
