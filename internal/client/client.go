package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://robot-ws.your-server.de"
)

// Client is the Hetzner Robot API client
type Client struct {
	Username   string
	Password   string
	HTTPClient *http.Client
}

// NewClient creates a new Hetzner Robot API client
func NewClient(username, password string) *Client {
	return &Client{
		Username: username,
		Password: password,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// Request makes an HTTP request to the Hetzner Robot API
func (c *Client) Request(method, path string, body interface{}, result interface{}) error {
	var requestBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		requestBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", baseURL, path), requestBody)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errorResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return fmt.Errorf("HTTP error %d", resp.StatusCode)
		}
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, errorResp.Error.Message)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return err
		}
	}

	return nil
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
