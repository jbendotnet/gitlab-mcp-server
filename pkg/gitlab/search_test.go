package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func TestSearchProjects(t *testing.T) {
	tests := []struct {
		name          string
		args          map[string]interface{}
		mockResponse  []*gitlab.Project
		mockError     error
		expectedError string
	}{
		{
			name: "successful search projects",
			args: map[string]interface{}{
				"query": "test",
			},
			mockResponse: []*gitlab.Project{
				{
					Name:        "test-project-1",
					Description: "Test Repository 1",
					Visibility:  "private",
				},
				{
					Name:        "test-project-2",
					Description: "Test Repository 2",
					Visibility:  "public",
				},
			},
			expectedError: "",
		},
		{
			name:          "missing required parameter",
			args:          map[string]interface{}{},
			expectedError: "missing required parameter: query",
		},
		{
			name: "GitLab API error",
			args: map[string]interface{}{
				"query": "test",
			},
			mockError:     fmt.Errorf("API error"),
			expectedError: "failed to search projects: API error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock client function
			getClient := func(ctx context.Context) (*gitlab.Client, error) {
				return &gitlab.Client{
					Search: &mockSearchService{
						searchProjectFunc: func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Project, *gitlab.Response, error) {
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
			_, handler := SearchProjects(getClient, translationHelper)

			// Create a request
			request := mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: tc.args,
				},
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

			// Verify the response contains the expected projects
			var projects []*gitlab.Project
			err = json.Unmarshal([]byte(result.Content[0].(*mcp.TextContent).Text), &projects)
			require.NoError(t, err)
			assert.Len(t, projects, len(tc.mockResponse))
			for i, project := range projects {
				assert.Equal(t, tc.mockResponse[i].Name, project.Name)
				assert.Equal(t, tc.mockResponse[i].Description, project.Description)
				assert.Equal(t, tc.mockResponse[i].Visibility, project.Visibility)
			}
		})
	}
}

func TestSearchMergeRequests(t *testing.T) {
	tests := []struct {
		name          string
		args          map[string]interface{}
		mockResponse  []*gitlab.MergeRequest
		mockError     error
		expectedError string
	}{
		{
			name: "successful search merge requests",
			args: map[string]interface{}{
				"query": "test",
			},
			mockResponse: []*gitlab.MergeRequest{
				{
					BasicMergeRequest: gitlab.BasicMergeRequest{
						IID:   1,
						Title: "Test MR 1",
						State: "opened",
					},
				},
				{
					BasicMergeRequest: gitlab.BasicMergeRequest{
						IID:   2,
						Title: "Test MR 2",
						State: "merged",
					},
				},
			},
			expectedError: "",
		},
		{
			name:          "missing required parameter",
			args:          map[string]interface{}{},
			expectedError: "missing required parameter: query",
		},
		{
			name: "GitLab API error",
			args: map[string]interface{}{
				"query": "test",
			},
			mockError:     fmt.Errorf("API error"),
			expectedError: "failed to search merge requests: API error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock client function
			getClient := func(ctx context.Context) (*gitlab.Client, error) {
				return &gitlab.Client{
					Search: &mockSearchService{
						searchMergeRequestFunc: func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
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
			_, handler := SearchMergeRequests(getClient, translationHelper)

			// Create a request
			request := mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: tc.args,
				},
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

			// Verify the response contains the expected merge requests
			var mrs []*gitlab.MergeRequest
			err = json.Unmarshal([]byte(result.Content[0].(*mcp.TextContent).Text), &mrs)
			require.NoError(t, err)
			assert.Len(t, mrs, len(tc.mockResponse))
			for i, mr := range mrs {
				assert.Equal(t, tc.mockResponse[i].Title, mr.Title)
				assert.Equal(t, tc.mockResponse[i].State, mr.State)
			}
		})
	}
}

func TestSearchUsers(t *testing.T) {
	tests := []struct {
		name          string
		args          map[string]interface{}
		mockResponse  []*gitlab.User
		mockError     error
		expectedError string
	}{
		{
			name: "successful search users",
			args: map[string]interface{}{
				"query": "test",
			},
			mockResponse: []*gitlab.User{
				{
					ID:       1,
					Name:     "Test User 1",
					Username: "testuser1",
				},
				{
					ID:       2,
					Name:     "Test User 2",
					Username: "testuser2",
				},
			},
			expectedError: "",
		},
		{
			name:          "missing required parameter",
			args:          map[string]interface{}{},
			expectedError: "missing required parameter: query",
		},
		{
			name: "GitLab API error",
			args: map[string]interface{}{
				"query": "test",
			},
			mockError:     fmt.Errorf("API error"),
			expectedError: "failed to search users: API error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock client function
			getClient := func(ctx context.Context) (*gitlab.Client, error) {
				return &gitlab.Client{
					Search: &mockSearchService{
						searchUserFunc: func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.User, *gitlab.Response, error) {
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
			_, handler := SearchUsers(getClient, translationHelper)

			// Create a request
			request := mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: tc.args,
				},
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

			// Verify the response contains the expected users
			var users []*gitlab.User
			err = json.Unmarshal([]byte(result.Content[0].(*mcp.TextContent).Text), &users)
			require.NoError(t, err)
			assert.Len(t, users, len(tc.mockResponse))
			for i, user := range users {
				assert.Equal(t, tc.mockResponse[i].Name, user.Name)
				assert.Equal(t, tc.mockResponse[i].Username, user.Username)
			}
		})
	}
}
