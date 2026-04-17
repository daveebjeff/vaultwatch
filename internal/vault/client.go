package vault

import (
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client.
type Client struct {
	api *vaultapi.Client
}

// NewClient creates a new Vault client from the given address and token.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}

	api.SetToken(token)

	return &Client{api: api}, nil
}

// IsAuthenticated checks whether the current token is valid.
func (c *Client) IsAuthenticated() error {
	_, err := c.api.Auth().Token().LookupSelf()
	if err != nil {
		return fmt.Errorf("vault authentication check failed: %w", err)
	}
	return nil
}

// ReadSecret reads a KV secret at the given path.
func (c *Client) ReadSecret(path string) (*vaultapi.Secret, error) {
	secret, err := c.api.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found at path %q", path)
	}
	return secret, nil
}
