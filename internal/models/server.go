package models

import "time"

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port            string        `json:"port"`
	ReadTimeout     time.Duration `json:"readTimeout"`
	WriteTimeout    time.Duration `json:"writeTimeout"`
	ShutdownTimeout time.Duration `json:"shutdownTimeout"`
	Version         string        `json:"version"`
}

// NewDefaultServerConfig creates a new server configuration with default values
func NewDefaultServerConfig(version string) *ServerConfig {
	return &ServerConfig{
		Port:            "8080",
		ReadTimeout:     time.Second * 15,
		WriteTimeout:    time.Second * 15,
		ShutdownTimeout: time.Second * 30,
		Version:         version,
	}
}

// MCPRequest represents a request to the MCP server
type MCPRequest struct {
	ID        string                 `json:"id"`
	Method    string                 `json:"method"`
	Params    map[string]interface{} `json:"params"`
	Timestamp time.Time              `json:"timestamp"`
}

// MCPResponse represents a response from the MCP server
type MCPResponse struct {
	ID        string                 `json:"id"`
	Result    map[string]interface{} `json:"result,omitempty"`
	Error     *MCPError              `json:"error,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// MCPError represents an error response
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
