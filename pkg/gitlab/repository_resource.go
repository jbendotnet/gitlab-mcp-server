package gitlab

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/jbendotnet/gitlab-mcp-server/pkg/translations"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// GetRepositoryResourceContent defines the resource template and handler for getting repository content
func GetRepositoryResourceContent(getClient GetClientFn, t translations.TranslationHelperFunc) (mcp.ResourceTemplate, server.ResourceTemplateHandlerFunc) {
	return mcp.NewResourceTemplate(
			"repo://{namespace}/{project}/contents{/path*}",
			t("RESOURCE_REPOSITORY_CONTENT_DESCRIPTION", "Repository Content"),
		),
		RepositoryResourceContentsHandler(getClient)
}

// GetRepositoryResourceBranchContent defines the resource template and handler for getting repository content for a branch
func GetRepositoryResourceBranchContent(getClient GetClientFn, t translations.TranslationHelperFunc) (mcp.ResourceTemplate, server.ResourceTemplateHandlerFunc) {
	return mcp.NewResourceTemplate(
			"repo://{namespace}/{project}/refs/heads/{branch}/contents{/path*}",
			t("RESOURCE_REPOSITORY_CONTENT_BRANCH_DESCRIPTION", "Repository Content for specific branch"),
		),
		RepositoryResourceContentsHandler(getClient)
}

// GetRepositoryResourceCommitContent defines the resource template and handler for getting repository content for a commit
func GetRepositoryResourceCommitContent(getClient GetClientFn, t translations.TranslationHelperFunc) (mcp.ResourceTemplate, server.ResourceTemplateHandlerFunc) {
	return mcp.NewResourceTemplate(
			"repo://{namespace}/{project}/sha/{sha}/contents{/path*}",
			t("RESOURCE_REPOSITORY_CONTENT_COMMIT_DESCRIPTION", "Repository Content for specific commit"),
		),
		RepositoryResourceContentsHandler(getClient)
}

// GetRepositoryResourceTagContent defines the resource template and handler for getting repository content for a tag
func GetRepositoryResourceTagContent(getClient GetClientFn, t translations.TranslationHelperFunc) (mcp.ResourceTemplate, server.ResourceTemplateHandlerFunc) {
	return mcp.NewResourceTemplate(
			"repo://{namespace}/{project}/refs/tags/{tag}/contents{/path*}",
			t("RESOURCE_REPOSITORY_CONTENT_TAG_DESCRIPTION", "Repository Content for specific tag"),
		),
		RepositoryResourceContentsHandler(getClient)
}

// GetRepositoryResourceMergeRequestContent defines the resource template and handler for getting merge request content
func GetRepositoryResourceMergeRequestContent(getClient GetClientFn, t translations.TranslationHelperFunc) (mcp.ResourceTemplate, server.ResourceTemplateHandlerFunc) {
	return mcp.NewResourceTemplate(
			"repo://{namespace}/{project}/merge_requests/{id}",
			t("RESOURCE_REPOSITORY_MERGE_REQUEST_DESCRIPTION", "Merge Request"),
		),
		RepositoryResourceMergeRequestHandler(getClient)
}

// RepositoryResourceContentsHandler handles repository content requests
func RepositoryResourceContentsHandler(getClient GetClientFn) func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get GitLab client: %w", err)
		}

		namespace := request.Params.Arguments["namespace"].(string)
		project := request.Params.Arguments["project"].(string)
		path := request.Params.Arguments["path"].(string)
		ref := request.Params.Arguments["ref"].(string)

		file, _, err := client.RepositoryFiles.GetFile(
			fmt.Sprintf("%s/%s", namespace, project),
			path,
			&gitlab.GetFileOptions{
				Ref: &ref,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get file: %w", err)
		}

		content, err := base64.StdEncoding.DecodeString(file.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to decode file content: %w", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/plain",
				Text:     string(content),
			},
		}, nil
	}
}

// RepositoryResourceMergeRequestHandler handles merge request requests
func RepositoryResourceMergeRequestHandler(getClient GetClientFn) func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get GitLab client: %w", err)
		}

		namespace := request.Params.Arguments["namespace"].(string)
		project := request.Params.Arguments["project"].(string)
		idStr := request.Params.Arguments["id"].(string)

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid merge request ID: %w", err)
		}

		mr, _, err := client.MergeRequests.GetMergeRequest(
			fmt.Sprintf("%s/%s", namespace, project),
			id,
			nil, // No options needed for basic get
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get merge request: %w", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/plain",
				Text:     fmt.Sprintf("Title: %s\nDescription: %s\nState: %s", mr.Title, mr.Description, mr.State),
			},
		}, nil
	}
}
