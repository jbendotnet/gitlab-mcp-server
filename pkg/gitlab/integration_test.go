package gitlab

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// TestIntegrationGitLab tests the integration with a real GitLab instance
func TestIntegrationGitLab(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("GITLAB_INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test. Set GITLAB_INTEGRATION_TEST=true to run.")
	}

	// Get GitLab credentials from environment
	token := os.Getenv("GITLAB_TOKEN")
	baseURL := os.Getenv("GITLAB_BASE_URL")
	if token == "" || baseURL == "" {
		t.Fatal("GITLAB_TOKEN and GITLAB_BASE_URL environment variables must be set for integration tests")
	}

	// Create GitLab client
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	require.NoError(t, err)

	// Create client function
	getClient := func(ctx context.Context) (*gitlab.Client, error) {
		return client, nil
	}

	// Create translation helper
	translationHelper := func(key string, defaultValue string) string {
		return defaultValue
	}

	// Create server
	s := NewServer(getClient, "1.0.0", false, translationHelper)

	// Test repository operations
	t.Run("Repository Operations", func(t *testing.T) {
		// Create a test project
		project, _, err := client.Projects.CreateProject(&gitlab.CreateProjectOptions{
			Name:        gitlab.String("test-project"),
			Description: gitlab.String("Test project for integration tests"),
			Visibility:  gitlab.Visibility(gitlab.PrivateVisibility),
		})
		require.NoError(t, err)
		defer func() {
			_, err := client.Projects.DeleteProject(project.ID, nil)
			assert.NoError(t, err)
		}()

		// Test getting repository
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_repository",
				Arguments: map[string]interface{}{
					"namespace": project.Namespace.FullPath,
					"project":   project.Path,
				},
			},
		}

		requestBytes, err := json.Marshal(request)
		require.NoError(t, err)

		result := s.HandleMessage(context.Background(), requestBytes)
		require.NotNil(t, result)
	})

	// Test issue operations
	t.Run("Issue Operations", func(t *testing.T) {
		// Create a test project
		project, _, err := client.Projects.CreateProject(&gitlab.CreateProjectOptions{
			Name:        gitlab.String("test-issues"),
			Description: gitlab.String("Test project for issue operations"),
			Visibility:  gitlab.Visibility(gitlab.PrivateVisibility),
		})
		require.NoError(t, err)
		defer func() {
			_, err := client.Projects.DeleteProject(project.ID, nil)
			assert.NoError(t, err)
		}()

		// Create a test issue
		issue, _, err := client.Issues.CreateIssue(project.ID, &gitlab.CreateIssueOptions{
			Title:       gitlab.String("Test Issue"),
			Description: gitlab.String("This is a test issue"),
		})
		require.NoError(t, err)

		// Test getting issue
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_issue",
				Arguments: map[string]interface{}{
					"namespace": project.Namespace.FullPath,
					"project":   project.Path,
					"id":        strconv.Itoa(issue.IID),
				},
			},
		}

		requestBytes, err := json.Marshal(request)
		require.NoError(t, err)

		result := s.HandleMessage(context.Background(), requestBytes)
		require.NotNil(t, result)
	})

	// Test merge request operations
	t.Run("Merge Request Operations", func(t *testing.T) {
		// Create a test project
		project, _, err := client.Projects.CreateProject(&gitlab.CreateProjectOptions{
			Name:        gitlab.String("test-mrs"),
			Description: gitlab.String("Test project for merge request operations"),
			Visibility:  gitlab.Visibility(gitlab.PrivateVisibility),
		})
		require.NoError(t, err)
		defer func() {
			_, err := client.Projects.DeleteProject(project.ID, nil)
			assert.NoError(t, err)
		}()

		// Create a test branch
		branch, _, err := client.Branches.CreateBranch(project.ID, &gitlab.CreateBranchOptions{
			Branch: gitlab.String("test-branch"),
			Ref:    gitlab.String("main"),
		})
		require.NoError(t, err)

		// Create a test merge request
		mr, _, err := client.MergeRequests.CreateMergeRequest(project.ID, &gitlab.CreateMergeRequestOptions{
			Title:        gitlab.String("Test MR"),
			Description:  gitlab.String("This is a test merge request"),
			SourceBranch: gitlab.String(branch.Name),
			TargetBranch: gitlab.String("main"),
		})
		require.NoError(t, err)

		// Test getting merge request
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_merge_request",
				Arguments: map[string]interface{}{
					"namespace": project.Namespace.FullPath,
					"project":   project.Path,
					"id":        strconv.Itoa(mr.IID),
				},
			},
		}

		requestBytes, err := json.Marshal(request)
		require.NoError(t, err)

		result := s.HandleMessage(context.Background(), requestBytes)
		require.NotNil(t, result)
	})

	// Test search operations
	t.Run("Search Operations", func(t *testing.T) {
		// Test searching for projects
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "search_projects",
				Arguments: map[string]interface{}{
					"query": "test",
				},
			},
		}

		requestBytes, err := json.Marshal(request)
		require.NoError(t, err)

		result := s.HandleMessage(context.Background(), requestBytes)
		require.NotNil(t, result)
	})

	// Test user operations
	t.Run("User Operations", func(t *testing.T) {
		// Test getting current user
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_me",
			},
		}

		requestBytes, err := json.Marshal(request)
		require.NoError(t, err)

		result := s.HandleMessage(context.Background(), requestBytes)
		require.NotNil(t, result)
	})
}
