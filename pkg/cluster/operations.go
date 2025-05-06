package cluster

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nstreama-ai/nstream-ai-cli/pkg/client"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/config"
	"github.com/nstreama-ai/nstream-ai-cli/pkg/utils"
	clusterproto "github.com/nstreama-ai/nstream-ai-cli/proto/cluster"
)

// Operations handles cluster-related operations
type Operations struct {
	client *client.Client
	config *config.Config
}

// NewOperations creates a new Operations instance
func NewOperations() (*Operations, error) {
	c, err := client.NewClient("", true, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	return &Operations{
		client: c,
		config: cfg,
	}, nil
}

// ListClusters lists all available clusters
func (o *Operations) ListClusters(ctx context.Context) ([]*clusterproto.Cluster, error) {
	if o.config.User.AuthToken == "" {
		return nil, fmt.Errorf("authentication token is missing. Please sign in first")
	}

	listResp, err := o.client.ClusterClient.ListClusters(ctx, &clusterproto.ListClustersRequest{
		AuthToken: o.config.User.AuthToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %v", err)
	}

	return listResp.Clusters, nil
}

// GetClusterDetails gets details for a specific cluster
func (o *Operations) GetClusterDetails(ctx context.Context, clusterName string) (*clusterproto.ClusterConfig, error) {
	if o.config.User.AuthToken == "" {
		return nil, fmt.Errorf("authentication token is missing. Please sign in first")
	}

	detailsResp, err := o.client.ClusterClient.GetClusterDetails(ctx, &clusterproto.GetClusterDetailsRequest{
		ClusterName: clusterName,
		AuthToken:   o.config.User.AuthToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster details: %v", err)
	}

	if detailsResp.Error != "" {
		return nil, fmt.Errorf("failed to get cluster details: %s", detailsResp.Error)
	}

	return detailsResp.Config, nil
}

// DisplayClusters displays clusters in a table format
func (o *Operations) DisplayClusters(clusters []*clusterproto.Cluster) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, utils.TableHeaderCluster)
	for i, cluster := range clusters {
		fmt.Fprintf(w, "%d. %s\t%s\t%s\t%s\t%s\n",
			i+1,
			cluster.Id,
			cluster.Region,
			cluster.CloudProvider,
			cluster.Bucket,
			cluster.Role,
		)
	}
	w.Flush()
}

// UpdateConfig updates the config with cluster details
func (o *Operations) UpdateConfig(clusterName string, details *clusterproto.ClusterConfig) error {
	o.config.Cluster = config.ClusterConfig{
		Name:          clusterName,
		Region:        details.Region,
		CloudProvider: details.CloudProvider,
		Bucket:        details.Bucket,
		Role:          details.Role,
		ClusterToken:  details.ClusterToken,
	}
	return config.SaveConfig(o.config)
}

// Close closes the client connection
func (o *Operations) Close() {
	if o.client != nil {
		o.client.Close()
	}
}
