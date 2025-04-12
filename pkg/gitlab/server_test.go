package gitlab

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		readOnly  bool
		wantTools int
	}{
		{
			name:      "read-only mode",
			version:   "1.0.0",
			readOnly:  true,
			wantTools: 14, // Number of tools in read-only mode
		},
		{
			name:      "read-write mode",
			version:   "1.0.0",
			readOnly:  false,
			wantTools: 20, // Number of tools in read-write mode
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock client function
			getClient := GetMockClientFn(t)

			// Create a mock translation helper
			translationHelper := func(key string, defaultValue string) string {
				return defaultValue
			}

			// Create the server
			s := NewServer(getClient, tc.version, tc.readOnly, translationHelper)

			// Create a context for the request
			ctx := context.Background()

			// Create a ListTools request
			request := struct {
				JSONRPC string        `json:"jsonrpc"`
				ID      int           `json:"id"`
				Method  mcp.MCPMethod `json:"method"`
			}{
				JSONRPC: mcp.JSONRPC_VERSION,
				ID:      1,
				Method:  mcp.MethodToolsList,
			}

			// Send the request
			rawRequest, err := json.Marshal(request)
			require.NoError(t, err)

			// Handle the request
			response := s.HandleMessage(ctx, rawRequest)
			require.NotNil(t, response)

			// Parse the response
			jsonResponse, ok := response.(mcp.JSONRPCResponse)
			require.True(t, ok)

			var result struct {
				Tools []mcp.Tool `json:"tools"`
			}
			resultBytes, err := json.Marshal(jsonResponse.Result)
			require.NoError(t, err)
			err = json.Unmarshal(resultBytes, &result)
			require.NoError(t, err)

			// Verify the server has the correct number of tools
			assert.Equal(t, tc.wantTools, len(result.Tools))
		})
	}
}

func TestServerWithOptions(t *testing.T) {
	tests := []struct {
		name          string
		options       []server.ServerOption
		expectedError bool
	}{
		{
			name:          "valid options",
			options:       []server.ServerOption{server.WithLogging()},
			expectedError: false,
		},
		{
			name:          "invalid options",
			options:       []server.ServerOption{nil},
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock client function
			getClient := GetMockClientFn(t)

			// Create a mock translation helper
			translationHelper := func(key string, defaultValue string) string {
				return defaultValue
			}

			// Create the server with options
			defer func() {
				if r := recover(); r != nil {
					if !tc.expectedError {
						t.Errorf("unexpected panic: %v", r)
					}
				} else if tc.expectedError {
					t.Error("expected panic but got none")
				}
			}()

			s := NewServer(getClient, "1.0.0", false, translationHelper, tc.options...)
			if !tc.expectedError {
				assert.NotNil(t, s)
			}
		})
	}
}

func TestServerResourceTemplates(t *testing.T) {
	// Create a mock client function
	getClient := func(ctx context.Context) (*gitlab.Client, error) {
		return &gitlab.Client{}, nil
	}

	// Create a mock translation helper
	translationHelper := func(key string, defaultValue string) string {
		return defaultValue
	}

	// Create the server with resource capabilities
	s := NewServer(getClient, "1.0.0", false, translationHelper)

	// Create a context for the request
	ctx := context.Background()

	// Now try to list the resource templates
	request := struct {
		JSONRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Method  string `json:"method"`
		Params  struct {
			Arguments map[string]interface{} `json:"arguments,omitempty"`
		} `json:"params"`
	}{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      1,
		Method:  string(mcp.MethodResourcesTemplatesList),
		Params: struct {
			Arguments map[string]interface{} `json:"arguments,omitempty"`
		}{
			Arguments: map[string]interface{}{},
		},
	}

	// Send the request
	rawRequest, err := json.Marshal(request)
	require.NoError(t, err)

	// Handle the request
	response := s.HandleMessage(ctx, rawRequest)
	require.NotNil(t, response)

	// Parse the response
	jsonResponse, ok := response.(mcp.JSONRPCResponse)
	require.True(t, ok)

	var result struct {
		ResourceTemplates []struct {
			URITemplate string `json:"uriTemplate"`
			Name        string `json:"name"`
		} `json:"resourceTemplates"`
	}
	resultBytes, err := json.Marshal(jsonResponse.Result)
	require.NoError(t, err)
	err = json.Unmarshal(resultBytes, &result)
	require.NoError(t, err)

	// Verify that the server has the correct number of resource templates
	expectedTemplates := 5 // Number of resource templates
	assert.Equal(t, expectedTemplates, len(result.ResourceTemplates))

	// Verify that each resource template has a valid pattern
	for _, template := range result.ResourceTemplates {
		assert.NotEmpty(t, template.URITemplate)
		assert.NotEmpty(t, template.Name)
	}
}

func TestServerToolHandlers(t *testing.T) {
	// Create a mock client function
	getClient := GetMockClientFn(t)

	// Create a mock translation helper
	translationHelper := func(key string, defaultValue string) string {
		return defaultValue
	}

	// Create the server
	s := NewServer(getClient, "1.0.0", false, translationHelper)

	// Create a context for the request
	ctx := context.Background()

	// Create a ListTools request
	request := struct {
		JSONRPC string        `json:"jsonrpc"`
		ID      int           `json:"id"`
		Method  mcp.MCPMethod `json:"method"`
	}{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      1,
		Method:  mcp.MethodToolsList,
	}

	// Send the request
	rawRequest, err := json.Marshal(request)
	require.NoError(t, err)

	// Handle the request
	response := s.HandleMessage(ctx, rawRequest)
	require.NotNil(t, response)

	// Parse the response
	jsonResponse, ok := response.(mcp.JSONRPCResponse)
	require.True(t, ok)

	var result struct {
		Tools []mcp.Tool `json:"tools"`
	}
	resultBytes, err := json.Marshal(jsonResponse.Result)
	require.NoError(t, err)
	err = json.Unmarshal(resultBytes, &result)
	require.NoError(t, err)

	// Verify that each tool has a valid handler
	for _, tool := range result.Tools {
		// We can't directly check the handler since it's not exposed in the ListTools response
		// Instead, we'll verify that the tool has a name and description
		assert.NotEmpty(t, tool.Name, "tool name should not be empty")
		assert.NotEmpty(t, tool.Description, "tool description should not be empty")
	}
}
