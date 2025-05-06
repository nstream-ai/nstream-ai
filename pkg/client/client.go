package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	authproto "github.com/nstreama-ai/nstream-ai-cli/proto/auth"
	clusterproto "github.com/nstreama-ai/nstream-ai-cli/proto/cluster"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Client represents the mothership gRPC client
type Client struct {
	conn *grpc.ClientConn

	// Service clients
	AuthClient    authproto.AuthServiceClient
	ClusterClient clusterproto.ClusterServiceClient
	BucketClient  clusterproto.BucketServiceClient
	// InitClient            authproto.InitServiceClient
	// BaseModelClient       authproto.BaseModelServiceClient
	// MegaModelClient       authproto.MegaModelServiceClient
	// EmbeddingModelClient  authproto.EmbeddingModelServiceClient
	// StreamFinetunerClient authproto.StreamFineTunerServiceClient
	// StreamGraphClient     authproto.StreamGraphServiceClient
	// StreamConnectorClient authproto.StreamConnectorServiceClient
	// KnowledgeBaseClient   authproto.KnowledgeBaseServiceClient
}

// NewClient creates a new mothership client
func NewClient(serverAddr string, useTLS bool, caCertPath string) (*Client, error) {
	var opts []grpc.DialOption

	// Set default server address if none provided
	if serverAddr == "" {
		serverAddr = "localhost:8080"
	}

	if useTLS {
		// Create certificate pool
		certPool := x509.NewCertPool()

		// Load CA certificate from certs directory
		caCert, err := os.ReadFile("certs/ca.crt")
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %v", err)
		}
		if !certPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA certificate")
		}

		// Load client certificate and key
		clientCert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %v", err)
		}

		// Extract hostname from server address
		hostname := serverAddr
		if strings.Contains(hostname, ":") {
			hostname = strings.Split(hostname, ":")[0]
		}

		// Create TLS credentials with mTLS
		creds := credentials.NewTLS(&tls.Config{
			RootCAs:      certPool,
			Certificates: []tls.Certificate{clientCert},
			ServerName:   hostname, // Use the actual server hostname
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Connect to the server
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	// Create service clients
	client := &Client{
		conn:          conn,
		AuthClient:    authproto.NewAuthServiceClient(conn),
		ClusterClient: clusterproto.NewClusterServiceClient(conn),
		BucketClient:  clusterproto.NewBucketServiceClient(conn),
		// BaseModelClient:       proto.NewBaseModelServiceClient(conn),
		// MegaModelClient:       proto.NewMegaModelServiceClient(conn),
		// EmbeddingModelClient:  proto.NewEmbeddingModelServiceClient(conn),
		// StreamFinetunerClient: proto.NewStreamFineTunerServiceClient(conn),
		// StreamGraphClient:     proto.NewStreamGraphServiceClient(conn),
		// StreamConnectorClient: proto.NewStreamConnectorServiceClient(conn),
		// KnowledgeBaseClient:   proto.NewKnowledgeBaseServiceClient(conn),
	}

	return client, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// WithContext returns a context with timeout
func (c *Client) WithContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 2*time.Minute)
}
