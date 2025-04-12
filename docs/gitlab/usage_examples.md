# GitLab Integration Usage Examples

## Overview
This document provides practical examples of using the GitLab integration with the MCP server.

## Basic Setup

### Initialize the GitLab Client
```go
import (
    "context"
    "github.com/jbendotnet/gitlab-mcp-server/pkg/gitlab"
    gitlab "gitlab.com/gitlab-org/api/client-go"
)

func getClient(ctx context.Context) (*gitlab.Client, error) {
    return gitlab.NewClient("your-gitlab-access-token", gitlab.WithBaseURL("https://gitlab.example.com"))
}

// Create the server
server := gitlab.NewServer(
    getClient,
    "1.0.0",
    false, // read-only mode
    translations.DefaultTranslationHelper,
)
```

## Repository Operations

### Get Repository Content
```go
// Access repository content
content, err := server.GetResource(ctx, "repo://namespace/project/contents/path/to/file")
```

### List Repositories
```go
// List all accessible repositories
repos, err := server.CallTool(ctx, "list_repositories", nil)
```

### Search Repositories
```go
// Search for repositories
result, err := server.CallTool(ctx, "search_repositories", map[string]interface{}{
    "query": "search term",
})
```

## Merge Request Operations

### Get Merge Request
```go
// Get a specific merge request
mr, err := server.CallTool(ctx, "get_merge_request", map[string]interface{}{
    "namespace": "namespace",
    "project": "project",
    "id": 123,
})
```

### Create Merge Request
```go
// Create a new merge request
result, err := server.CallTool(ctx, "create_merge_request", map[string]interface{}{
    "namespace": "namespace",
    "project": "project",
    "title": "New feature",
    "description": "Description of changes",
    "source_branch": "feature-branch",
    "target_branch": "main",
})
```

### List Merge Request Comments
```go
// Get comments for a merge request
comments, err := server.CallTool(ctx, "get_merge_request_comments", map[string]interface{}{
    "namespace": "namespace",
    "project": "project",
    "id": 123,
})
```

## Issue Operations

### Get Issue
```go
// Get a specific issue
issue, err := server.CallTool(ctx, "get_issue", map[string]interface{}{
    "namespace": "namespace",
    "project": "project",
    "id": 456,
})
```

### Create Issue
```go
// Create a new issue
result, err := server.CallTool(ctx, "create_issue", map[string]interface{}{
    "namespace": "namespace",
    "project": "project",
    "title": "Bug report",
    "description": "Description of the bug",
})
```

### Search Issues
```go
// Search for issues
result, err := server.CallTool(ctx, "search_issues", map[string]interface{}{
    "query": "bug OR error",
})
```

## Search Operations

### Search Projects
```go
// Search for projects
result, err := server.CallTool(ctx, "search_projects", map[string]interface{}{
    "query": "project name",
})
```

### Search Users
```go
// Search for users
result, err := server.CallTool(ctx, "search_users", map[string]interface{}{
    "query": "username",
})
```

## User Operations

### Get Current User
```go
// Get information about the authenticated user
user, err := server.CallTool(ctx, "get_me", nil)
```

## Error Handling

### Basic Error Handling
```go
result, err := server.CallTool(ctx, "get_repository", map[string]interface{}{
    "namespace": "namespace",
    "project": "project",
})
if err != nil {
    switch {
    case strings.Contains(err.Error(), "404"):
        // Handle not found
    case strings.Contains(err.Error(), "401"):
        // Handle unauthorized
    case strings.Contains(err.Error(), "403"):
        // Handle forbidden
    case strings.Contains(err.Error(), "429"):
        // Handle rate limit
    default:
        // Handle other errors
    }
}
```

### Rate Limit Handling
```go
func callWithRetry(ctx context.Context, server *server.MCPServer, tool string, params map[string]interface{}) (*mcp.CallToolResult, error) {
    var result *mcp.CallToolResult
    var err error
    
    for i := 0; i < 3; i++ {
        result, err = server.CallTool(ctx, tool, params)
        if err == nil {
            return result, nil
        }
        
        if strings.Contains(err.Error(), "429") {
            // Wait before retrying
            time.Sleep(time.Second * time.Duration(i+1))
            continue
        }
        
        return nil, err
    }
    
    return nil, err
}
```

## Advanced Examples

### Batch Operations
```go
func batchCreateIssues(ctx context.Context, server *server.MCPServer, issues []Issue) error {
    for _, issue := range issues {
        _, err := server.CallTool(ctx, "create_issue", map[string]interface{}{
            "namespace": issue.Namespace,
            "project": issue.Project,
            "title": issue.Title,
            "description": issue.Description,
        })
        if err != nil {
            return fmt.Errorf("failed to create issue %s: %w", issue.Title, err)
        }
    }
    return nil
}
```

### Monitoring Merge Requests
```go
func monitorMergeRequests(ctx context.Context, server *server.MCPServer, namespace, project string) error {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            mrs, err := server.CallTool(ctx, "list_merge_requests", map[string]interface{}{
                "namespace": namespace,
                "project": project,
            })
            if err != nil {
                log.Printf("Error listing merge requests: %v", err)
                continue
            }
            
            // Process merge requests
            processMergeRequests(mrs)
        }
    }
}
```

## Best Practices

1. **Error Handling**
   - Always check for errors
   - Implement retry logic for rate limits
   - Log errors appropriately

2. **Rate Limiting**
   - Implement backoff strategies
   - Monitor rate limit usage
   - Cache responses when appropriate

3. **Resource Management**
   - Close resources properly
   - Use context for cancellation
   - Implement timeouts

4. **Security**
   - Never log sensitive information
   - Validate input parameters
   - Use appropriate token permissions

## Additional Resources
- [API Documentation](./api.md)
- [Configuration Guide](./configuration.md)
- [Implementation Plan](./implementation_plan.md) 