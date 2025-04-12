# GitLab Integration Authentication Guide

## Overview
This guide explains how to set up authentication for the GitLab integration with the MCP server.

## Authentication Methods

### Personal Access Token
The primary method of authentication is using a GitLab Personal Access Token (PAT).

#### Creating a Personal Access Token
1. Log in to your GitLab instance
2. Go to User Settings > Access Tokens
3. Create a new token with the following scopes:
   - `api` - Full API access
   - `read_api` - Read-only API access
   - `read_repository` - Repository read access
   - `write_repository` - Repository write access

#### Token Scopes
The required token scopes depend on your use case:

##### Read-Only Operations
- `read_api`
- `read_repository`

##### Read-Write Operations
- `api`
- `read_repository`
- `write_repository`

### OAuth2 Authentication
The integration also supports OAuth2 authentication for more complex scenarios.

#### Setting Up OAuth2
1. Register your application in GitLab:
   - Go to User Settings > Applications
   - Create a new application
   - Set the redirect URI
   - Note the Application ID and Secret

2. Configure the integration to use OAuth2:
```yaml
gitlab:
  auth:
    type: "oauth2"
    client_id: "your-application-id"
    client_secret: "your-application-secret"
    redirect_uri: "your-redirect-uri"
```

## Configuration

### Environment Variables
```bash
# Personal Access Token
export GITLAB_TOKEN="your-gitlab-access-token"

# OAuth2
export GITLAB_CLIENT_ID="your-application-id"
export GITLAB_CLIENT_SECRET="your-application-secret"
export GITLAB_REDIRECT_URI="your-redirect-uri"
```

### Configuration File
```yaml
gitlab:
  auth:
    # Personal Access Token
    token: "your-gitlab-access-token"
    
    # OAuth2
    type: "oauth2"
    client_id: "your-application-id"
    client_secret: "your-application-secret"
    redirect_uri: "your-redirect-uri"
```

## Security Best Practices

### Token Security
1. **Storage**
   - Never commit tokens to version control
   - Use environment variables or secure secret management
   - Store tokens in encrypted form

2. **Rotation**
   - Rotate tokens regularly
   - Set token expiration
   - Monitor token usage

3. **Access Control**
   - Use the principle of least privilege
   - Limit token scopes to required permissions
   - Monitor token access

### OAuth2 Security
1. **Client Security**
   - Keep client secrets secure
   - Use HTTPS for redirect URIs
   - Validate state parameter

2. **Token Management**
   - Store refresh tokens securely
   - Implement token refresh logic
   - Handle token expiration

## Troubleshooting

### Common Authentication Issues

1. **Invalid Token**
   - Verify token is correct
   - Check token expiration
   - Ensure token has required scopes

2. **Permission Denied**
   - Verify token scopes
   - Check user permissions
   - Ensure project access

3. **OAuth2 Errors**
   - Verify client credentials
   - Check redirect URI
   - Validate state parameter

### Debugging
Enable debug logging to troubleshoot authentication issues:
```yaml
gitlab:
  logging:
    level: "debug"
```

## Examples

### Personal Access Token Setup
```go
import (
    "context"
    "github.com/jbendotnet/gitlab-mcp-server/pkg/gitlab"
    gitlab "gitlab.com/gitlab-org/api/client-go"
)

func getClient(ctx context.Context) (*gitlab.Client, error) {
    return gitlab.NewClient("your-gitlab-access-token", gitlab.WithBaseURL("https://gitlab.example.com"))
}
```

### OAuth2 Setup
```go
import (
    "context"
    "github.com/jbendotnet/gitlab-mcp-server/pkg/gitlab"
    gitlab "gitlab.com/gitlab-org/api/client-go"
)

func getClient(ctx context.Context) (*gitlab.Client, error) {
    return gitlab.NewOAuthClient(
        "your-application-id",
        "your-application-secret",
        "your-redirect-uri",
        gitlab.WithBaseURL("https://gitlab.example.com"),
    )
}
```

## Additional Resources
- [API Documentation](./api.md)
- [Configuration Guide](./configuration.md)
- [Usage Examples](./usage_examples.md)
- [GitLab OAuth2 Documentation](https://docs.gitlab.com/ee/api/oauth2.html) 