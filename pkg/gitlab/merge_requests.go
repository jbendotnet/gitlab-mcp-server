package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/jbendotnet/gitlab-mcp-server/pkg/translations"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// GetMergeRequest returns a tool for getting a specific merge request
func GetMergeRequest(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"get_merge_request",
		mcp.WithDescription(t("TOOL_GET_MERGE_REQUEST_DESCRIPTION", "Get a specific merge request")),
		mcp.WithString("namespace",
			mcp.Required(),
			mcp.Description(t("PARAM_NAMESPACE_DESCRIPTION", "The namespace of the project")),
		),
		mcp.WithString("project",
			mcp.Required(),
			mcp.Description(t("PARAM_PROJECT_DESCRIPTION", "The name of the project")),
		),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description(t("PARAM_MERGE_REQUEST_ID_DESCRIPTION", "The ID of the merge request")),
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
		project, err := RequiredString(r, "project")
		if err != nil {
			return nil, err
		}
		id, err := RequiredInt(r, "id")
		if err != nil {
			return nil, err
		}

		mr, _, err := client.MergeRequests.GetMergeRequest(
			fmt.Sprintf("%s/%s", namespace, project),
			id,
			nil, // No options needed for basic get
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get merge request: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Title: %s\nDescription: %s\nState: %s\nAuthor: %s", mr.Title, mr.Description, mr.State, mr.Author.Name),
				},
			},
		}, nil
	}

	return tool, handler
}

// ListMergeRequests returns a tool for listing merge requests
func ListMergeRequests(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"list_merge_requests",
		mcp.WithDescription(t("TOOL_LIST_MERGE_REQUESTS_DESCRIPTION", "List merge requests in a project")),
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

		namespace, err := requiredParam[string](r, "namespace")
		if err != nil {
			return nil, err
		}

		project, err := requiredParam[string](r, "project")
		if err != nil {
			return nil, err
		}

		opts := &gitlab.ListProjectMergeRequestsOptions{}

		if state, ok := r.Params.Arguments["state"].(string); ok {
			opts.State = gitlab.Ptr(state)
		}

		if orderBy, ok := r.Params.Arguments["order_by"].(string); ok {
			opts.OrderBy = gitlab.Ptr(orderBy)
		}

		if sort, ok := r.Params.Arguments["sort"].(string); ok {
			opts.Sort = gitlab.Ptr(sort)
		}

		mrs, _, err := client.MergeRequests.ListProjectMergeRequests(
			fmt.Sprintf("%s/%s", namespace, project),
			opts,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to list merge requests: %w", err)
		}

		response, err := json.Marshal(mrs)
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

// GetMergeRequestComments returns a tool for getting merge request comments
func GetMergeRequestComments(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"get_merge_request_comments",
		mcp.WithDescription(t("TOOL_GET_MERGE_REQUEST_COMMENTS_DESCRIPTION", "Get comments for a merge request")),
		mcp.WithString("namespace",
			mcp.Required(),
			mcp.Description(t("PARAM_NAMESPACE_DESCRIPTION", "The namespace of the project")),
		),
		mcp.WithString("project",
			mcp.Required(),
			mcp.Description(t("PARAM_PROJECT_DESCRIPTION", "The name of the project")),
		),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description(t("PARAM_MERGE_REQUEST_ID_DESCRIPTION", "The ID of the merge request")),
		),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to get GitLab client: %w", err).Error()), nil
		}

		namespace, err := requiredParam[string](r, "namespace")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		project, err := requiredParam[string](r, "project")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		id, err := requiredParam[string](r, "id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		mrID, err := strconv.Atoi(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("invalid merge request ID: %w", err).Error()), nil
		}

		notes, resp, err := client.Notes.ListMergeRequestNotes(
			fmt.Sprintf("%s/%s", namespace, project),
			mrID,
			&gitlab.ListMergeRequestNotesOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to get merge request comments: %w", err).Error()), nil
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return mcp.NewToolResultError(fmt.Errorf("failed to read response body: %w", err).Error()), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("failed to get merge request comments: %s", string(body))), nil
		}

		jsonData, err := json.Marshal(notes)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to marshal response: %w", err).Error()), nil
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	}

	return tool, handler
}

// CreateMergeRequest returns a tool for creating a new merge request
func CreateMergeRequest(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"create_merge_request",
		mcp.WithDescription(t("TOOL_CREATE_MERGE_REQUEST_DESCRIPTION", "Create a new merge request")),
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
			mcp.Description(t("PARAM_MERGE_REQUEST_TITLE_DESCRIPTION", "The title of the merge request")),
		),
		mcp.WithString("description",
			mcp.Required(),
			mcp.Description(t("PARAM_MERGE_REQUEST_DESCRIPTION_DESCRIPTION", "The description of the merge request")),
		),
		mcp.WithString("source_branch",
			mcp.Required(),
			mcp.Description(t("PARAM_SOURCE_BRANCH_DESCRIPTION", "The source branch")),
		),
		mcp.WithString("target_branch",
			mcp.Required(),
			mcp.Description(t("PARAM_TARGET_BRANCH_DESCRIPTION", "The target branch")),
		),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to get GitLab client: %w", err).Error()), nil
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
		description, err := requiredParam[string](r, "description")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		sourceBranch, err := requiredParam[string](r, "source_branch")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		targetBranch, err := requiredParam[string](r, "target_branch")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		mr, _, err := client.MergeRequests.CreateMergeRequest(
			fmt.Sprintf("%s/%s", namespace, project),
			&gitlab.CreateMergeRequestOptions{
				Title:        &title,
				Description:  &description,
				SourceBranch: &sourceBranch,
				TargetBranch: &targetBranch,
			},
		)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to create merge request: %w", err).Error()), nil
		}

		jsonData, err := json.Marshal(mr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to marshal response: %w", err).Error()), nil
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	}

	return tool, handler
}

// AddMergeRequestComment returns a tool for adding a comment to a merge request
func AddMergeRequestComment(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"add_merge_request_comment",
		mcp.WithDescription(t("TOOL_ADD_MERGE_REQUEST_COMMENT_DESCRIPTION", "Add a comment to a merge request")),
		mcp.WithString("namespace",
			mcp.Required(),
			mcp.Description(t("PARAM_NAMESPACE_DESCRIPTION", "The namespace of the project")),
		),
		mcp.WithString("project",
			mcp.Required(),
			mcp.Description(t("PARAM_PROJECT_DESCRIPTION", "The name of the project")),
		),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description(t("PARAM_MERGE_REQUEST_ID_DESCRIPTION", "The ID of the merge request")),
		),
		mcp.WithString("body",
			mcp.Required(),
			mcp.Description(t("PARAM_COMMENT_BODY_DESCRIPTION", "The body of the comment")),
		),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to get GitLab client: %w", err).Error()), nil
		}

		namespace, err := requiredParam[string](r, "namespace")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		project, err := requiredParam[string](r, "project")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		id, err := requiredParam[string](r, "id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		body, err := requiredParam[string](r, "body")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		mrID, err := strconv.Atoi(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("invalid merge request ID: %w", err).Error()), nil
		}

		note, resp, err := client.Notes.CreateMergeRequestNote(
			fmt.Sprintf("%s/%s", namespace, project),
			mrID,
			&gitlab.CreateMergeRequestNoteOptions{
				Body: &body,
			},
		)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to add comment: %w", err).Error()), nil
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusCreated {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return mcp.NewToolResultError(fmt.Errorf("failed to read response body: %w", err).Error()), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("failed to add comment: %s", string(body))), nil
		}

		jsonData, err := json.Marshal(note)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to marshal response: %w", err).Error()), nil
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	}

	return tool, handler
}

// UpdateMergeRequest returns a tool for updating a merge request
func UpdateMergeRequest(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"update_merge_request",
		mcp.WithDescription(t("TOOL_UPDATE_MERGE_REQUEST_DESCRIPTION", "Update a merge request")),
		mcp.WithString("namespace",
			mcp.Required(),
			mcp.Description(t("PARAM_NAMESPACE_DESCRIPTION", "The namespace of the project")),
		),
		mcp.WithString("project",
			mcp.Required(),
			mcp.Description(t("PARAM_PROJECT_DESCRIPTION", "The name of the project")),
		),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description(t("PARAM_MERGE_REQUEST_ID_DESCRIPTION", "The ID of the merge request")),
		),
		mcp.WithString("title",
			mcp.Description(t("PARAM_MERGE_REQUEST_TITLE_DESCRIPTION", "The new title of the merge request")),
		),
		mcp.WithString("description",
			mcp.Description(t("PARAM_MERGE_REQUEST_DESCRIPTION_DESCRIPTION", "The new description of the merge request")),
		),
		mcp.WithString("state_event",
			mcp.Description(t("PARAM_MERGE_REQUEST_STATE_DESCRIPTION", "The new state of the merge request (opened/closed)")),
		),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to get GitLab client: %w", err).Error()), nil
		}

		namespace, err := requiredParam[string](r, "namespace")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		project, err := requiredParam[string](r, "project")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		id, err := requiredParam[string](r, "id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		mrID, err := strconv.Atoi(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("invalid merge request ID: %w", err).Error()), nil
		}

		opts := &gitlab.UpdateMergeRequestOptions{}

		if title, err := OptionalParam[string](r, "title"); err == nil && title != "" {
			opts.Title = &title
		}
		if description, err := OptionalParam[string](r, "description"); err == nil && description != "" {
			opts.Description = &description
		}
		if stateEvent, err := OptionalParam[string](r, "state_event"); err == nil && stateEvent != "" {
			opts.StateEvent = &stateEvent
		}

		mr, resp, err := client.MergeRequests.UpdateMergeRequest(
			fmt.Sprintf("%s/%s", namespace, project),
			mrID,
			opts,
		)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to update merge request: %w", err).Error()), nil
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return mcp.NewToolResultError(fmt.Errorf("failed to read response body: %w", err).Error()), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("failed to update merge request: %s", string(body))), nil
		}

		jsonData, err := json.Marshal(mr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to marshal response: %w", err).Error()), nil
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	}

	return tool, handler
}
