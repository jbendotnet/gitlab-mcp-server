package gitlab

import (
	"context"
	"fmt"

	"github.com/jbendotnet/gitlab-mcp-server/pkg/translations"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// GetUserProfile implements the get user profile tool
func GetUserProfile(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"get_user_profile",
		mcp.WithDescription(t("TOOL_GET_USER_PROFILE_DESCRIPTION", "Get a GitLab user's profile")),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get GitLab client: %w", err)
		}

		username, err := RequiredString(r, "username")
		if err != nil {
			return nil, err
		}

		// First get the user ID by searching for the username
		users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{
			Username: &username,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get user profile: %w", err)
		}

		if len(users) == 0 {
			return nil, fmt.Errorf("user not found")
		}

		user := users[0]

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Name: %s\nUsername: %s\nEmail: %s\nBio: %s", user.Name, user.Username, user.Email, user.Bio),
				},
			},
		}, nil
	}

	return tool, handler
}

// ListUserGroups implements the list user groups tool
func ListUserGroups(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"list_user_groups",
		mcp.WithDescription(t("TOOL_LIST_USER_GROUPS_DESCRIPTION", "List groups a user belongs to")),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get GitLab client: %w", err)
		}

		username, ok := r.Params.Arguments["username"].(string)
		if !ok {
			return nil, fmt.Errorf("username parameter is required")
		}

		// First get the user ID by searching for the username
		users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{
			Username: &username,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to find user: %w", err)
		}

		if len(users) == 0 {
			return nil, fmt.Errorf("user not found")
		}

		userID := users[0].ID

		// Then get the user's groups
		groups, _, err := client.Groups.ListGroups(&gitlab.ListGroupsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
			},
			AllAvailable: gitlab.Bool(true),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list groups: %w", err)
		}

		// Filter groups where the user is a member
		userGroups := make([]*gitlab.Group, 0)
		for _, group := range groups {
			members, _, err := client.Groups.ListGroupMembers(group.ID, &gitlab.ListGroupMembersOptions{
				ListOptions: gitlab.ListOptions{
					PerPage: 100,
				},
			})
			if err != nil {
				continue
			}
			for _, member := range members {
				if member.ID == userID {
					userGroups = append(userGroups, group)
					break
				}
			}
		}

		groupList := "Groups:\n"
		for _, group := range userGroups {
			groupList += fmt.Sprintf("- %s (ID: %d)\n", group.Name, group.ID)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: groupList,
				},
			},
		}, nil
	}

	return tool, handler
}

// GetUserPermissions implements the get user permissions tool
func GetUserPermissions(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"get_user_permissions",
		mcp.WithDescription(t("TOOL_GET_USER_PERMISSIONS_DESCRIPTION", "Get a user's permissions in a project")),
	)

	handler = func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get GitLab client: %w", err)
		}

		projectID, ok := r.Params.Arguments["project_id"].(string)
		if !ok {
			return nil, fmt.Errorf("project_id parameter is required")
		}

		username, ok := r.Params.Arguments["username"].(string)
		if !ok {
			return nil, fmt.Errorf("username parameter is required")
		}

		// First get the user ID by searching for the username
		users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{
			Username: &username,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to find user: %w", err)
		}

		if len(users) == 0 {
			return nil, fmt.Errorf("user not found")
		}

		userID := users[0].ID

		// Get project members to find the user's access level
		members, _, err := client.ProjectMembers.ListProjectMembers(projectID, &gitlab.ListProjectMembersOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get project members: %w", err)
		}

		var userMember *gitlab.ProjectMember
		for _, member := range members {
			if member.ID == userID {
				userMember = member
				break
			}
		}

		if userMember == nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("User %s has no permissions in project %s", username, projectID),
					},
				},
			}, nil
		}

		accessLevel := userMember.AccessLevel
		permission := "Unknown"
		switch accessLevel {
		case gitlab.GuestPermissions:
			permission = "Guest"
		case gitlab.ReporterPermissions:
			permission = "Reporter"
		case gitlab.DeveloperPermissions:
			permission = "Developer"
		case gitlab.MaintainerPermissions:
			permission = "Maintainer"
		case gitlab.OwnerPermissions:
			permission = "Owner"
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("User %s has %s permissions in project %s", username, permission, projectID),
				},
			},
		}, nil
	}

	return tool, handler
}
