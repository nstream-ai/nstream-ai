# Mothership gRPC Server Implementation Guide

This document details the complete gRPC server implementation, including service definitions, routes, request/response structures, and implementation details.

## Server Overview

The server implements three main services:
1. Auth Service - User authentication and authorization
2. Cluster Service - Cluster management
3. Bucket Service - Storage bucket management

## Service Definitions

### 1. Auth Service

#### Routes
```
/auth.AuthService/SignIn
/auth.AuthService/VerifySignIn
/auth.AuthService/SignUp
/auth.AuthService/VerifySignUp
/auth.AuthService/ValidateUser
/auth.AuthService/ValidateToken
/auth.AuthService/ValidateClusterToken
```

#### Request/Response Structures

1. **SignIn**
   ```protobuf
   // Request
   message SignInRequest {
     string email = 1;
   }
   
   // Response
   message SignInResponse {
     bool success = 1;
     string error = 2;
   }
   ```

2. **VerifySignIn**
   ```protobuf
   // Request
   message VerifySignInRequest {
     string email = 1;
     string otp = 2;
   }
   
   // Response
   message VerifySignInResponse {
     string auth_token = 1;
     google.protobuf.Timestamp expires_at = 2;
     UserInfo user_info = 3;
     string error = 4;
   }
   ```

3. **SignUp**
   ```protobuf
   // Request
   message SignUpRequest {
     string email = 1;
     string name = 2;
     string organization = 3;
     string role = 4;
   }
   
   // Response
   message SignUpResponse {
     bool success = 1;
     string error = 2;
   }
   ```

4. **VerifySignUp**
   ```protobuf
   // Request
   message VerifySignUpRequest {
     string email = 1;
     string otp = 2;
   }
   
   // Response
   message VerifySignUpResponse {
     string auth_token = 1;
     google.protobuf.Timestamp expires_at = 2;
     UserInfo user_info = 3;
     string error = 4;
   }
   ```

5. **ValidateUser**
   ```protobuf
   // Request
   message ValidateUserRequest {
     string email = 1;
   }
   
   // Response
   message ValidateUserResponse {
     bool valid = 1;
     string error = 2;
   }
   ```

6. **ValidateToken**
   ```protobuf
   // Request
   message ValidateTokenRequest {
     string token = 1;
   }
   
   // Response
   message ValidateTokenResponse {
     bool valid = 1;
     google.protobuf.Timestamp expires_at = 2;
     string error = 3;
   }
   ```

7. **ValidateClusterToken**
   ```protobuf
   // Request
   message ValidateClusterTokenRequest {
     string token = 1;
   }
   
   // Response
   message ValidateClusterTokenResponse {
     bool valid = 1;
     google.protobuf.Timestamp expires_at = 2;
     string error = 3;
   }
   ```

### 2. Cluster Service

#### Routes
```
/cluster.ClusterService/ListClusters
/cluster.ClusterService/VerifyClusterExists
/cluster.ClusterService/GetClusterDetails
/cluster.ClusterService/CreateCluster
```

#### Request/Response Structures

1. **ListClusters**
   ```protobuf
   // Request
   message ListClustersRequest {}
   
   // Response
   message ListClustersResponse {
     repeated Cluster clusters = 1;
   }
   ```

2. **VerifyClusterExists**
   ```protobuf
   // Request
   message VerifyClusterExistsRequest {
     string cluster_name = 1;
   }
   
   // Response
   message VerifyClusterExistsResponse {
     bool exists = 1;
     string error = 2;
   }
   ```

3. **GetClusterDetails**
   ```protobuf
   // Request
   message GetClusterDetailsRequest {
     string cluster_name = 1;
   }
   
   // Response
   message GetClusterDetailsResponse {
     ClusterConfig config = 1;
     string error = 2;
   }
   ```

4. **CreateCluster**
   ```protobuf
   // Request
   message CreateClusterRequest {
     string name = 1;
     string type = 2;
     string cloud_provider = 3;
     string region = 4;
     string bucket = 5;
     string role = 6;
   }
   
   // Response
   message CreateClusterResponse {
     ClusterConfig config = 1;
     string error = 2;
   }
   ```

### 3. Bucket Service

#### Routes
```
/cluster.BucketService/ListBuckets
/cluster.BucketService/VerifyBucketAccess
/cluster.BucketService/CheckResourceReadiness
```

#### Request/Response Structures

1. **ListBuckets**
   ```protobuf
   // Request
   message ListBucketsRequest {
     string cloud_provider = 1;
   }
   
   // Response
   message ListBucketsResponse {
     repeated Bucket buckets = 1;
   }
   ```

2. **VerifyBucketAccess**
   ```protobuf
   // Request
   message VerifyBucketAccessRequest {
     string cloud_provider = 1;
     string bucket = 2;
     string role = 3;
   }
   
   // Response
   message VerifyBucketAccessResponse {
     bool has_access = 1;
     string error = 2;
   }
   ```

3. **CheckResourceReadiness**
   ```protobuf
   // Request
   message CheckResourceReadinessRequest {
     string cloud_provider = 1;
     string bucket = 2;
     string role = 3;
   }
   
   // Response
   message CheckResourceReadinessResponse {
     bool ready = 1;
     string error = 2;
   }
   ```

## Common Message Types

```protobuf
message UserInfo {
  string email = 1;
  string organization = 2;
  string role = 3;
  string current_cluster = 4;
}

message Cluster {
  string id = 1;
  string region = 2;
  string cloud_provider = 3;
  string bucket = 4;
  string role = 5;
}

message ClusterConfig {
  string name = 1;
  string region = 2;
  string cloud_provider = 3;
  string bucket = 4;
  string role = 5;
  string cluster_token = 6;
}

message Bucket {
  string name = 1;
  string region = 2;
  string provider = 3;
  string size = 4;
  google.protobuf.Timestamp created_at = 5;
}
```

## Server Implementation

### Directory Structure
```
.
├── main.go                 # Server entry point
├── proto/                  # Protocol buffer definitions
│   ├── auth.proto
│   └── cluster.proto
├── pkg/
│   ├── server/            # Service implementations
│   │   ├── auth/
│   │   │   └── auth.go
│   │   └── cluster/
│   │       └── cluster.go
│   └── tls/               # TLS configuration
│       └── tls.go
└── scripts/               # Utility scripts
    ├── generate_certs.sh
    └── run.sh
```

### Server Configuration

1. **TLS Configuration**
   ```go
   tlsConfig := &tls.Config{
       CertFile:   os.Getenv("TLS_CERT_FILE"),
       KeyFile:    os.Getenv("TLS_KEY_FILE"),
       CAFile:     os.Getenv("TLS_CA_FILE"),
       ServerName: os.Getenv("TLS_SERVER_NAME"),
   }
   ```

2. **gRPC Server Setup**
   ```go
   grpcServer := grpc.NewServer(
       grpc.Creds(creds),
       grpc.UnaryInterceptor(loggingInterceptor),
   )
   ```

### Error Handling

1. **Standard Error Format**
   ```go
   type Error struct {
       Code    int32
       Message string
   }
   ```

2. **Error Codes**
   - 0: Success
   - 1: Invalid Request
   - 2: Authentication Failed
   - 3: Authorization Failed
   - 4: Resource Not Found
   - 5: Internal Server Error

### Logging

The server implements structured logging with the following information:
- Request method
- Request parameters
- Response status
- Error details (if any)
- Timestamp
- Client information

### Security

1. **mTLS Authentication**
   - Server certificate verification
   - Client certificate verification
   - CA-based trust chain

2. **Authorization**
   - Token-based access control
   - Role-based permissions
   - Cluster-specific access rights

## Running the Server

1. **Generate Certificates**
   ```bash
   ./scripts/generate_certs.sh
   ```

2. **Set Environment Variables**
   ```bash
   export TLS_CERT_FILE=/path/to/server.crt
   export TLS_KEY_FILE=/path/to/server.key
   export TLS_CA_FILE=/path/to/ca.crt
   export TLS_SERVER_NAME=mothership-server
   ```

3. **Start Server**
   ```bash
   ./scripts/run.sh
   ```

## Monitoring and Debugging

1. **Available Routes**
   ```
   🚀 Available gRPC Routes
   ==================================================
   
   📦 Auth Service
   ==============
     • SignIn                         → auth.go
     • VerifySignIn                   → auth.go
     ...
   ```

2. **Request Logging**
   ```
   🔍 Request: /auth.AuthService/SignIn
   ✅ Success: /auth.AuthService/SignIn
   ```

3. **Error Logging**
   ```
   ❌ Error: authentication failed
   ```

## Best Practices

1. **Error Handling**
   - Use appropriate error codes
   - Provide descriptive error messages
   - Log errors with context

2. **Security**
   - Validate all input
   - Implement rate limiting
   - Use secure defaults
   - Regular security audits

3. **Performance**
   - Implement connection pooling
   - Use appropriate timeouts
   - Monitor resource usage

4. **Maintenance**
   - Regular certificate rotation
   - Version control for proto files
   - Documentation updates
   - Regular testing 