# Signer CLI

A command-line tool for cryptographic signing operations supporting multiple algorithms including RSA-PSS, ECDSA, and Ed25519.

## Building

To build the CLI binary:

```bash
make build-cli
```

The binary will be created at `bin/signer`.

## Installation

```bash
# Build and install to your system
go install
```

## Usage

### Help

```bash
signer --help
signer [command] --help
```

### List Supported Algorithms

```bash
signer algorithms
```

Output:
```
Supported Algorithms:
=====================
1. ECDSA-P256-SHA256
   ECDSA with P-256 curve and SHA-256 hash

2. RSA-PSS-SHA256
   RSA-PSS with SHA-256 hash

3. Ed25519
   Ed25519 elliptic curve signature scheme
```

### Sign Data

Sign data using a private key. Input can be provided as a file path or directly as a string.

```bash
# Sign from a string
signer sign --algorithm ECDSA-P256-SHA256 --key-id my-key --input "Hello World"

# Short flags
signer sign -a ECDSA-P256-SHA256 -k my-key -i "Hello World"

# Sign from a file
signer sign -a ECDSA-P256-SHA256 -k my-key -i /path/to/data.txt

# Save signature to a file
signer sign -a ECDSA-P256-SHA256 -k my-key -i "Hello World" -o signature.b64
```

**Flags:**
- `-a, --algorithm string` - Signing algorithm (required)
- `-k, --key-id string` - Key identifier (required)
- `-i, --input string` - Input data or file path (required)
- `-o, --output string` - Output file for signature (optional)

**Output:**
```
Signature (base64): MEYCIQDoUi3jGhCoDFjecwvkI68lAhgMcXOMkMjokEMC78pyHwIhALVfaf7McwHb6zbm5ytcH/JEuE2NIAKsEN9bV5ffejJE
Algorithm: ECDSA-P256-SHA256
Key ID: my-key
```

### Verify Signature

Verify that a signature is valid for the given data.

```bash
# Verify with signature as argument
signer verify -a ECDSA-P256-SHA256 -k my-key -i "Hello World" -s "MEYCIQDoUi3jGhCoDFjecwvkI68lAhgMcXOMkMjokEMC78pyHwIhALVfaf7McwHb6zbm5ytcH/JEuE2NIAKsEN9bV5ffejJE"

# Verify with signature from file
signer verify -a ECDSA-P256-SHA256 -k my-key -i "Hello World" -s signature.b64

# Verify data from file
signer verify -a ECDSA-P256-SHA256 -k my-key -i data.txt -s signature.b64
```

**Flags:**
- `-a, --algorithm string` - Signing algorithm (required)
- `-k, --key-id string` - Key identifier (required)
- `-i, --input string` - Input data or file path (required)
- `-s, --signature string` - Signature (base64) or file path (required)

**Output (success):**
```
✓ Signature verification successful
Algorithm: ECDSA-P256-SHA256
Key ID: my-key
```

**Output (failure):**
```
Verification failed: signature is invalid
```

## Examples

### Complete Workflow

```bash
# 1. See available algorithms
./bin/signer algorithms

# 2. Sign a message
SIG=$(./bin/signer sign -a ECDSA-P256-SHA256 -k my-key -i "Important message" | grep "Signature" | cut -d' ' -f3)

# 3. Verify the signature
./bin/signer verify -a ECDSA-P256-SHA256 -k my-key -i "Important message" -s "$SIG"
```

### File-based Operations

```bash
# Create a file to sign
echo "This is my document" > document.txt

# Sign it
./bin/signer sign -a RSA-PSS-SHA256 -k my-rsa-key -i document.txt -o document.sig

# Verify the signature
./bin/signer verify -a RSA-PSS-SHA256 -k my-rsa-key -i document.txt -s document.sig
```

## Key Management

Currently, the CLI uses the local signing backend which loads keys from:
- RSA keys: `backend/local/keys/rsa/signing-key.p12`
- ECDSA keys: `backend/local/keys/ecdsa/signing-key.p12`

The default password is `password1!` (configured in the code).

## Troubleshooting

### "Invalid algorithm" error
Make sure you use one of the supported algorithms:
- `ECDSA-P256-SHA256`
- `RSA-PSS-SHA256`
- `Ed25519`

### "Signing failed" error
- Check that the key files exist in the expected locations
- Verify the key ID is correct
- Ensure the input data is valid

### "Signature is invalid" error
- The signature or input data may have been modified
- Ensure you're using the same input data and algorithm for verification as was used for signing
- Verify the signature is properly formatted (base64)

## Development

### Run Tests

```bash
make test
```

### Build Everything

```bash
make build
make build-cli
```

### Clean

```bash
make clean
```

## See Also

- [Signer SDK Documentation](../README.md)
- Supported algorithms and implementations in `sdk/signers/`
