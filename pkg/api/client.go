package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseURL = "https://api.nstream.ai/v1"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

type SignUpRequest struct {
	Email   string `json:"email"`
	OrgName string `json:"org_name"`
	Role    string `json:"role"`
}

type SignInRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ClusterRequest struct {
	Name          string `json:"name"`
	CloudProvider string `json:"cloud_provider"`
	Region        string `json:"region"`
	Bucket        string `json:"bucket"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ClusterResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func (c *Client) SignUp(req SignUpRequest) (*AuthResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/auth/signup", c.baseURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return nil, err
	}

	return &authResp, nil
}

func (c *Client) SignIn(req SignInRequest) (*AuthResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/auth/signin", c.baseURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return nil, err
	}

	return &authResp, nil
}

func (c *Client) CreateCluster(req ClusterRequest, authToken string) (*ClusterResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/clusters", c.baseURL),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var clusterResp ClusterResponse
	err = json.NewDecoder(resp.Body).Decode(&clusterResp)
	if err != nil {
		return nil, err
	}

	return &clusterResp, nil
}

func (c *Client) ListClusters(authToken string) ([]ClusterResponse, error) {
	httpReq, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/clusters", c.baseURL),
		nil,
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var clusters []ClusterResponse
	err = json.NewDecoder(resp.Body).Decode(&clusters)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}
