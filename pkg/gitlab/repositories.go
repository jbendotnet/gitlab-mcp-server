package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jbendotnet/gitlab-mcp-server/pkg/translations"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// RequiredString is a helper function to get a required string parameter from a request
func RequiredString(r mcp.CallToolRequest, p string) (string, error) {
	val, ok := r.Params.Arguments[p]
	if !ok {
		return "", fmt.Errorf("missing required parameter: %s", p)
	}

	strVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("parameter %s is not a string", p)
	}

	return strVal, nil
}

// GetRepository returns a tool for getting a specific repository
func GetRepository(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"get_repository",
		mcp.WithDescription(t("TOOL_GET_REPOSITORY_DESCRIPTION", "Get a specific repository")),
		mcp.WithString("namespace",
			mcp.Required(),
			mcp.Description(t("PARAM_NAMESPACE_DESCRIPTION", "The namespace of the project")),
		),
		mcp.WithString("project",
			mcp.Required(),
			mcp.Description(t("PARAM_PROJECT_DESCRIPTION", "The name of the project")),
		),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get GitLab client: %w", err)
		}

		namespace, err := RequiredString(r, "namespace")
		if err != nil {
			return nil, err
		}
		projectName, err := RequiredString(r, "project")
		if err != nil {
			return nil, err
		}

		query := fmt.Sprintf("%s/%s", namespace, projectName)
		projects, _, err := client.Search.Projects(query, &gitlab.SearchOptions{})
		if err != nil {
			return nil, err
		}

		if len(projects) == 0 {
			return nil, fmt.Errorf("repository not found")
		}

		project := projects[0]
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Name: %s\nPath: %s\nURL: %s", project.Name, project.PathWithNamespace, project.WebURL),
				},
			},
		}, nil
	}

	return tool, handler
}

// ListRepositories returns a tool for listing repositories
func ListRepositories(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"list_repositories",
		mcp.WithDescription(t("TOOL_LIST_REPOSITORIES_DESCRIPTION", "List repositories")),
		mcp.WithString("namespace",
			mcp.Description(t("PARAM_NAMESPACE_DESCRIPTION", "The namespace to list repositories from")),
		),
		mcp.WithString("search",
			mcp.Description(t("PARAM_SEARCH_DESCRIPTION", "Search query")),
		),
		mcp.WithString("order_by",
			mcp.Description(t("PARAM_ORDER_BY_DESCRIPTION", "Order by field")),
		),
		mcp.WithString("sort",
			mcp.Description(t("PARAM_SORT_DESCRIPTION", "Sort order")),
		),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get GitLab client: %v", err)), nil
		}

		opts := &gitlab.ListProjectsOptions{}

		if namespace, err := OptionalParam[string](r, "namespace"); err == nil && namespace != "" {
			owned := true
			opts.Owned = &owned
		}
		if search, err := OptionalParam[string](r, "search"); err == nil && search != "" {
			opts.Search = &search
		}
		if orderBy, err := OptionalParam[string](r, "order_by"); err == nil && orderBy != "" {
			opts.OrderBy = &orderBy
		}
		if sort, err := OptionalParam[string](r, "sort"); err == nil && sort != "" {
			opts.Sort = &sort
		}

		projects, resp, err := client.Projects.ListProjects(opts)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to list repositories: %v", err)), nil
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to read response body: %v", err)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("failed to list repositories: %s", string(body))), nil
		}

		jsonData, err := json.Marshal(projects)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	}

	return tool, handler
}

// SearchRepositories returns a tool for searching repositories
func SearchRepositories(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"search_repositories",
		mcp.WithDescription(t("TOOL_SEARCH_REPOSITORIES_DESCRIPTION", "Search for repositories")),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description(t("PARAM_QUERY_DESCRIPTION", "Search query")),
		),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get GitLab client: %v", err)), nil
		}

		query, err := requiredParam[string](r, "query")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		projects, resp, err := client.Projects.ListProjects(&gitlab.ListProjectsOptions{
			Search: &query,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to search repositories: %v", err)), nil
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to read response body: %v", err)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("failed to search repositories: %s", string(body))), nil
		}

		jsonData, err := json.Marshal(projects)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	}

	return tool, handler
}
