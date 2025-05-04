package create

import (
	"fmt"
	"time"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
)

const (
	redColor     = "\033[31m"
	boldColor    = "\033[1m"
	resetColor   = "\033[0m"
	blinkColor   = "\033[5m"
	reverseColor = "\033[7m"
)

// ShowLoading displays a colorful and interactive loading spinner with the given message
func ShowLoading(message string, done chan bool) {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	colors := []string{redColor, boldColor, blinkColor, reverseColor}
	i := 0
	j := 0
	lastTime := time.Now()
	elapsed := 0

	for {
		select {
		case <-done:
			fmt.Printf("\r%s%s%s %s✓%s\n", boldColor, redColor, message, blinkColor, resetColor)
			return
		default:
			now := time.Now()
			elapsed += int(now.Sub(lastTime).Milliseconds())
			lastTime = now

			if elapsed >= 500 {
				j = (j + 1) % len(colors)
				elapsed = 0
			}

			dynamicMessage := fmt.Sprintf("%s%s%s%s", colors[j], boldColor, message, resetColor)

			extra := ""
			if elapsed < 100 {
				extra = " 🔥"
			} else if elapsed < 200 {
				extra = " ⚡"
			} else if elapsed < 300 {
				extra = " 💫"
			} else if elapsed < 400 {
				extra = " ✨"
			}

			fmt.Printf("\r%s %s%s", dynamicMessage, spinner[i], extra)
			time.Sleep(100 * time.Millisecond)
			i = (i + 1) % len(spinner)
		}
	}
}

// DummyCheckCredits simulates checking user credits
func DummyCheckCredits() (bool, error) {
	time.Sleep(1 * time.Second)
	return true, nil // For now, assume user has credits
}

// DummyCreateCluster simulates creating a cluster
func DummyCreateCluster(name, clusterType, cloudProvider, region, bucket, role string) (*config.ClusterConfig, error) {
	// Simulate cluster creation process
	time.Sleep(5 * time.Second)

	return &config.ClusterConfig{
		Name:          name,
		CloudProvider: cloudProvider,
		Region:        region,
		Bucket:        bucket,
		ClusterToken:  "dummy-cluster-token-1234567890",
	}, nil
}

// GetCloudRegions returns available regions for a cloud provider
func GetCloudRegions(provider string) []string {
	switch provider {
	case "aws":
		return []string{"us-east-1", "us-west-2", "eu-west-1", "ap-southeast-1"}
	case "gcp":
		return []string{"us-central1", "us-east1", "europe-west1", "asia-east1"}
	case "azure":
		return []string{"eastus", "westus", "northeurope", "southeastasia"}
	default:
		return []string{}
	}
}

// DummyGetServiceRole simulates getting the service role from gRPC service
func DummyGetServiceRole(provider string) (string, error) {
	time.Sleep(1 * time.Second)
	switch provider {
	case "aws":
		return "arn:aws:iam::123456789012:role/nstream-service-role", nil
	case "gcp":
		return "nstream-service@nstream-ai.iam.gserviceaccount.com", nil
	case "azure":
		return "12345678-1234-1234-1234-123456789012", nil
	default:
		return "", fmt.Errorf("unsupported cloud provider: %s", provider)
	}
}

// GetCloudSetupInstructions returns setup instructions for the cloud provider
func GetCloudSetupInstructions(provider, serviceRole, userRole string) string {
	switch provider {
	case "aws":
		return fmt.Sprintf(`AWS S3 Bucket Access Setup:
1. Create an IAM role named '%s' for S3 bucket access with the following trust policy:
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "AWS": "%s"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}

2. Add the following S3 bucket permissions to the role:
- s3:ListBucket
- s3:GetObject
- s3:PutObject
- s3:DeleteObject

3. Create a trust relationship between your role and the NStream service role for S3 access:
aws iam update-assume-role-policy --role-name %s --policy-document '{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "s3.amazonaws.com",
                "AWS": "%s"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}'`, userRole, serviceRole, userRole, serviceRole)
	case "gcp":
		return fmt.Sprintf(`GCP Cloud Storage Bucket Access Setup:
1. Create a service account named '%s' for bucket access with the following roles:
- Storage Object Viewer
- Storage Object Creator
- Storage Object Admin

2. Create and download a key for your service account:
gcloud iam service-accounts keys create key.json --iam-account=%s

3. Grant the NStream service account access to your service account for bucket operations:
gcloud iam service-accounts add-iam-policy-binding %s \
    --member="serviceAccount:%s" \
    --role="roles/iam.serviceAccountUser"

4. Grant both service accounts access to your bucket:
gsutil iam ch serviceAccount:%s:objectViewer,objectCreator gs://YOUR_BUCKET_NAME
gsutil iam ch serviceAccount:%s:objectViewer,objectCreator gs://YOUR_BUCKET_NAME`, userRole, userRole, userRole, serviceRole, userRole, serviceRole)
	case "azure":
		return fmt.Sprintf(`Azure Blob Storage Access Setup:
1. Create a service principal named '%s' for blob storage access:
az ad sp create-for-rbac --name %s

2. Assign the Storage Blob Data Contributor role to your service principal:
az role assignment create --assignee %s \
    --role "Storage Blob Data Contributor" \
    --scope "/subscriptions/YOUR_SUBSCRIPTION_ID/resourceGroups/YOUR_RESOURCE_GROUP/providers/Microsoft.Storage/storageAccounts/YOUR_STORAGE_ACCOUNT"

3. Grant the NStream service principal access to your blob storage:
az role assignment create --assignee %s \
    --role "Storage Blob Data Contributor" \
    --scope "/subscriptions/YOUR_SUBSCRIPTION_ID/resourceGroups/YOUR_RESOURCE_GROUP/providers/Microsoft.Storage/storageAccounts/YOUR_STORAGE_ACCOUNT"

4. Add the NStream service principal to your app registration for blob access:
az ad app credential list --id %s
az ad app credential add --id %s --cert "@/path/to/certificate.pem"`, userRole, userRole, userRole, serviceRole, userRole, serviceRole)
	default:
		return "No setup instructions available for this cloud provider"
	}
}

// VerifyBucketAccess verifies if the bucket is accessible with the provided credentials
func VerifyBucketAccess(provider, bucket, role string) error {
	// Simulate bucket access verification
	time.Sleep(3 * time.Second)
	return nil
}

// CheckResourceReadiness checks if all required resources are ready
func CheckResourceReadiness(provider, bucket, role string) error {
	// Check bucket attachment readiness
	if err := checkBucketAttachmentReady(provider, bucket, role); err != nil {
		return fmt.Errorf("bucket attachment not ready: %v", err)
	}

	// Check base models availability
	if err := checkBaseModelsReady(provider); err != nil {
		return fmt.Errorf("base models not ready: %v", err)
	}

	// Check knowledgebases readiness
	if err := checkKnowledgebasesReady(provider); err != nil {
		return fmt.Errorf("knowledgebases not ready: %v", err)
	}

	// Check customer resources readiness
	if err := checkCustomerResourcesReady(provider); err != nil {
		return fmt.Errorf("customer resources not ready: %v", err)
	}

	return nil
}

// Helper functions for resource readiness checks
func checkBucketAttachmentReady(provider, bucket, role string) error {
	// Simulate checking bucket attachment
	time.Sleep(2 * time.Second)
	return nil
}

func checkBaseModelsReady(provider string) error {
	// Simulate checking base models
	time.Sleep(2 * time.Second)
	return nil
}

func checkKnowledgebasesReady(provider string) error {
	// Simulate checking knowledgebases
	time.Sleep(2 * time.Second)
	return nil
}

func checkCustomerResourcesReady(provider string) error {
	// Simulate checking customer resources
	time.Sleep(2 * time.Second)
	return nil
}

// VerifyUserAuth checks if the user's authentication token is valid
func VerifyUserAuth() error {
	// Load config to get auth token
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	if cfg.User.AuthToken == "" {
		return fmt.Errorf("no authentication token found. Please login first")
	}

	// Simulate token validation
	time.Sleep(1 * time.Second)
	return nil
}

// VerifyUserExists checks if the user exists in the system
func VerifyUserExists() error {
	// Load config to get user details
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	if cfg.User.Email == "" {
		return fmt.Errorf("no user found. Please login first")
	}

	// Simulate user existence check
	time.Sleep(1 * time.Second)
	return nil
}
