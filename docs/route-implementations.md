# NStream AI CLI Route Implementations

This document provides a comprehensive overview of how the gRPC routes defined in the server implementation are mapped to CLI commands.

## Auth Service Routes

| Route | CLI Command | File | Implementation Status | Notes |
|-------|-------------|------|----------------------|-------|
| `/auth.AuthService/SignIn` | `nsai auth signin` | `pkg/cmd/auth/signin.go` | ✅ Implemented | Handles initial sign-in request |
| `/auth.AuthService/VerifySignIn` | `nsai auth signin` | `pkg/cmd/auth/signin.go` | ✅ Implemented | Handles OTP verification |
| `/auth.AuthService/SignUp` | `nsai auth signup` | `pkg/cmd/auth/signup.go` | ✅ Implemented | Handles new user registration |
| `/auth.AuthService/VerifySignUp` | `nsai auth signup` | `pkg/cmd/auth/signup.go` | ✅ Implemented | Handles signup verification |
| `/auth.AuthService/ValidateUser` | N/A | N/A | 🔄 Implicit | Handled during signin process |
| `/auth.AuthService/ValidateToken` | N/A | N/A | 🔄 Implicit | Handled during signin process |
| `/auth.AuthService/ValidateClusterToken` | N/A | N/A | 🔄 Implicit | Handled during cluster operations |

## Cluster Service Routes

| Route | CLI Command | File | Implementation Status | Notes |
|-------|-------------|------|----------------------|-------|
| `/cluster.ClusterService/ListClusters` | `nsai create cluster` | `pkg/cmd/create/cluster.go` | ✅ Implemented | Lists available clusters |
| `/cluster.ClusterService/VerifyClusterExists` | N/A | N/A | 🔄 Implicit | Used internally by other commands |
| `/cluster.ClusterService/GetClusterDetails` | `nsai create cluster` | `pkg/cmd/create/cluster.go` | ✅ Implemented | Gets detailed cluster information |
| `/cluster.ClusterService/CreateCluster` | `nsai create cluster` | `pkg/cmd/create/cluster.go` | ✅ Implemented | Creates new cluster |

## Bucket Service Routes

| Route | CLI Command | File | Implementation Status | Notes |
|-------|-------------|------|----------------------|-------|
| `/cluster.BucketService/ListBuckets` | `nsai create bucket` | `pkg/cmd/create/bucket.go` | ✅ Implemented | Lists available buckets |
| `/cluster.BucketService/VerifyBucketAccess` | `nsai create cluster` | `pkg/cmd/create/cluster.go` | ✅ Implemented | Used in cluster creation workflow |
| `/cluster.BucketService/CheckResourceReadiness` | `nsai create cluster` | `pkg/cmd/create/cluster.go` | ✅ Implemented | Used in cluster creation workflow |

## Implementation Status Legend

- ✅ Implemented: Route is fully implemented as a CLI command
- 🔄 Implicit: Route is used internally but not exposed as a direct command
- ❌ Not Implemented: Route is not currently implemented

## Notes

1. Some routes are used implicitly by other commands rather than being exposed as direct CLI commands. This is by design to provide a more streamlined user experience.

2. The bucket creation functionality currently uses the `CreateCluster` route, as there isn't a dedicated bucket creation route in the server implementation.

3. Validation routes are typically handled during the signin process or other operations rather than being exposed as separate commands.

4. The implementation follows a hierarchical command structure:
   - `nsai auth [signin|signup]` for authentication
   - `nsai create [cluster|bucket]` for resource creation
   - `nsai get [cluster|bucket]` for resource information

5. Bucket verification and resource readiness checks are now implemented as part of the cluster creation workflow in `pkg/cmd/create/cluster.go`.

## Future Improvements

1. Consider implementing dedicated commands for validation routes if direct access is needed
2. Add a dedicated bucket creation route in the server implementation
3. Consider exposing bucket access verification and resource readiness checks as standalone commands for troubleshooting purposes 