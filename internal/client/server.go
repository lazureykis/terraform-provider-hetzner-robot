package client

import (
	"fmt"
	"net/http"
)

// Server represents a Hetzner dedicated server
type Server struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ServerIP    string `json:"server_ip"`
	ServerIPv6  string `json:"server_ipv6,omitempty"`
	Product     string `json:"product"`
	Datacenter  string `json:"dc"`
	Status      string `json:"status"`
	RescueOS    string `json:"rescue,omitempty"`
	Cancelled   bool   `json:"cancelled,omitempty"`
	Paid        bool   `json:"paid"`
	Traffic     string `json:"traffic,omitempty"`
	Flatrate    bool   `json:"flatrate"`
	Distributed bool   `json:"distributed"`
}

// ServerResetOptions represents options for resetting a server
type ServerResetOptions struct {
	Type string `json:"type"` // e.g., "hw" for hardware reset, "sw" for software reset
}

// ServerRescueOptions represents options for configuring rescue mode
type ServerRescueOptions struct {
	OS      string   `json:"os"` // e.g., "linux64", "linux32", "freebsd64"
	SSHKeys []string `json:"ssh_keys,omitempty"`
}

// GetServer retrieves a server by its ID
func (c *Client) GetServer(id string) (*Server, error) {
	path := fmt.Sprintf("/server/%s", id)
	var response struct {
		Server Server `json:"server"`
	}

	err := c.Request(http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Server, nil
}

// ListServers retrieves all servers
func (c *Client) ListServers() ([]Server, error) {
	var response struct {
		Servers []Server `json:"server"`
	}

	err := c.Request(http.MethodGet, "/server", nil, &response)
	if err != nil {
		return nil, err
	}

	return response.Servers, nil
}

// UpdateServer updates a server's name
func (c *Client) UpdateServer(id string, name string) (*Server, error) {
	path := fmt.Sprintf("/server/%s", id)
	body := map[string]string{
		"name": name,
	}

	var response struct {
		Server Server `json:"server"`
	}

	err := c.Request(http.MethodPost, path, body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Server, nil
}

// ResetServer performs a reset on the server
func (c *Client) ResetServer(id string, options ServerResetOptions) error {
	path := fmt.Sprintf("/server/%s/reset", id)
	return c.Request(http.MethodPost, path, options, nil)
}

// EnableRescueMode enables rescue mode on a server
func (c *Client) EnableRescueMode(id string, options ServerRescueOptions) error {
	path := fmt.Sprintf("/server/%s/rescue", id)
	return c.Request(http.MethodPost, path, options, nil)
}

// DisableRescueMode disables rescue mode on a server
func (c *Client) DisableRescueMode(id string) error {
	path := fmt.Sprintf("/server/%s/rescue", id)
	return c.Request(http.MethodDelete, path, nil, nil)
}
