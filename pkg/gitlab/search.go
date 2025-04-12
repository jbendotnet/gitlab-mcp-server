package gitlab

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jbendotnet/gitlab-mcp-server/pkg/translations"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// SearchProjects implements the search projects tool
func SearchProjects(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"search_projects",
		mcp.WithDescription(t("TOOL_SEARCH_PROJECTS_DESCRIPTION", "Search for GitLab projects")),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description(t("PARAM_QUERY_DESCRIPTION", "Search query")),
		),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get GitLab client: %w", err)
		}

		query, err := RequiredString(r, "query")
		if err != nil {
			return nil, err
		}

		projects, _, err := client.Search.Projects(query, &gitlab.SearchOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to search projects: %w", err)
		}

		response, err := json.Marshal(projects)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: string(response),
				},
			},
		}, nil
	}

	return tool, handler
}

// SearchMergeRequests implements the search merge requests tool
func SearchMergeRequests(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"search_merge_requests",
		mcp.WithDescription(t("TOOL_SEARCH_MERGE_REQUESTS_DESCRIPTION", "Search for GitLab merge requests")),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description(t("PARAM_QUERY_DESCRIPTION", "Search query")),
		),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get GitLab client: %w", err)
		}

		query, err := RequiredString(r, "query")
		if err != nil {
			return nil, err
		}

		mergeRequests, _, err := client.Search.MergeRequests(query, &gitlab.SearchOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to search merge requests: %w", err)
		}

		response, err := json.Marshal(mergeRequests)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: string(response),
				},
			},
		}, nil
	}

	return tool, handler
}

// SearchUsers implements the search users tool
func SearchUsers(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"search_users",
		mcp.WithDescription(t("TOOL_SEARCH_USERS_DESCRIPTION", "Search for GitLab users")),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description(t("PARAM_QUERY_DESCRIPTION", "Search query")),
		),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get GitLab client: %w", err)
		}

		query, err := RequiredString(r, "query")
		if err != nil {
			return nil, err
		}

		users, _, err := client.Search.Users(query, &gitlab.SearchOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to search users: %w", err)
		}

		response, err := json.Marshal(users)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: string(response),
				},
			},
		}, nil
	}

	return tool, handler
}
