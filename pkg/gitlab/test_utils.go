package gitlab

import (
	"io"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// mockSearchService is a mock implementation of the GitLab search service
type mockSearchService struct {
	searchProjectFunc      func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Project, *gitlab.Response, error)
	searchMergeRequestFunc func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequest, *gitlab.Response, error)
	searchUserFunc         func(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.User, *gitlab.Response, error)
}

// ensure mockSearchService implements the gitlab.SearchServiceInterface
var _ gitlab.SearchServiceInterface = &mockSearchService{}

func (m *mockSearchService) Projects(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Project, *gitlab.Response, error) {
	return m.searchProjectFunc(query, opt, options...)
}

func (m *mockSearchService) MergeRequests(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
	return m.searchMergeRequestFunc(query, opt, options...)
}

func (m *mockSearchService) Users(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.User, *gitlab.Response, error) {
	return m.searchUserFunc(query, opt, options...)
}

func (m *mockSearchService) ProjectsByGroup(gid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Project, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) Issues(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Issue, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) IssuesByGroup(gid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Issue, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) IssuesByProject(pid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Issue, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) MergeRequestsByGroup(gid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) MergeRequestsByProject(pid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) Milestones(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Milestone, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) MilestonesByGroup(gid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Milestone, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) MilestonesByProject(pid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Milestone, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) Blobs(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Blob, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) BlobsByGroup(gid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Blob, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) BlobsByProject(pid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Blob, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) Commits(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Commit, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) CommitsByGroup(gid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Commit, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) CommitsByProject(pid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Commit, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) NotesByProject(pid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Note, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) SnippetBlobs(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Snippet, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) SnippetBlobsByProject(pid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Snippet, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) WikiBlobs(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Wiki, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) WikiBlobsByProject(pid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Wiki, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) WikiBlobsByGroup(gid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Wiki, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) SnippetTitles(query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Snippet, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) UsersByGroup(gid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.User, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockSearchService) UsersByProject(pid interface{}, query string, opt *gitlab.SearchOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.User, *gitlab.Response, error) {
	return nil, nil, nil
}

// mockUsersService is a mock implementation of the GitLab users service
type mockUsersService struct {
	currentUserFunc func(options ...gitlab.RequestOptionFunc) (*gitlab.User, *gitlab.Response, error)
}

// ensure mockUsersService implements the gitlab.UsersServiceInterface
var _ gitlab.UsersServiceInterface = &mockUsersService{}

func (m *mockUsersService) CurrentUser(options ...gitlab.RequestOptionFunc) (*gitlab.User, *gitlab.Response, error) {
	return m.currentUserFunc(options...)
}

func (m *mockUsersService) ActivateUser(user int, options ...gitlab.RequestOptionFunc) error {
	return nil
}

func (m *mockUsersService) ApproveUser(user int, options ...gitlab.RequestOptionFunc) error {
	return nil
}

func (m *mockUsersService) BlockUser(user int, options ...gitlab.RequestOptionFunc) error {
	return nil
}

func (m *mockUsersService) CreateUser(opt *gitlab.CreateUserOptions, options ...gitlab.RequestOptionFunc) (*gitlab.User, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) DeleteUser(user int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	return nil, nil
}

func (m *mockUsersService) GetUser(user int, opt gitlab.GetUsersOptions, options ...gitlab.RequestOptionFunc) (*gitlab.User, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) ListUsers(opt *gitlab.ListUsersOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.User, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) UnblockUser(user int, options ...gitlab.RequestOptionFunc) error {
	return nil
}

func (m *mockUsersService) UpdateUser(user int, opt *gitlab.ModifyUserOptions, options ...gitlab.RequestOptionFunc) (*gitlab.User, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) AddEmail(opt *gitlab.AddEmailOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Email, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) AddEmailForUser(user int, opt *gitlab.AddEmailOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Email, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) AddGPGKey(opt *gitlab.AddGPGKeyOptions, options ...gitlab.RequestOptionFunc) (*gitlab.GPGKey, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) AddGPGKeyForUser(user int, opt *gitlab.AddGPGKeyOptions, options ...gitlab.RequestOptionFunc) (*gitlab.GPGKey, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) AddSSHKey(opt *gitlab.AddSSHKeyOptions, options ...gitlab.RequestOptionFunc) (*gitlab.SSHKey, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) AddSSHKeyForUser(user int, opt *gitlab.AddSSHKeyOptions, options ...gitlab.RequestOptionFunc) (*gitlab.SSHKey, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) CurrentUserStatus(options ...gitlab.RequestOptionFunc) (*gitlab.UserStatus, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) GetUserStatus(user int, options ...gitlab.RequestOptionFunc) (*gitlab.UserStatus, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) SetUserStatus(opt *gitlab.UserStatusOptions, options ...gitlab.RequestOptionFunc) (*gitlab.UserStatus, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) GetUserAssociationsCount(user int, options ...gitlab.RequestOptionFunc) (*gitlab.UserAssociationsCount, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) ListSSHKeys(opt *gitlab.ListSSHKeysOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.SSHKey, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) ListSSHKeysForUser(uid interface{}, opt *gitlab.ListSSHKeysForUserOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.SSHKey, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) GetSSHKey(key int, options ...gitlab.RequestOptionFunc) (*gitlab.SSHKey, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) DeleteSSHKey(key int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	return nil, nil
}

func (m *mockUsersService) ListGPGKeys(options ...gitlab.RequestOptionFunc) ([]*gitlab.GPGKey, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) GetGPGKey(key int, options ...gitlab.RequestOptionFunc) (*gitlab.GPGKey, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) DeleteGPGKey(key int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	return nil, nil
}

func (m *mockUsersService) ListEmails(options ...gitlab.RequestOptionFunc) ([]*gitlab.Email, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) ListEmailsForUser(user int, opt *gitlab.ListEmailsForUserOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Email, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) GetEmail(email int, options ...gitlab.RequestOptionFunc) (*gitlab.Email, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) DeleteEmail(email int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	return nil, nil
}

func (m *mockUsersService) DeleteEmailForUser(user, email int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	return nil, nil
}

func (m *mockUsersService) DeactivateUser(user int, options ...gitlab.RequestOptionFunc) error {
	return nil
}

func (m *mockUsersService) RejectUser(user int, options ...gitlab.RequestOptionFunc) error {
	return nil
}

func (m *mockUsersService) GetAllImpersonationTokens(user int, opt *gitlab.GetAllImpersonationTokensOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.ImpersonationToken, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) GetImpersonationToken(user, token int, options ...gitlab.RequestOptionFunc) (*gitlab.ImpersonationToken, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) CreateImpersonationToken(user int, opt *gitlab.CreateImpersonationTokenOptions, options ...gitlab.RequestOptionFunc) (*gitlab.ImpersonationToken, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) RevokeImpersonationToken(user, token int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	return nil, nil
}

func (m *mockUsersService) CreatePersonalAccessToken(user int, opt *gitlab.CreatePersonalAccessTokenOptions, options ...gitlab.RequestOptionFunc) (*gitlab.PersonalAccessToken, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) CreatePersonalAccessTokenForCurrentUser(opt *gitlab.CreatePersonalAccessTokenForCurrentUserOptions, options ...gitlab.RequestOptionFunc) (*gitlab.PersonalAccessToken, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) GetUserActivities(opt *gitlab.GetUserActivitiesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.UserActivity, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) GetUserMemberships(user int, opt *gitlab.GetUserMembershipOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.UserMembership, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) DisableTwoFactor(user int, options ...gitlab.RequestOptionFunc) error {
	return nil
}

func (m *mockUsersService) CreateUserRunner(opts *gitlab.CreateUserRunnerOptions, options ...gitlab.RequestOptionFunc) (*gitlab.UserRunner, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) CreateServiceAccountUser(opts *gitlab.CreateServiceAccountUserOptions, options ...gitlab.RequestOptionFunc) (*gitlab.User, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) ListServiceAccounts(opt *gitlab.ListServiceAccountsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.ServiceAccount, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) UploadAvatar(avatar io.Reader, filename string, options ...gitlab.RequestOptionFunc) (*gitlab.User, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) DeleteUserIdentity(user int, provider string, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	return nil, nil
}

func (m *mockUsersService) ListUserContributionEvents(uid interface{}, opt *gitlab.ListContributionEventsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.ContributionEvent, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) BanUser(user int, options ...gitlab.RequestOptionFunc) error {
	return nil
}

func (m *mockUsersService) DeleteGPGKeyForUser(user, key int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	return nil, nil
}

func (m *mockUsersService) DeleteSSHKeyForUser(user, key int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	return nil, nil
}

func (m *mockUsersService) GetGPGKeyForUser(user, key int, options ...gitlab.RequestOptionFunc) (*gitlab.GPGKey, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) GetSSHKeyForUser(user, key int, options ...gitlab.RequestOptionFunc) (*gitlab.SSHKey, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) ListGPGKeysForUser(user int, options ...gitlab.RequestOptionFunc) ([]*gitlab.GPGKey, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) ModifyUser(user int, opt *gitlab.ModifyUserOptions, options ...gitlab.RequestOptionFunc) (*gitlab.User, *gitlab.Response, error) {
	return nil, nil, nil
}

func (m *mockUsersService) UnbanUser(user int, options ...gitlab.RequestOptionFunc) error {
	return nil
}
