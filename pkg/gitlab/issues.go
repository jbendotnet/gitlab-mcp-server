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

// requiredParam is a helper function that can be used to fetch a requested parameter from the request.
// It does the following checks:
// 1. Checks if the parameter is present in the request.
// 2. Checks if the parameter is of the expected type.
// 3. Checks if the parameter is not empty, i.e: non-zero value
func requiredParam[T comparable](r mcp.CallToolRequest, p string) (T, error) {
	var zero T

	// Check if the parameter is present in the request
	if _, ok := r.Params.Arguments[p]; !ok {
		return zero, fmt.Errorf("missing required parameter: %s", p)
	}

	// Check if the parameter is of the expected type
	if _, ok := r.Params.Arguments[p].(T); !ok {
		return zero, fmt.Errorf("parameter %s is not of type %T", p, zero)
	}

	if r.Params.Arguments[p].(T) == zero {
		return zero, fmt.Errorf("missing required parameter: %s", p)
	}

	return r.Params.Arguments[p].(T), nil
}

// RequiredInt is a helper function that can be used to fetch a requested parameter from the request.
// It does the following checks:
// 1. Checks if the parameter is present in the request.
// 2. Checks if the parameter is of the expected type.
// 3. Checks if the parameter is not empty, i.e: non-zero value
func RequiredInt(r mcp.CallToolRequest, p string) (int, error) {
	v, err := requiredParam[float64](r, p)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

// OptionalParam is a helper function that can be used to fetch a requested parameter from the request.
// It does the following checks:
// 1. Checks if the parameter is present in the request, if not, it returns its zero-value
// 2. If it is present, it checks if the parameter is of the expected type and returns it
func OptionalParam[T any](r mcp.CallToolRequest, p string) (T, error) {
	var zero T

	// Check if the parameter is present in the request
	if _, ok := r.Params.Arguments[p]; !ok {
		return zero, nil
	}

	// Check if the parameter is of the expected type
	if _, ok := r.Params.Arguments[p].(T); !ok {
		return zero, fmt.Errorf("parameter %s is not of type %T, is %T", p, zero, r.Params.Arguments[p])
	}

	return r.Params.Arguments[p].(T), nil
}

// GetIssue returns a tool for getting a specific issue
func GetIssue(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"get_issue",
		mcp.WithDescription(t("TOOL_GET_ISSUE_DESCRIPTION", "Get a specific issue")),
		mcp.WithString("namespace",
			mcp.Required(),
			mcp.Description(t("PARAM_NAMESPACE_DESCRIPTION", "The namespace of the project")),
		),
		mcp.WithString("project",
			mcp.Required(),
			mcp.Description(t("PARAM_PROJECT_DESCRIPTION", "The name of the project")),
		),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description(t("PARAM_ISSUE_ID_DESCRIPTION", "The ID of the issue")),
		),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get GitLab client: %w", err)
		}

		namespace, err := requiredParam[string](r, "namespace")
		if err != nil {
			return nil, err
		}

		project, err := requiredParam[string](r, "project")
		if err != nil {
			return nil, err
		}

		issueID, err := requiredParam[float64](r, "id")
		if err != nil {
			return nil, err
		}

		projectID := fmt.Sprintf("%s/%s", namespace, project)
		issue, _, err := client.Issues.GetIssue(projectID, int(issueID))
		if err != nil {
			return nil, fmt.Errorf("failed to get issue: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Title: %s\nState: %s\nAuthor: %s", issue.Title, issue.State, issue.Author.Name),
				},
			},
		}, nil
	}

	return tool, handler
}

// ListIssues returns a tool for listing issues in a project
func ListIssues(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_issues",
			mcp.WithDescription(t("TOOL_LIST_ISSUES_DESCRIPTION", "List issues in a project")),
			mcp.WithString("namespace",
				mcp.Required(),
				mcp.Description(t("PARAM_NAMESPACE_DESCRIPTION", "The namespace of the project")),
			),
			mcp.WithString("project",
				mcp.Required(),
				mcp.Description(t("PARAM_PROJECT_DESCRIPTION", "The name of the project")),
			),
		),
		func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitLab client: %w", err)
			}

			namespace, err := requiredParam[string](r, "namespace")
			if err != nil {
				return nil, err
			}

			project, err := requiredParam[string](r, "project")
			if err != nil {
				return nil, err
			}

			issues, _, err := client.Issues.ListProjectIssues(
				fmt.Sprintf("%s/%s", namespace, project),
				&gitlab.ListProjectIssuesOptions{},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to list issues: %w", err)
			}

			response, err := json.Marshal(issues)
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
}

// SearchIssues returns a tool for searching issues in a project
func SearchIssues(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("search_issues",
			mcp.WithDescription(t("TOOL_SEARCH_ISSUES_DESCRIPTION", "Search for issues in a project")),
			mcp.WithString("namespace",
				mcp.Required(),
				mcp.Description(t("PARAM_NAMESPACE_DESCRIPTION", "The namespace of the project")),
			),
			mcp.WithString("project",
				mcp.Required(),
				mcp.Description(t("PARAM_PROJECT_DESCRIPTION", "The name of the project")),
			),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description(t("PARAM_SEARCH_QUERY_DESCRIPTION", "The search query")),
			),
		),
		func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitLab client: %w", err)
			}

			namespace, err := requiredParam[string](r, "namespace")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			project, err := requiredParam[string](r, "project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			query, err := requiredParam[string](r, "query")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			issues, _, err := client.Issues.ListProjectIssues(
				fmt.Sprintf("%s/%s", namespace, project),
				&gitlab.ListProjectIssuesOptions{
					Search: &query,
				},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to search issues: %w", err)
			}

			response, err := json.Marshal(issues)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(response)), nil
		}
}

// GetIssueComments returns a tool for getting comments on an issue
func GetIssueComments(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_issue_comments",
			mcp.WithDescription(t("TOOL_GET_ISSUE_COMMENTS_DESCRIPTION", "Get comments on an issue")),
			mcp.WithString("namespace",
				mcp.Required(),
				mcp.Description(t("PARAM_NAMESPACE_DESCRIPTION", "The namespace of the project")),
			),
			mcp.WithString("project",
				mcp.Required(),
				mcp.Description(t("PARAM_PROJECT_DESCRIPTION", "The name of the project")),
			),
			mcp.WithNumber("id",
				mcp.Required(),
				mcp.Description(t("PARAM_ISSUE_ID_DESCRIPTION", "The ID of the issue")),
			),
		),
		func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitLab client: %w", err)
			}

			namespace, err := requiredParam[string](r, "namespace")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			project, err := requiredParam[string](r, "project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			id, err := RequiredInt(r, "id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			notes, _, err := client.Notes.ListIssueNotes(
				fmt.Sprintf("%s/%s", namespace, project),
				id,
				&gitlab.ListIssueNotesOptions{},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to get issue comments: %w", err)
			}

			response, err := json.Marshal(notes)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(response)), nil
		}
}

// CreateIssue returns a tool for creating a new issue
func CreateIssue(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("create_issue",
			mcp.WithDescription(t("TOOL_CREATE_ISSUE_DESCRIPTION", "Create a new issue")),
			mcp.WithString("namespace",
				mcp.Required(),
				mcp.Description(t("PARAM_NAMESPACE_DESCRIPTION", "The namespace of the project")),
			),
			mcp.WithString("project",
				mcp.Required(),
				mcp.Description(t("PARAM_PROJECT_DESCRIPTION", "The name of the project")),
			),
			mcp.WithString("title",
				mcp.Required(),
				mcp.Description(t("PARAM_ISSUE_TITLE_DESCRIPTION", "The title of the issue")),
			),
			mcp.WithString("description",
				mcp.Description(t("PARAM_ISSUE_DESCRIPTION_DESCRIPTION", "The description of the issue")),
			),
		),
		func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitLab client: %w", err)
			}

			namespace, err := requiredParam[string](r, "namespace")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			project, err := requiredParam[string](r, "project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			title, err := requiredParam[string](r, "title")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			description, _ := requiredParam[string](r, "description")

			issue, _, err := client.Issues.CreateIssue(
				fmt.Sprintf("%s/%s", namespace, project),
				&gitlab.CreateIssueOptions{
					Title:       &title,
					Description: &description,
				},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create issue: %w", err)
			}

			response, err := json.Marshal(issue)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(response)), nil
		}
}

// AddIssueComment returns a tool for adding a comment to an issue
func AddIssueComment(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("add_issue_comment",
			mcp.WithDescription(t("TOOL_ADD_ISSUE_COMMENT_DESCRIPTION", "Add a comment to an issue")),
			mcp.WithString("namespace",
				mcp.Required(),
				mcp.Description(t("PARAM_NAMESPACE_DESCRIPTION", "The namespace of the project")),
			),
			mcp.WithString("project",
				mcp.Required(),
				mcp.Description(t("PARAM_PROJECT_DESCRIPTION", "The name of the project")),
			),
			mcp.WithNumber("id",
				mcp.Required(),
				mcp.Description(t("PARAM_ISSUE_ID_DESCRIPTION", "The ID of the issue")),
			),
			mcp.WithString("body",
				mcp.Required(),
				mcp.Description(t("PARAM_COMMENT_BODY_DESCRIPTION", "The comment text")),
			),
		),
		func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitLab client: %w", err)
			}

			namespace, err := requiredParam[string](r, "namespace")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			project, err := requiredParam[string](r, "project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			id, err := RequiredInt(r, "id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			body, err := requiredParam[string](r, "body")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			note, _, err := client.Notes.CreateIssueNote(
				fmt.Sprintf("%s/%s", namespace, project),
				id,
				&gitlab.CreateIssueNoteOptions{
					Body: &body,
				},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to add issue comment: %w", err)
			}

			response, err := json.Marshal(note)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(response)), nil
		}
}

// UpdateIssue returns a tool for updating an issue
func UpdateIssue(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("update_issue",
			mcp.WithDescription(t("TOOL_UPDATE_ISSUE_DESCRIPTION", "Update an issue")),
			mcp.WithString("namespace",
				mcp.Required(),
				mcp.Description(t("PARAM_NAMESPACE_DESCRIPTION", "The namespace of the project")),
			),
			mcp.WithString("project",
				mcp.Required(),
				mcp.Description(t("PARAM_PROJECT_DESCRIPTION", "The name of the project")),
			),
			mcp.WithNumber("id",
				mcp.Required(),
				mcp.Description(t("PARAM_ISSUE_ID_DESCRIPTION", "The ID of the issue")),
			),
			mcp.WithString("title",
				mcp.Description(t("PARAM_ISSUE_TITLE_DESCRIPTION", "The new title of the issue")),
			),
			mcp.WithString("description",
				mcp.Description(t("PARAM_ISSUE_DESCRIPTION_DESCRIPTION", "The new description of the issue")),
			),
			mcp.WithString("state_event",
				mcp.Description(t("PARAM_ISSUE_STATE_DESCRIPTION", "The new state of the issue (close/reopen)")),
			),
		),
		func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitLab client: %w", err)
			}

			namespace, err := requiredParam[string](r, "namespace")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			project, err := requiredParam[string](r, "project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			id, err := RequiredInt(r, "id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			opts := &gitlab.UpdateIssueOptions{}

			if title, err := requiredParam[string](r, "title"); err == nil {
				opts.Title = &title
			}

			if description, err := requiredParam[string](r, "description"); err == nil {
				opts.Description = &description
			}

			if state, err := requiredParam[string](r, "state_event"); err == nil {
				opts.StateEvent = &state
			}

			issue, _, err := client.Issues.UpdateIssue(
				fmt.Sprintf("%s/%s", namespace, project),
				id,
				opts,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to update issue: %w", err)
			}

			response, err := json.Marshal(issue)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(response)), nil
		}
}
