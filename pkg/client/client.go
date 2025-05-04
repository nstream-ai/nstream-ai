package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"github.com/nstreama-ai/nstream-ai-mothership/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Client represents the mothership gRPC client
type Client struct {
	conn *grpc.ClientConn

	// Service clients
	AuthClient            proto.AuthServiceClient
	InitClient            proto.InitServiceClient
	BaseModelClient       proto.BaseModelServiceClient
	MegaModelClient       proto.MegaModelServiceClient
	EmbeddingModelClient  proto.EmbeddingModelServiceClient
	StreamFinetunerClient proto.StreamFineTunerServiceClient
	StreamGraphClient     proto.StreamGraphServiceClient
	StreamConnectorClient proto.StreamConnectorServiceClient
	KnowledgeBaseClient   proto.KnowledgeBaseServiceClient
}

// NewClient creates a new mothership client
func NewClient(serverAddr string, useTLS bool, caCertPath string) (*Client, error) {
	var opts []grpc.DialOption

	if useTLS {
		// Load CA certificate
		caCert, err := os.ReadFile(caCertPath)
		if err != nil {
			return nil, err
		}

		// Create certificate pool
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caCert) {
			return nil, err
		}

		// Create TLS credentials
		creds := credentials.NewTLS(&tls.Config{
			RootCAs: certPool,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Connect to the server
	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		return nil, err
	}

	// Create service clients
	client := &Client{
		conn:                  conn,
		AuthClient:            proto.NewAuthServiceClient(conn),
		InitClient:            proto.NewInitServiceClient(conn),
		BaseModelClient:       proto.NewBaseModelServiceClient(conn),
		MegaModelClient:       proto.NewMegaModelServiceClient(conn),
		EmbeddingModelClient:  proto.NewEmbeddingModelServiceClient(conn),
		StreamFinetunerClient: proto.NewStreamFineTunerServiceClient(conn),
		StreamGraphClient:     proto.NewStreamGraphServiceClient(conn),
		StreamConnectorClient: proto.NewStreamConnectorServiceClient(conn),
		KnowledgeBaseClient:   proto.NewKnowledgeBaseServiceClient(conn),
	}

	return client, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// WithContext returns a context with timeout
func (c *Client) WithContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 30*time.Second)
}
