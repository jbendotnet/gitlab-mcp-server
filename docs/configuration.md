# GitLab Integration Configuration Guide

## Overview
This guide explains how to configure the GitLab integration for the MCP server.

## Prerequisites
- GitLab instance URL
- GitLab access token with appropriate permissions
- MCP server installation

## Configuration Options

### Required Configuration

#### GitLab Instance URL
```yaml
gitlab:
  base_url: "https://gitlab.example.com"  # Your GitLab instance URL
```

#### Authentication Token
```yaml
gitlab:
  token: "your-gitlab-access-token"  # GitLab personal access token
```

### Optional Configuration

#### Read-Only Mode
```yaml
gitlab:
  read_only: true  # Set to true to disable write operations
```

#### Rate Limiting
```yaml
gitlab:
  rate_limit:
    requests_per_second: 10  # Maximum requests per second
    burst_size: 20          # Maximum burst size
```

#### Logging
```yaml
gitlab:
  logging:
    level: "info"           # Log level (debug, info, warn, error)
    format: "json"          # Log format (text, json)
```

## Authentication Setup

### Creating a GitLab Access Token
1. Log in to your GitLab instance
2. Go to User Settings > Access Tokens
3. Create a new token with the following scopes:
   - `api` - Full API access
   - `read_api` - Read-only API access (if using read-only mode)
   - `read_repository` - Repository read access
   - `write_repository` - Repository write access (if not in read-only mode)

### Token Permissions
The required token permissions depend on your use case:

#### Read-Only Operations
- `read_api`
- `read_repository`

#### Read-Write Operations
- `api`
- `read_repository`
- `write_repository`

## Environment Variables

You can also configure the integration using environment variables:

```bash
export GITLAB_BASE_URL="https://gitlab.example.com"
export GITLAB_TOKEN="your-gitlab-access-token"
export GITLAB_READ_ONLY="true"
export GITLAB_RATE_LIMIT_REQUESTS_PER_SECOND="10"
export GITLAB_RATE_LIMIT_BURST_SIZE="20"
export GITLAB_LOG_LEVEL="info"
export GITLAB_LOG_FORMAT="json"
```

## Configuration File Example

Here's a complete configuration file example:

```yaml
gitlab:
  base_url: "https://gitlab.example.com"
  token: "your-gitlab-access-token"
  read_only: false
  rate_limit:
    requests_per_second: 10
    burst_size: 20
  logging:
    level: "info"
    format: "json"
```

## Security Considerations

1. **Token Security**
   - Never commit access tokens to version control
   - Use environment variables or secure secret management
   - Rotate tokens regularly

2. **Rate Limiting**
   - Configure appropriate rate limits based on your GitLab instance's limits
   - Monitor rate limit usage

3. **Access Control**
   - Use read-only mode when write access is not required
   - Limit token permissions to the minimum required

## Troubleshooting

### Common Issues

1. **Authentication Errors**
   - Verify token is valid and has correct permissions
   - Check token expiration
   - Ensure base URL is correct

2. **Rate Limit Errors**
   - Adjust rate limit configuration
   - Implement retry logic
   - Monitor usage patterns

3. **Permission Errors**
   - Verify token scopes
   - Check project access permissions
   - Ensure user has required access

### Logging
Enable debug logging to troubleshoot issues:

```yaml
gitlab:
  logging:
    level: "debug"
```

## Support
For additional support:
- Check the [API Documentation](./api.md)
- Review the [Implementation Plan](./implementation_plan.md)
- Open an issue in the project repository 