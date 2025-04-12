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

// mockIssuesService is a mock implementation of the GitLab issues service
type mockIssuesService struct {
	getFunc         func(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error)
	listProjectFunc func(pid interface{}, opt *gitlab.ListProjectIssuesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Issue, *gitlab.Response, error)
}

func (m *mockIssuesService) GetIssue(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error) {
	return m.getFunc(pid, issue, options...)
}

func (m *mockIssuesService) ListProjectIssues(pid interface{}, opt *gitlab.ListProjectIssuesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Issue, *gitlab.Response, error) {
	return m.listProjectFunc(pid, opt, options...)
}

func (m *mockIssuesService) AddSpentTime(pid interface{}, issue int, opt *gitlab.AddSpentTimeOptions, options ...gitlab.RequestOptionFunc) (*gitlab.TimeStats, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) CreateIssue(pid interface{}, opt *gitlab.CreateIssueOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) CreateTodo(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.Todo, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) DeleteIssue(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	return nil, nil
}

func (m *mockIssuesService) GetIssueByID(issue int, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) GetParticipants(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) ([]*gitlab.BasicUser, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) GetTimeSpent(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.TimeStats, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) ListGroupIssues(pid interface{}, opt *gitlab.ListGroupIssuesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Issue, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) ListIssues(opt *gitlab.ListIssuesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Issue, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) ListMergeRequestsClosingIssue(pid interface{}, issue int, opt *gitlab.ListMergeRequestsClosingIssueOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.BasicMergeRequest, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) ListMergeRequestsRelatedToIssue(pid interface{}, issue int, opt *gitlab.ListMergeRequestsRelatedToIssueOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.BasicMergeRequest, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) MoveIssue(pid interface{}, issue int, opt *gitlab.MoveIssueOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) ReorderIssue(pid interface{}, issue int, opt *gitlab.ReorderIssueOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) ResetSpentTime(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.TimeStats, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) ResetTimeEstimate(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.TimeStats, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) SetTimeEstimate(pid interface{}, issue int, opt *gitlab.SetTimeEstimateOptions, options ...gitlab.RequestOptionFunc) (*gitlab.TimeStats, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) SubscribeToIssue(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) UnsubscribeFromIssue(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockIssuesService) UpdateIssue(pid interface{}, issue int, opt *gitlab.UpdateIssueOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error) {
	return nil, nil, nil
}

func TestGetIssue(t *testing.T) {
	tests := []struct {
		name          string
		args          map[string]interface{}
		mockResponse  *gitlab.Issue
		mockError     error
		expectedError string
	}{
		{
			name: "successful get issue",
			args: map[string]interface{}{
				"namespace": "test-namespace",
				"project":   "test-project",
				"id":        float64(1),
			},
			mockResponse: &gitlab.Issue{
				IID:   1,
				Title: "Test Issue",
				State: "opened",
				Author: &gitlab.IssueAuthor{
					Name: "Test User",
				},
			},
			expectedError: "",
		},
		{
			name: "missing required parameter",
			args: map[string]interface{}{
				"namespace": "test-namespace",
				"project":   "test-project",
			},
			expectedError: "missing required parameter: id",
		},
		{
			name: "invalid parameter type",
			args: map[string]interface{}{
				"namespace": "test-namespace",
				"project":   "test-project",
				"id":        "not-a-number",
			},
			expectedError: "parameter id is not of type float64",
		},
		{
			name: "GitLab API error",
			args: map[string]interface{}{
				"namespace": "test-namespace",
				"project":   "test-project",
				"id":        float64(1),
			},
			mockError:     fmt.Errorf("API error"),
			expectedError: "failed to get issue: API error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock client function
			getClient := func(ctx context.Context) (*gitlab.Client, error) {
				return &gitlab.Client{
					Issues: &mockIssuesService{
						getFunc: func(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error) {
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
			_, handler := GetIssue(getClient, translationHelper)

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
			require.Len(t, result.Content, 1)

			textContent, ok := result.Content[0].(*mcp.TextContent)
			require.True(t, ok)
			assert.Contains(t, textContent.Text, tc.mockResponse.Title)
			assert.Contains(t, textContent.Text, tc.mockResponse.State)
			assert.Contains(t, textContent.Text, tc.mockResponse.Author.Name)
		})
	}
}

func TestListIssues(t *testing.T) {
	tests := []struct {
		name          string
		args          map[string]interface{}
		mockResponse  []*gitlab.Issue
		mockError     error
		expectedError string
	}{
		{
			name: "successful list issues",
			args: map[string]interface{}{
				"namespace": "test-namespace",
				"project":   "test-project",
			},
			mockResponse: []*gitlab.Issue{
				{
					IID:   1,
					Title: "Test Issue 1",
					State: "opened",
				},
				{
					IID:   2,
					Title: "Test Issue 2",
					State: "closed",
				},
			},
			expectedError: "",
		},
		{
			name: "successful list with filters",
			args: map[string]interface{}{
				"namespace": "test-namespace",
				"project":   "test-project",
				"state":     "opened",
			},
			mockResponse: []*gitlab.Issue{
				{
					IID:   1,
					Title: "Test Issue",
					State: "opened",
				},
			},
			expectedError: "",
		},
		{
			name: "missing required parameter",
			args: map[string]interface{}{
				"namespace": "test-namespace",
			},
			expectedError: "missing required parameter: project",
		},
		{
			name: "GitLab API error",
			args: map[string]interface{}{
				"namespace": "test-namespace",
				"project":   "test-project",
			},
			mockError:     fmt.Errorf("API error"),
			expectedError: "failed to list issues: API error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock client function
			getClient := func(ctx context.Context) (*gitlab.Client, error) {
				return &gitlab.Client{
					Issues: &mockIssuesService{
						listProjectFunc: func(pid interface{}, opt *gitlab.ListProjectIssuesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Issue, *gitlab.Response, error) {
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
			_, handler := ListIssues(getClient, translationHelper)

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

			// Verify the response contains the expected issues
			var issues []*gitlab.Issue
			err = json.Unmarshal([]byte(result.Content[0].(*mcp.TextContent).Text), &issues)
			require.NoError(t, err)
			assert.Len(t, issues, len(tc.mockResponse))
			for i, issue := range issues {
				assert.Equal(t, tc.mockResponse[i].Title, issue.Title)
				assert.Equal(t, tc.mockResponse[i].State, issue.State)
			}
		})
	}
}
