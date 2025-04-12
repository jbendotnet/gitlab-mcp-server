# GitLab Integration API Documentation

## Overview
This document describes the API endpoints and operations available through the GitLab MCP server integration.

## Authentication
The GitLab integration supports token-based authentication. All API requests require a valid GitLab access token.

## Resource Templates

### Repository Content
```
repo://{namespace}/{project}/contents{/path*}
```
- **Description**: Access repository content
- **Methods**: GET
- **Parameters**:
  - `namespace`: GitLab namespace/group
  - `project`: Project name
  - `path`: Optional path within the repository

### Repository Branch Content
```
repo://{namespace}/{project}/refs/heads/{branch}/contents{/path*}
```
- **Description**: Access content from a specific branch
- **Methods**: GET
- **Parameters**:
  - `namespace`: GitLab namespace/group
  - `project`: Project name
  - `branch`: Branch name
  - `path`: Optional path within the repository

### Repository Commit Content
```
repo://{namespace}/{project}/sha/{sha}/contents{/path*}
```
- **Description**: Access content at a specific commit
- **Methods**: GET
- **Parameters**:
  - `namespace`: GitLab namespace/group
  - `project`: Project name
  - `sha`: Commit SHA
  - `path`: Optional path within the repository

### Repository Tag Content
```
repo://{namespace}/{project}/refs/tags/{tag}/contents{/path*}
```
- **Description**: Access content at a specific tag
- **Methods**: GET
- **Parameters**:
  - `namespace`: GitLab namespace/group
  - `project`: Project name
  - `tag`: Tag name
  - `path`: Optional path within the repository

### Merge Request Content
```
repo://{namespace}/{project}/merge_requests/{id}
```
- **Description**: Access merge request content
- **Methods**: GET
- **Parameters**:
  - `namespace`: GitLab namespace/group
  - `project`: Project name
  - `id`: Merge request ID

## Tools

### Repository Operations

#### Get Repository
- **Tool Name**: `get_repository`
- **Description**: Get information about a specific repository
- **Parameters**:
  - `namespace`: GitLab namespace/group
  - `project`: Project name

#### List Repositories
- **Tool Name**: `list_repositories`
- **Description**: List repositories accessible to the authenticated user
- **Parameters**: None

#### Search Repositories
- **Tool Name**: `search_repositories`
- **Description**: Search for repositories
- **Parameters**:
  - `query`: Search query string

### Merge Request Operations

#### Get Merge Request
- **Tool Name**: `get_merge_request`
- **Description**: Get information about a specific merge request
- **Parameters**:
  - `namespace`: GitLab namespace/group
  - `project`: Project name
  - `id`: Merge request ID

#### List Merge Requests
- **Tool Name**: `list_merge_requests`
- **Description**: List merge requests for a repository
- **Parameters**:
  - `namespace`: GitLab namespace/group
  - `project`: Project name

#### Get Merge Request Comments
- **Tool Name**: `get_merge_request_comments`
- **Description**: Get comments for a merge request
- **Parameters**:
  - `namespace`: GitLab namespace/group
  - `project`: Project name
  - `id`: Merge request ID

#### Create Merge Request (Read-Write Mode)
- **Tool Name**: `create_merge_request`
- **Description**: Create a new merge request
- **Parameters**:
  - `namespace`: GitLab namespace/group
  - `project`: Project name
  - `title`: Merge request title
  - `description`: Merge request description
  - `source_branch`: Source branch
  - `target_branch`: Target branch

### Issue Operations

#### Get Issue
- **Tool Name**: `get_issue`
- **Description**: Get information about a specific issue
- **Parameters**:
  - `namespace`: GitLab namespace/group
  - `project`: Project name
  - `id`: Issue ID

#### List Issues
- **Tool Name**: `list_issues`
- **Description**: List issues for a repository
- **Parameters**:
  - `namespace`: GitLab namespace/group
  - `project`: Project name

#### Search Issues
- **Tool Name**: `search_issues`
- **Description**: Search for issues
- **Parameters**:
  - `query`: Search query string

#### Get Issue Comments
- **Tool Name**: `get_issue_comments`
- **Description**: Get comments for an issue
- **Parameters**:
  - `namespace`: GitLab namespace/group
  - `project`: Project name
  - `id`: Issue ID

### Search Operations

#### Search Projects
- **Tool Name**: `search_projects`
- **Description**: Search for projects
- **Parameters**:
  - `query`: Search query string

#### Search Users
- **Tool Name**: `search_users`
- **Description**: Search for users
- **Parameters**:
  - `query`: Search query string

### User Operations

#### Get Current User
- **Tool Name**: `get_me`
- **Description**: Get information about the authenticated user
- **Parameters**: None

## Error Handling
The API returns standard HTTP status codes and includes error messages in the response body when operations fail.

Common error codes:
- `400`: Bad Request
- `401`: Unauthorized
- `403`: Forbidden
- `404`: Not Found
- `429`: Too Many Requests
- `500`: Internal Server Error

## Rate Limiting
The GitLab API has rate limits that are enforced. The integration includes handling for rate limits and will return appropriate error messages when limits are reached.

## Read-Only Mode
When the server is configured in read-only mode, write operations (create, update, delete) will be disabled. 