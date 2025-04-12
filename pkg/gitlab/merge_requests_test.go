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

// mockMergeRequestsService is a mock implementation of the GitLab merge requests service
type mockMergeRequestsService struct {
	getFunc         func(pid interface{}, mr int, opt *gitlab.GetMergeRequestsOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
	listProjectFunc func(pid interface{}, opt *gitlab.ListProjectMergeRequestsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.BasicMergeRequest, *gitlab.Response, error)
}

func (m *mockMergeRequestsService) GetMergeRequest(pid interface{}, mr int, opt *gitlab.GetMergeRequestsOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	return m.getFunc(pid, mr, opt, options...)
}

func (m *mockMergeRequestsService) ListProjectMergeRequests(pid interface{}, opt *gitlab.ListProjectMergeRequestsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.BasicMergeRequest, *gitlab.Response, error) {
	return m.listProjectFunc(pid, opt, options...)
}

func (m *mockMergeRequestsService) AcceptMergeRequest(pid interface{}, mr int, opt *gitlab.AcceptMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) AddSpentTime(pid interface{}, mr int, opt *gitlab.AddSpentTimeOptions, options ...gitlab.RequestOptionFunc) (*gitlab.TimeStats, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) CancelMergeWhenPipelineSucceeds(pid interface{}, mr int, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) CreateMergeRequest(pid interface{}, opt *gitlab.CreateMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) CreateMergeRequestDependency(pid interface{}, mergeRequest int, opts gitlab.CreateMergeRequestDependencyOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequestDependency, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) CreateMergeRequestPipeline(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) (*gitlab.PipelineInfo, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) DeleteMergeRequest(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	return nil, nil
}

func (m *mockMergeRequestsService) DeleteMergeRequestDependency(pid interface{}, mergeRequest int, blockingMergeRequest int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	return nil, nil
}

func (m *mockMergeRequestsService) GetMergeRequestApprovals(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequestApprovals, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) GetMergeRequestChanges(pid interface{}, mergeRequest int, opt *gitlab.GetMergeRequestChangesOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) GetMergeRequestCommits(pid interface{}, mergeRequest int, opt *gitlab.GetMergeRequestCommitsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Commit, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) GetMergeRequestDependencies(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) ([]gitlab.MergeRequestDependency, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) GetSingleMergeRequestDiffVersion(pid interface{}, mergeRequest, version int, opt *gitlab.GetSingleMergeRequestDiffVersionOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequestDiffVersion, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) GetTimeSpent(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) (*gitlab.TimeStats, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) ListGroupMergeRequests(gid interface{}, opt *gitlab.ListGroupMergeRequestsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.BasicMergeRequest, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) ListMergeRequestDiffs(pid interface{}, mergeRequest int, opt *gitlab.ListMergeRequestDiffsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequestDiff, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) ListMergeRequests(opt *gitlab.ListMergeRequestsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.BasicMergeRequest, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) ResetSpentTime(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) (*gitlab.TimeStats, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) ResetTimeEstimate(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) (*gitlab.TimeStats, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) SetTimeEstimate(pid interface{}, mergeRequest int, opt *gitlab.SetTimeEstimateOptions, options ...gitlab.RequestOptionFunc) (*gitlab.TimeStats, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) SubscribeToMergeRequest(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) UnsubscribeFromMergeRequest(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) UpdateMergeRequest(pid interface{}, mergeRequest int, opt *gitlab.UpdateMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) CreateTodo(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) (*gitlab.Todo, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) GetIssuesClosedOnMerge(pid interface{}, mergeRequest int, opt *gitlab.GetIssuesClosedOnMergeOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Issue, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) GetMergeRequestDiffVersions(pid interface{}, mergeRequest int, opt *gitlab.GetMergeRequestDiffVersionsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequestDiffVersion, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) GetMergeRequestParticipants(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) ([]*gitlab.BasicUser, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) GetMergeRequestReviewers(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequestReviewer, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) ListMergeRequestPipelines(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) ([]*gitlab.PipelineInfo, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockMergeRequestsService) RebaseMergeRequest(pid interface{}, mergeRequest int, opt *gitlab.RebaseMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	return nil, nil
}

func (m *mockMergeRequestsService) ShowMergeRequestRawDiffs(pid interface{}, mergeRequest int, opt *gitlab.ShowMergeRequestRawDiffsOptions, options ...gitlab.RequestOptionFunc) ([]byte, *gitlab.Response, error) {
	return nil, nil, nil
}

func TestGetMergeRequest(t *testing.T) {
	tests := []struct {
		name          string
		args          map[string]interface{}
		mockResponse  *gitlab.MergeRequest
		mockError     error
		expectedError string
	}{
		{
			name: "successful get merge request",
			args: map[string]interface{}{
				"namespace": "test-namespace",
				"project":   "test-project",
				"id":        float64(1),
			},
			mockResponse: &gitlab.MergeRequest{
				BasicMergeRequest: gitlab.BasicMergeRequest{
					IID:   1,
					Title: "Test MR",
					State: "opened",
					Author: &gitlab.BasicUser{
						Name: "Test User",
					},
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
			expectedError: "failed to get merge request: API error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock client function
			getClient := func(ctx context.Context) (*gitlab.Client, error) {
				return &gitlab.Client{
					MergeRequests: &mockMergeRequestsService{
						getFunc: func(pid interface{}, mr int, opt *gitlab.GetMergeRequestsOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
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
			_, handler := GetMergeRequest(getClient, translationHelper)

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

func TestListMergeRequests(t *testing.T) {
	tests := []struct {
		name          string
		args          map[string]interface{}
		mockResponse  []*gitlab.BasicMergeRequest
		mockError     error
		expectedError string
	}{
		{
			name: "successful list merge requests",
			args: map[string]interface{}{
				"namespace": "test-namespace",
				"project":   "test-project",
			},
			mockResponse: []*gitlab.BasicMergeRequest{
				{
					IID:   1,
					Title: "Test MR 1",
					State: "opened",
				},
				{
					IID:   2,
					Title: "Test MR 2",
					State: "merged",
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
				"order_by":  "created_at",
				"sort":      "desc",
			},
			mockResponse: []*gitlab.BasicMergeRequest{
				{
					IID:   1,
					Title: "Test MR",
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
			expectedError: "failed to list merge requests: API error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock client function
			getClient := func(ctx context.Context) (*gitlab.Client, error) {
				return &gitlab.Client{
					MergeRequests: &mockMergeRequestsService{
						listProjectFunc: func(pid interface{}, opt *gitlab.ListProjectMergeRequestsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.BasicMergeRequest, *gitlab.Response, error) {
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
			_, handler := ListMergeRequests(getClient, translationHelper)

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
			var mrs []*gitlab.BasicMergeRequest
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
