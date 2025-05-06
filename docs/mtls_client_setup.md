# mTLS Client Setup Guide

This guide explains how to update your gRPC clients to use mutual TLS (mTLS) authentication with the Mothership server.

## Prerequisites

1. Client certificates (obtain from server administrator):
   - `client.crt` - Client certificate
   - `client.key` - Client private key
   - `ca.crt` - CA certificate

2. Required Go packages:
   ```go
   import (
       "google.golang.org/grpc"
       "google.golang.org/grpc/credentials"
   )
   ```

## Client Setup

### 1. Load TLS Configuration

```go
import "github.com/nstream-ai/nstream-ai-mothership/pkg/tls"

// Create TLS config
tlsConfig := &tls.Config{
    CertFile:   "/path/to/client.crt",
    KeyFile:    "/path/to/client.key",
    CAFile:     "/path/to/ca.crt",
    ServerName: "mothership-server",
}

// Load client TLS configuration
clientTLSConfig, err := tls.LoadClientTLSConfig(tlsConfig)
if err != nil {
    log.Fatalf("Failed to load TLS configuration: %v", err)
}

// Create credentials
creds := credentials.NewTLS(clientTLSConfig)
```

### 2. Create gRPC Connection

```go
// Create connection with TLS credentials
conn, err := grpc.Dial(
    "localhost:8080",
    grpc.WithTransportCredentials(creds),
)
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}
defer conn.Close()
```

### 3. Create Client

```go
// Create client using the secure connection
client := pb.NewYourServiceClient(conn)
```

## Complete Example

Here's a complete example for an auth client:

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/nstream-ai/nstream-ai-mothership/pkg/tls"
    "github.com/nstream-ai/nstream-ai-mothership/proto/auth"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

func main() {
    // Load TLS configuration
    tlsConfig := &tls.Config{
        CertFile:   os.Getenv("TLS_CERT_FILE"),
        KeyFile:    os.Getenv("TLS_KEY_FILE"),
        CAFile:     os.Getenv("TLS_CA_FILE"),
        ServerName: os.Getenv("TLS_SERVER_NAME"),
    }

    clientTLSConfig, err := tls.LoadClientTLSConfig(tlsConfig)
    if err != nil {
        log.Fatalf("Failed to load TLS configuration: %v", err)
    }

    // Create connection with TLS credentials
    conn, err := grpc.Dial(
        "localhost:8080",
        grpc.WithTransportCredentials(credentials.NewTLS(clientTLSConfig)),
    )
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    // Create auth client
    client := auth.NewAuthServiceClient(conn)

    // Use the client
    resp, err := client.SignIn(context.Background(), &auth.SignInRequest{
        Email: "user@example.com",
    })
    if err != nil {
        log.Fatalf("Failed to sign in: %v", err)
    }
    log.Printf("Sign in response: %v", resp)
}
```

## Environment Variables

Set these environment variables in your client environment:

```bash
export TLS_CERT_FILE=/path/to/client.crt
export TLS_KEY_FILE=/path/to/client.key
export TLS_CA_FILE=/path/to/ca.crt
export TLS_SERVER_NAME=mothership-server
```

## Troubleshooting

1. **Certificate Errors**
   - Ensure all certificate files are readable
   - Verify certificate paths are correct
   - Check certificate permissions (600 for keys, 644 for certs)

2. **Connection Errors**
   - Verify server is running and accessible
   - Check server name matches certificate
   - Ensure client certificate is signed by the CA

3. **Common Issues**
   - "certificate signed by unknown authority" - CA certificate not properly loaded
   - "bad certificate" - Client certificate not properly signed
   - "connection refused" - Server not running or wrong address

## Security Notes

1. Keep private keys secure:
   - Never commit keys to version control
   - Use secure key storage
   - Rotate certificates regularly

2. Certificate Management:
   - Monitor certificate expiration
   - Implement certificate rotation
   - Keep CA private key secure

3. Best Practices:
   - Use strong key sizes (4096 bits)
   - Implement proper error handling
   - Log security events
   - Regular security audits 