package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jbendotnet/gitlab-mcp-server/pkg/translations"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type mockClient struct {
	Search gitlab.SearchServiceInterface
}

func (m *mockClient) SearchService() gitlab.SearchServiceInterface {
	return m.Search
}

func TestGetRepository(t *testing.T) {
	tests := []struct {
		name          string
		projectID     string
		searchFunc    func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Project, *gitlab.Response, error)
		expectedRepo  *gitlab.Project
		expectedError error
	}{
		{
			name:      "success",
			projectID: "123",
			searchFunc: func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Project, *gitlab.Response, error) {
				return []*gitlab.Project{
					{
						ID:                123,
						Name:              "test-repo",
						PathWithNamespace: "test-group/test-repo",
						WebURL:            "https://gitlab.com/test-group/test-repo",
					},
				}, &gitlab.Response{Response: &http.Response{StatusCode: 200}}, nil
			},
			expectedRepo: &gitlab.Project{
				ID:                123,
				Name:              "test-repo",
				PathWithNamespace: "test-group/test-repo",
				WebURL:            "https://gitlab.com/test-group/test-repo",
			},
			expectedError: nil,
		},
		{
			name:      "not found",
			projectID: "123",
			searchFunc: func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Project, *gitlab.Response, error) {
				return []*gitlab.Project{}, &gitlab.Response{Response: &http.Response{StatusCode: 200}}, nil
			},
			expectedRepo:  nil,
			expectedError: fmt.Errorf("repository not found"),
		},
		{
			name:      "api error",
			projectID: "123",
			searchFunc: func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Project, *gitlab.Response, error) {
				return nil, nil, assert.AnError
			},
			expectedRepo:  nil,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSearch := &mockSearchService{
				searchProjectFunc: tt.searchFunc,
			}

			getClient := func(ctx context.Context) (*gitlab.Client, error) {
				return &gitlab.Client{
					Search: mockSearch,
				}, nil
			}

			mockTranslationHelper := translations.TranslationHelperFunc(func(key string, defaultValue string) string {
				return defaultValue
			})

			tool, handler := GetRepository(getClient, mockTranslationHelper)

			assert.NotNil(t, tool)
			assert.NotNil(t, handler)

			result, err := handler(context.Background(), mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: map[string]interface{}{
						"namespace": "test-group",
						"project":   "test-repo",
					},
				},
			})

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestSearchProjectsTool(t *testing.T) {
	tests := []struct {
		name       string
		projectID  string
		searchFunc func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Project, *gitlab.Response, error)
		wantErr    bool
	}{
		{
			name:      "success",
			projectID: "123",
			searchFunc: func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Project, *gitlab.Response, error) {
				return []*gitlab.Project{
					{
						ID:          123,
						Name:        "test-project",
						Description: "Test project",
						Path:        "test/project",
					},
				}, nil, nil
			},
			wantErr: false,
		},
		{
			name:      "not found",
			projectID: "456",
			searchFunc: func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Project, *gitlab.Response, error) {
				return []*gitlab.Project{}, nil, nil
			},
			wantErr: false,
		},
		{
			name:      "api error",
			projectID: "789",
			searchFunc: func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Project, *gitlab.Response, error) {
				return nil, nil, assert.AnError
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSearch := &mockSearchService{
				searchProjectFunc: tt.searchFunc,
			}

			getClient := func(ctx context.Context) (*gitlab.Client, error) {
				return &gitlab.Client{
					Search: mockSearch,
				}, nil
			}

			mockTranslationHelper := translations.TranslationHelperFunc(func(key string, defaultValue string) string {
				return defaultValue
			})

			tool, handler := SearchProjects(getClient, mockTranslationHelper)

			assert.NotNil(t, tool)
			assert.NotNil(t, handler)

			result, err := handler(context.Background(), mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: map[string]interface{}{
						"query": tt.projectID,
					},
				},
			})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}
