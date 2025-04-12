package gitlab

import (
	"context"
	"fmt"

	"github.com/jbendotnet/gitlab-mcp-server/pkg/translations"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// GetClientFn is a function type that returns a GitLab client
type GetClientFn func(context.Context) (*gitlab.Client, error)

// NewServer creates a new GitLab MCP server with the specified client and logger
func NewServer(getClient GetClientFn, version string, readOnly bool, t translations.TranslationHelperFunc, opts ...server.ServerOption) *server.MCPServer {
	// Add default options
	defaultOpts := []server.ServerOption{
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	}
	opts = append(defaultOpts, opts...)

	// Create a new MCP server
	s := server.NewMCPServer(
		"gitlab-mcp-server",
		version,
		opts...,
	)

	// Add GitLab Resources
	template, handler := GetRepositoryResourceContent(getClient, t)
	s.AddResourceTemplate(template, handler)

	template, handler = GetRepositoryResourceBranchContent(getClient, t)
	s.AddResourceTemplate(template, handler)

	template, handler = GetRepositoryResourceCommitContent(getClient, t)
	s.AddResourceTemplate(template, handler)

	template, handler = GetRepositoryResourceTagContent(getClient, t)
	s.AddResourceTemplate(template, handler)

	template, handler = GetRepositoryResourceMergeRequestContent(getClient, t)
	s.AddResourceTemplate(template, handler)

	// Add GitLab tools - Issues
	tool, toolHandler := GetIssue(getClient, t)
	s.AddTool(tool, toolHandler)

	tool, toolHandler = SearchIssues(getClient, t)
	s.AddTool(tool, toolHandler)

	tool, toolHandler = ListIssues(getClient, t)
	s.AddTool(tool, toolHandler)

	tool, toolHandler = GetIssueComments(getClient, t)
	s.AddTool(tool, toolHandler)

	if !readOnly {
		tool, toolHandler = CreateIssue(getClient, t)
		s.AddTool(tool, toolHandler)

		tool, toolHandler = AddIssueComment(getClient, t)
		s.AddTool(tool, toolHandler)

		tool, toolHandler = UpdateIssue(getClient, t)
		s.AddTool(tool, toolHandler)
	}

	// Add GitLab tools - Merge Requests
	tool, toolHandler = GetMergeRequest(getClient, t)
	s.AddTool(tool, toolHandler)

	tool, toolHandler = ListMergeRequests(getClient, t)
	s.AddTool(tool, toolHandler)

	tool, toolHandler = GetMergeRequestComments(getClient, t)
	s.AddTool(tool, toolHandler)

	if !readOnly {
		tool, toolHandler = CreateMergeRequest(getClient, t)
		s.AddTool(tool, toolHandler)

		tool, toolHandler = AddMergeRequestComment(getClient, t)
		s.AddTool(tool, toolHandler)

		tool, toolHandler = UpdateMergeRequest(getClient, t)
		s.AddTool(tool, toolHandler)
	}

	// Add GitLab tools - Repository
	tool, toolHandler = GetRepository(getClient, t)
	s.AddTool(tool, toolHandler)

	tool, toolHandler = ListRepositories(getClient, t)
	s.AddTool(tool, toolHandler)

	tool, toolHandler = SearchRepositories(getClient, t)
	s.AddTool(tool, toolHandler)

	// Add GitLab tools - Search
	tool, toolHandler = SearchProjects(getClient, t)
	s.AddTool(tool, toolHandler)

	tool, toolHandler = SearchMergeRequests(getClient, t)
	s.AddTool(tool, toolHandler)

	tool, toolHandler = SearchUsers(getClient, t)
	s.AddTool(tool, toolHandler)

	// Add GitLab tools - User
	tool, toolHandler = GetMe(getClient, t)
	s.AddTool(tool, toolHandler)

	return s
}

// GetMe returns information about the authenticated user
func GetMe(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"get_me",
		mcp.WithDescription(t("TOOL_GET_ME_DESCRIPTION", "Get information about the authenticated user")),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get GitLab client: %w", err)
		}

		user, _, err := client.Users.CurrentUser()
		if err != nil {
			return nil, fmt.Errorf("failed to get current user: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Username: %s\nName: %s\nEmail: %s", user.Username, user.Name, user.Email),
				},
			},
		}, nil
	}

	return tool, handler
}
