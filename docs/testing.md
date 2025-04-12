# GitLab Integration Testing

This document describes how to run and configure the GitLab integration tests.

## Prerequisites

1. A GitLab instance (either GitLab.com or self-hosted)
2. A GitLab personal access token with appropriate permissions
3. Go 1.16 or later installed
4. The `github-mcp-server` repository cloned

## Setting Up Test Environment

### 1. Create a GitLab Personal Access Token

1. Log in to your GitLab instance
2. Go to User Settings → Access Tokens
3. Create a new token with the following scopes:
   - `api` - Full API access
   - `read_api` - Read API access
   - `read_repository` - Read repository access
   - `write_repository` - Write repository access

### 2. Configure Environment Variables

Set the following environment variables:

```bash
# Required for running integration tests
export GITLAB_INTEGRATION_TEST=true

# GitLab API credentials
export GITLAB_TOKEN=your_personal_access_token
export GITLAB_BASE_URL=https://gitlab.com/api/v4  # or your instance URL
```

You can add these to your shell's configuration file (e.g., `~/.bashrc`, `~/.zshrc`, or `~/.config/fish/config.fish`) to make them permanent.

## Running Tests

### Running All Tests

To run all tests, including integration tests:

```bash
go test ./pkg/gitlab/... -v
```

### Running Only Integration Tests

To run only the integration tests:

```bash
go test ./pkg/gitlab/... -v -run TestIntegrationGitLab
```

### Running Specific Test Cases

You can run specific test cases using the `-run` flag with a regex pattern:

```bash
# Run only repository tests
go test ./pkg/gitlab/... -v -run TestIntegrationGitLab/Repository

# Run only issue tests
go test ./pkg/gitlab/... -v -run TestIntegrationGitLab/Issue

# Run only merge request tests
go test ./pkg/gitlab/... -v -run TestIntegrationGitLab/Merge
```

## Test Coverage

To generate a coverage report:

```bash
go test ./pkg/gitlab/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Cleanup

The integration tests automatically clean up after themselves by:
1. Creating test projects with unique names
2. Deleting the projects after the tests complete
3. Using `defer` statements to ensure cleanup even if tests fail

## Troubleshooting

### Common Issues

1. **Authentication Errors**
   - Verify your `GITLAB_TOKEN` is correct and has the required scopes
   - Check if the token has expired

2. **API Rate Limiting**
   - GitLab.com has rate limits on API calls
   - If you hit rate limits, consider using a self-hosted instance for testing

3. **Project Creation Failures**
   - Verify you have permission to create projects
   - Check if the project name is already in use

4. **Test Timeouts**
   - Some operations (like creating merge requests) may take time
   - If tests timeout, you can increase the timeout with `-timeout` flag:
     ```bash
     go test ./pkg/gitlab/... -v -timeout 10m
     ```

### Debugging

To enable debug logging during tests:

```bash
export GITLAB_DEBUG=true
go test ./pkg/gitlab/... -v
```

## Best Practices

1. **Use Unique Names**
   - Test projects are created with unique names
   - Avoid using common prefixes in your test data

2. **Clean Up Manually**
   - If tests fail, you may need to clean up manually
   - Check your GitLab instance for any leftover test projects

3. **Test in Isolation**
   - Each test case runs in isolation
   - Tests don't depend on each other's state

4. **Use Read-Only Mode**
   - For read-only operations, use the read-only mode:
     ```go
     s := NewServer(getClient, "1.0.0", true, translationHelper)
     ```

## Contributing

When adding new integration tests:
1. Follow the existing test structure
2. Include proper cleanup
3. Add relevant documentation
4. Test both success and error cases
5. Consider rate limits and API quotas 