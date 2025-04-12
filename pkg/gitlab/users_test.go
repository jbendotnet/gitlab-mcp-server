package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func TestGetMe(t *testing.T) {
	tests := []struct {
		name          string
		mockResponse  *gitlab.User
		mockError     error
		expectedError string
	}{
		{
			name: "successful get current user",
			mockResponse: &gitlab.User{
				ID:       1,
				Name:     "Test User",
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectedError: "",
		},
		{
			name:          "GitLab API error",
			mockError:     fmt.Errorf("API error"),
			expectedError: "failed to get current user: API error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock client function
			getClient := func(ctx context.Context) (*gitlab.Client, error) {
				return &gitlab.Client{
					Users: &mockUsersService{
						currentUserFunc: func(options ...gitlab.RequestOptionFunc) (*gitlab.User, *gitlab.Response, error) {
							return tc.mockResponse, &gitlab.Response{Response: &http.Response{StatusCode: http.StatusOK}}, tc.mockError
						},
					},
				}, nil
			}

			// Create a mock translation helper
			translationHelper := func(key string, defaultValue string) string {
				return defaultValue
			}

			// Get the tool and handler
			_, handler := GetMe(getClient, translationHelper)

			// Create a request
			request := mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{},
			}

			// Call the handler
			result, err := handler(context.Background(), request)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			require.Len(t, result.Content, 1)

			textContent, ok := result.Content[0].(*mcp.TextContent)
			require.True(t, ok)
			assert.Contains(t, textContent.Text, tc.mockResponse.Username)
			assert.Contains(t, textContent.Text, tc.mockResponse.Name)
			assert.Contains(t, textContent.Text, tc.mockResponse.Email)
		})
	}
}
