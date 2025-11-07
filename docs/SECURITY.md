# Security Features

NebulaBox includes comprehensive security features for image scanning, signing, and verification.

## Image Scanning

Enhanced vulnerability scanning with an extended CVE database that checks container images for known security vulnerabilities.

### Features

- **Extended Vulnerability Database**: Includes CVEs for common packages (OpenSSL, glibc, nginx, PostgreSQL, Node.js, Python, Alpine)
- **Automatic Package Detection**: Detects installed packages based on image type
- **Severity Classification**: Categorizes vulnerabilities as CRITICAL, HIGH, MEDIUM, LOW, or UNKNOWN
- **Sorted Results**: Vulnerabilities sorted by severity (critical first)
- **Fixed Version Tracking**: Reports fixed versions when available

### API Endpoint

```bash
POST /api/images/scan
Content-Type: application/json

{
  "image": "nginx:latest"
}
```

Response:
```json
{
  "image": "nginx:latest",
  "scannedAt": "2024-01-15T10:30:00Z",
  "criticalCount": 1,
  "highCount": 2,
  "mediumCount": 1,
  "lowCount": 0,
  "unknownCount": 0,
  "vulnerabilities": [
    {
      "id": "CVE-2024-1111",
      "package": "nginx",
      "installed": "1.21.0",
      "fixedVersion": "1.25.3",
      "severity": "CRITICAL",
      "title": "nginx request smuggling vulnerability",
      "description": "Automated vulnerability detection",
      "source": "nebula-vuln-db"
    }
  ]
}
```

## Image Signing

Cryptographic signing of container images using RSA key pairs to ensure image integrity and authenticity.

### Features

- **RSA-SHA256 Signing**: Uses RSA 2048-bit keys with SHA-256 hashing
- **Key Management**: Generate, list, and manage signing keys
- **Signature Storage**: Signatures include metadata (signer, timestamp, key ID)
- **Public Key Export**: Export public keys in PEM format for distribution
- **User Attribution**: Tracks who signed each image

### API Endpoints

#### Generate Signing Key

```bash
POST /api/security/keys/generate
Content-Type: application/json

{
  "keyId": "production-key-2024"
}
```

Response:
```json
{
  "keyId": "production-key-2024",
  "createdAt": "2024-01-15T10:30:00Z",
  "createdBy": "admin",
  "publicKey": "-----BEGIN PUBLIC KEY-----\n..."
}
```

#### List Keys

```bash
GET /api/security/keys
```

Response:
```json
{
  "keys": [...],
  "count": 2
}
```

#### Sign Image

```bash
POST /api/security/sign
Content-Type: application/json

{
  "image": "myregistry/myapp",
  "tag": "v1.0.0",
  "digest": "sha256:abc123...",
  "keyId": "production-key-2024"
}
```

Response:
```json
{
  "image": "myregistry/myapp",
  "tag": "v1.0.0",
  "digest": "sha256:abc123...",
  "signedBy": "admin",
  "signedAt": "2024-01-15T10:30:00Z",
  "signature": "base64-encoded-signature",
  "publicKey": "base64-encoded-public-key",
  "algorithm": "RSA-SHA256",
  "keyId": "production-key-2024"
}
```

#### Verify Signature

```bash
POST /api/security/verify
Content-Type: application/json

{
  "image": "myregistry/myapp",
  "tag": "v1.0.0",
  "digest": "sha256:abc123...",
  "signature": {
    "image": "myregistry/myapp",
    "tag": "v1.0.0",
    "digest": "sha256:abc123...",
    "signedBy": "admin",
    "signedAt": "2024-01-15T10:30:00Z",
    "signature": "...",
    "publicKey": "...",
    "algorithm": "RSA-SHA256",
    "keyId": "production-key-2024"
  }
}
```

Response:
```json
{
  "valid": true,
  "image": "myregistry/myapp",
  "tag": "v1.0.0",
  "digest": "sha256:abc123..."
}
```

## Security Dashboard

The Security page (`/security`) provides a comprehensive interface for:

1. **Key Management**:
   - Generate new signing keys
   - List existing keys
   - View public keys (PEM format)
   - Copy keys to clipboard

2. **Image Signing**:
   - Sign images with selected keys
   - View signature details
   - Export signatures as JSON

3. **Signature Verification**:
   - Verify image signatures
   - Visual feedback (valid/invalid)
   - Error reporting

4. **Vulnerability Scanning**:
   - Integration with enhanced scanner
   - Links to Images page for scanning

## Implementation Details

### Signing Algorithm

- **Algorithm**: RSA-PKCS1v15 with SHA-256
- **Key Size**: 2048 bits
- **Hash Function**: SHA-256
- **Signature Format**: Base64-encoded PKCS#1 v1.5 signature

### Signature Payload

The signature covers:
- Image name
- Tag
- Digest (content-addressable hash)
- Signer username
- Timestamp

This ensures:
- **Integrity**: Any change to image content invalidates the signature
- **Authenticity**: Only holders of the private key can sign
- **Non-repudiation**: Signature proves who signed it
- **Timestamping**: Records when the image was signed

### Verification Process

1. Reconstruct signature payload from image metadata
2. Decode and parse public key from PEM
3. Verify RSA signature using public key
4. Validate payload matches (image, tag, digest)
5. Return verification result

## Security Best Practices

1. **Key Rotation**: Regularly rotate signing keys
2. **Key Storage**: Store private keys securely (production should use HSM/KMS)
3. **Access Control**: Limit key generation to authorized users
4. **Signature Storage**: Store signatures alongside images in registry
5. **Policy Enforcement**: Enforce signature verification before deployment
6. **Audit Logging**: Track all signing and verification operations
7. **Key Backup**: Backup private keys securely (encrypted)
8. **Public Key Distribution**: Distribute public keys to verification systems

## Integration with Registry

Signatures can be:
- Stored as metadata with image versions in the registry
- Retrieved with version information
- Validated automatically during image pull

## Future Enhancements

- Support for ECDSA keys (smaller, faster)
- Integration with Notary v2 or OCI distribution-spec signatures
- Signature expiration and revocation
- Key hierarchy (root, intermediate, signing keys)
- Hardware security module (HSM) support
- Automated scanning in CI/CD pipelines
- Security policy enforcement (block unsigned images)
- Integration with external CVE databases (NVD, CVE feeds)
- Real-time vulnerability updates
- Compliance reporting (vulnerability counts, fix rates)

## See Also

- [Registry Documentation](./REGISTRY.md)
- [CI/CD Documentation](./CI_CD.md)
- [Build Specification](./BUILD_SPEC.md)

