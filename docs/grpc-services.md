# NStream AI CLI gRPC Services Documentation

This document outlines all gRPC service calls used in the NStream AI CLI, including their request/response structures and expected behaviors.

## Authentication Services

### 1. SignIn
Initiates the sign-in process by sending a one-time password to the user's email.

**Request:**
```protobuf
message SignInRequest {
    string email = 1;
}
```

**Response:**
```protobuf
message SignInResponse {
    bool success = 1;
    string error = 2;
}
```

**Expected Behavior:**
- Sends a one-time password to the provided email
- Returns success if email exists in the system
- Returns error if:
  - Email doesn't exist
  - Account is locked
  - Rate limit exceeded
- Server should respond within 1s

### 2. VerifySignIn
Verifies the one-time password sent to user's email.

**Request:**
```protobuf
message VerifySignInRequest {
    string email = 1;
    string otp = 2;
}
```

**Response:**
```protobuf
message VerifySignInResponse {
    string auth_token = 1;
    google.protobuf.Timestamp expires_at = 2;
    UserInfo user_info = 3;
    string error = 4;
}

message UserInfo {
    string email = 1;
    string organization = 2;
    string role = 3;
    string current_cluster = 4;
}
```

**Expected Behavior:**
- Verifies the OTP sent to user's email
- Returns auth token and user info on success
- Returns error if:
  - OTP is invalid or expired
  - Rate limit exceeded
- Server should respond within 1s
- Auth token should be valid for 24 hours

### 3. SignUp
Initiates the sign-up process by sending a one-time password to the user's email.

**Request:**
```protobuf
message SignUpRequest {
    string email = 1;
    string name = 2;
    string organization = 3;
    string role = 4;
}
```

**Response:**
```protobuf
message SignUpResponse {
    bool success = 1;
    string error = 2;
}
```

**Expected Behavior:**
- Sends a one-time password to the provided email
- Returns success if email is valid and not already registered
- Returns error if:
  - Email is already registered
  - Email is invalid
  - Organization is invalid
  - Role is invalid
  - Name is invalid
  - Rate limit exceeded
- Server should respond within 1s

### 4. VerifySignUp
Verifies the one-time password and completes registration.

**Request:**
```protobuf
message VerifySignUpRequest {
    string email = 1;
    string otp = 2;
}
```

**Response:**
```protobuf
message VerifySignUpResponse {
    string auth_token = 1;
    google.protobuf.Timestamp expires_at = 2;
    UserInfo user_info = 3;
    string error = 4;
}
```

**Expected Behavior:**
- Verifies OTP and creates user account
- Returns auth token and user info on success
- Returns error if:
  - OTP is invalid or expired
  - Rate limit exceeded
- Server should respond within 1s
- Auth token should be valid for 24 hours

### 5. ValidateUser
Validates if a user exists in the system.

**Request:**
```protobuf
message ValidateUserRequest {
    string email = 1;
}
```

**Response:**
```protobuf
message ValidateUserResponse {
    bool valid = 1;
    string error = 2;
}
```

**Expected Behavior:**
- Returns `valid: true` if user exists
- Returns `valid: false` with error message if user doesn't exist
- Returns error if email is empty or invalid
- Server should respond within 500ms

### 6. ValidateToken
Validates an authentication token.

**Request:**
```protobuf
message ValidateTokenRequest {
    string token = 1;
}
```

**Response:**
```protobuf
message ValidateTokenResponse {
    bool valid = 1;
    google.protobuf.Timestamp expires_at = 2;
    string error = 3;
}
```

**Expected Behavior:**
- Returns `valid: true` and expiration time if token is valid
- Returns `valid: false` with error message if token is:
  - Empty
  - Expired
  - Invalid
  - Malformed
- Server should respond within 500ms

### 7. ValidateClusterToken
Validates a cluster-specific token.

**Request:**
```protobuf
message ValidateClusterTokenRequest {
    string token = 1;
}
```

**Response:**
```protobuf
message ValidateClusterTokenResponse {
    bool valid = 1;
    google.protobuf.Timestamp expires_at = 2;
    string error = 3;
}
```

**Expected Behavior:**
- Returns `valid: true` and expiration time if token is valid
- Returns `valid: false` with error message if token is:
  - Empty
  - Expired
  - Invalid
  - Malformed
- Server should respond within 500ms

## Security Requirements for Authentication

1. **OTP Requirements**:
   - 6-digit numeric code
   - Valid for 10 minutes
   - One-time use only
   - Rate limited to 3 attempts
   - Different OTPs for signin and signup

2. **Rate Limiting**:
   - SignIn: 3 attempts per 15 minutes
   - SignUp: 3 attempts per hour
   - VerifySignIn: 3 attempts per OTP
   - VerifySignUp: 3 attempts per OTP

3. **Token Security**:
   - JWT format
   - Signed with HS256
   - Contains user ID, email, and role
   - 24-hour expiration
   - Refresh token mechanism available

4. **Email Verification**:
   - Required for both signin and signup
   - OTP sent via email
   - HTML and plain text formats
   - Includes organization name
   - Includes security warnings
   - Clear distinction between signin and signup emails

## Cluster Services

### 1. ListClusters
Retrieves a list of available clusters.

**Request:**
```protobuf
message ListClustersRequest {
    // Empty request
}
```

**Response:**
```protobuf
message ListClustersResponse {
    repeated Cluster clusters = 1;
}

message Cluster {
    string id = 1;
    string region = 2;
    string cloud_provider = 3;
    string bucket = 4;
    string role = 5;
}
```

**Expected Behavior:**
- Returns list of all clusters user has access to
- Empty list if no clusters exist
- Server should respond within 1s

### 2. VerifyClusterExists
Checks if a specific cluster exists.

**Request:**
```protobuf
message VerifyClusterExistsRequest {
    string cluster_name = 1;
}
```

**Response:**
```protobuf
message VerifyClusterExistsResponse {
    bool exists = 1;
    string error = 2;
}
```

**Expected Behavior:**
- Returns `exists: true` if cluster exists
- Returns `exists: false` if cluster doesn't exist
- Returns error if cluster name is invalid
- Server should respond within 1s

### 3. GetClusterDetails
Retrieves detailed information about a specific cluster.

**Request:**
```protobuf
message GetClusterDetailsRequest {
    string cluster_name = 1;
}
```

**Response:**
```protobuf
message GetClusterDetailsResponse {
    ClusterConfig config = 1;
    string error = 2;
}

message ClusterConfig {
    string name = 1;
    string region = 2;
    string cloud_provider = 3;
    string bucket = 4;
    string role = 5;
    string cluster_token = 6;
}
```

**Expected Behavior:**
- Returns complete cluster configuration if cluster exists
- Returns error if cluster doesn't exist
- Server should respond within 1s

### 4. CreateCluster
Creates a new cluster.

**Request:**
```protobuf
message CreateClusterRequest {
    string name = 1;
    string type = 2;
    string cloud_provider = 3;
    string region = 4;
    string bucket = 5;
    string role = 6;
}
```

**Response:**
```protobuf
message CreateClusterResponse {
    ClusterConfig config = 1;
    string error = 2;
}
```

**Expected Behavior:**
- Creates new cluster with specified configuration
- Returns cluster configuration on success
- Returns error if:
  - Cluster name already exists
  - Invalid configuration
  - Insufficient permissions
  - Resource creation fails
- Server should respond within 2s

## Bucket Services

### 1. ListBuckets
Retrieves a list of available buckets.

**Request:**
```protobuf
message ListBucketsRequest {
    string cloud_provider = 1;
}
```

**Response:**
```protobuf
message ListBucketsResponse {
    repeated Bucket buckets = 1;
}

message Bucket {
    string name = 1;
    string region = 2;
    string provider = 3;
    string size = 4;
    google.protobuf.Timestamp created_at = 5;
}
```

**Expected Behavior:**
- Returns list of all buckets user has access to
- Filters by cloud provider if specified
- Empty list if no buckets exist
- Server should respond within 1s

### 2. VerifyBucketAccess
Verifies access to a specific bucket.

**Request:**
```protobuf
message VerifyBucketAccessRequest {
    string cloud_provider = 1;
    string bucket = 2;
    string role = 3;
}
```

**Response:**
```protobuf
message VerifyBucketAccessResponse {
    bool has_access = 1;
    string error = 2;
}
```

**Expected Behavior:**
- Returns `has_access: true` if role has proper permissions
- Returns `has_access: false` with error if:
  - Bucket doesn't exist
  - Role doesn't have required permissions
  - Invalid configuration
- Server should respond within 1s

### 3. CheckResourceReadiness
Checks if all required resources are ready.

**Request:**
```protobuf
message CheckResourceReadinessRequest {
    string cloud_provider = 1;
    string bucket = 2;
    string role = 3;
}
```

**Response:**
```protobuf
message CheckResourceReadinessResponse {
    bool ready = 1;
    string error = 2;
}
```

**Expected Behavior:**
- Returns `ready: true` if all resources are properly configured
- Returns `ready: false` with error if:
  - Resources are not ready
  - Configuration is incomplete
  - Permissions are insufficient
- Server should respond within 1s

## Error Handling

All gRPC services should follow these error handling guidelines:

1. **Timeout**: All services should respond within their specified timeouts
2. **Error Codes**:
   - `INVALID_ARGUMENT`: Invalid input parameters
   - `NOT_FOUND`: Resource doesn't exist
   - `PERMISSION_DENIED`: Insufficient permissions
   - `UNAUTHENTICATED`: Invalid or expired tokens
   - `RESOURCE_EXHAUSTED`: Rate limiting or quota exceeded
   - `INTERNAL`: Server-side errors

3. **Error Messages**: Should be clear and actionable
4. **Retry Logic**: Client should implement exponential backoff for retries

## Rate Limiting

1. **Authentication Services**: 100 requests per minute
2. **Cluster Services**: 50 requests per minute
3. **Bucket Services**: 50 requests per minute

## Security

1. All services require valid authentication token
2. Tokens should be validated on every request
3. Sensitive data should be encrypted in transit
4. Access control should be enforced at service level 