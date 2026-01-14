package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultBaseURL    = "https://api.anthropic.com"
	DefaultAPIVersion = "2023-06-01"
)

// Client is the Anthropic Admin API client
type Client struct {
	BaseURL    string
	AdminKey   string
	APIVersion string
	HTTPClient *http.Client
}

// NewClient creates a new Anthropic Admin API client
func NewClient(adminKey string) *Client {
	return &Client{
		BaseURL:    DefaultBaseURL,
		AdminKey:   adminKey,
		APIVersion: DefaultAPIVersion,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// WithBaseURL sets a custom base URL (useful for testing)
func (c *Client) WithBaseURL(baseURL string) *Client {
	c.BaseURL = baseURL
	return c
}

// APIError represents an error response from the Anthropic API
type APIError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Error   struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e *APIError) String() string {
	if e.Error.Message != "" {
		return fmt.Sprintf("%s: %s", e.Error.Type, e.Error.Message)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// doRequest performs an HTTP request to the Anthropic Admin API
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", c.AdminKey)
	req.Header.Set("anthropic-version", c.APIVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
		}
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, apiErr.String())
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// ListResponse is a generic paginated list response
type ListResponse[T any] struct {
	Data    []T     `json:"data"`
	HasMore bool    `json:"has_more"`
	FirstID *string `json:"first_id,omitempty"`
	LastID  *string `json:"last_id,omitempty"`
}

// ============================================================================
// Workspace Operations
// ============================================================================

// Workspace represents an Anthropic workspace
type Workspace struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
	ArchivedAt  string `json:"archived_at,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// CreateWorkspaceRequest represents the request to create a workspace
type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}

// UpdateWorkspaceRequest represents the request to update a workspace
type UpdateWorkspaceRequest struct {
	Name string `json:"name"`
}

// ListWorkspaces retrieves all workspaces
func (c *Client) ListWorkspaces(ctx context.Context, limit int, beforeID, afterID string) (*ListResponse[Workspace], error) {
	path := "/v1/organizations/workspaces"
	params := []string{}
	if limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", limit))
	}
	if beforeID != "" {
		params = append(params, fmt.Sprintf("before_id=%s", beforeID))
	}
	if afterID != "" {
		params = append(params, fmt.Sprintf("after_id=%s", afterID))
	}
	if len(params) > 0 {
		path += "?"
		for i, p := range params {
			if i > 0 {
				path += "&"
			}
			path += p
		}
	}

	var result ListResponse[Workspace]
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	return &result, err
}

// GetWorkspace retrieves a workspace by ID
func (c *Client) GetWorkspace(ctx context.Context, workspaceID string) (*Workspace, error) {
	var workspace Workspace
	err := c.doRequest(ctx, http.MethodGet, "/v1/organizations/workspaces/"+workspaceID, nil, &workspace)
	return &workspace, err
}

// CreateWorkspace creates a new workspace
func (c *Client) CreateWorkspace(ctx context.Context, req *CreateWorkspaceRequest) (*Workspace, error) {
	var workspace Workspace
	err := c.doRequest(ctx, http.MethodPost, "/v1/organizations/workspaces", req, &workspace)
	return &workspace, err
}

// UpdateWorkspace updates an existing workspace
func (c *Client) UpdateWorkspace(ctx context.Context, workspaceID string, req *UpdateWorkspaceRequest) (*Workspace, error) {
	var workspace Workspace
	err := c.doRequest(ctx, http.MethodPost, "/v1/organizations/workspaces/"+workspaceID, req, &workspace)
	return &workspace, err
}

// ArchiveWorkspace archives a workspace
func (c *Client) ArchiveWorkspace(ctx context.Context, workspaceID string) (*Workspace, error) {
	var workspace Workspace
	err := c.doRequest(ctx, http.MethodPost, "/v1/organizations/workspaces/"+workspaceID+"/archive", nil, &workspace)
	return &workspace, err
}

// ============================================================================
// API Key Operations
// ============================================================================

// APIKey represents an Anthropic API key
type APIKey struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Hint        string `json:"hint,omitempty"` // Last 4 characters
	CreatedAt   string `json:"created_at"`
	CreatedBy   *Actor `json:"created_by,omitempty"`
	Status      string `json:"status"` // active, inactive, archived
	WorkspaceID string `json:"workspace_id,omitempty"`
	// Only returned on creation
	Key string `json:"key,omitempty"`
}

// Actor represents who performed an action
type Actor struct {
	ID   string `json:"id"`
	Type string `json:"type"` // user, api_key, system
}

// CreateAPIKeyRequest represents the request to create an API key
type CreateAPIKeyRequest struct {
	Name        string `json:"name"`
	WorkspaceID string `json:"workspace_id,omitempty"`
}

// UpdateAPIKeyRequest represents the request to update an API key
type UpdateAPIKeyRequest struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"` // active, inactive
}

// ListAPIKeys retrieves all API keys
func (c *Client) ListAPIKeys(ctx context.Context, limit int, beforeID, afterID, status, workspaceID string) (*ListResponse[APIKey], error) {
	path := "/v1/organizations/api_keys"
	params := []string{}
	if limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", limit))
	}
	if beforeID != "" {
		params = append(params, fmt.Sprintf("before_id=%s", beforeID))
	}
	if afterID != "" {
		params = append(params, fmt.Sprintf("after_id=%s", afterID))
	}
	if status != "" {
		params = append(params, fmt.Sprintf("status=%s", status))
	}
	if workspaceID != "" {
		params = append(params, fmt.Sprintf("workspace_id=%s", workspaceID))
	}
	if len(params) > 0 {
		path += "?"
		for i, p := range params {
			if i > 0 {
				path += "&"
			}
			path += p
		}
	}

	var result ListResponse[APIKey]
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	return &result, err
}

// GetAPIKey retrieves an API key by ID
func (c *Client) GetAPIKey(ctx context.Context, apiKeyID string) (*APIKey, error) {
	var apiKey APIKey
	err := c.doRequest(ctx, http.MethodGet, "/v1/organizations/api_keys/"+apiKeyID, nil, &apiKey)
	return &apiKey, err
}

// CreateAPIKey creates a new API key
func (c *Client) CreateAPIKey(ctx context.Context, req *CreateAPIKeyRequest) (*APIKey, error) {
	var apiKey APIKey
	err := c.doRequest(ctx, http.MethodPost, "/v1/organizations/api_keys", req, &apiKey)
	return &apiKey, err
}

// UpdateAPIKey updates an existing API key
func (c *Client) UpdateAPIKey(ctx context.Context, apiKeyID string, req *UpdateAPIKeyRequest) (*APIKey, error) {
	var apiKey APIKey
	err := c.doRequest(ctx, http.MethodPost, "/v1/organizations/api_keys/"+apiKeyID, req, &apiKey)
	return &apiKey, err
}

// DeleteAPIKey deletes an API key (archives it)
func (c *Client) DeleteAPIKey(ctx context.Context, apiKeyID string) error {
	// Archive the key by setting status to archived
	_, err := c.UpdateAPIKey(ctx, apiKeyID, &UpdateAPIKeyRequest{Status: "archived"})
	return err
}

// ============================================================================
// Workspace Member Operations
// ============================================================================

// WorkspaceMember represents a user's membership in a workspace
type WorkspaceMember struct {
	UserID        string `json:"user_id"`
	WorkspaceID   string `json:"workspace_id"`
	WorkspaceRole string `json:"workspace_role"` // workspace_user, workspace_admin, workspace_developer
	Type          string `json:"type"`
}

// AddWorkspaceMemberRequest represents the request to add a member to a workspace
type AddWorkspaceMemberRequest struct {
	UserID        string `json:"user_id"`
	WorkspaceRole string `json:"workspace_role"`
}

// UpdateWorkspaceMemberRequest represents the request to update a workspace member
type UpdateWorkspaceMemberRequest struct {
	WorkspaceRole string `json:"workspace_role"`
}

// ListWorkspaceMembers retrieves all members of a workspace
func (c *Client) ListWorkspaceMembers(ctx context.Context, workspaceID string, limit int, beforeID, afterID string) (*ListResponse[WorkspaceMember], error) {
	path := fmt.Sprintf("/v1/organizations/workspaces/%s/members", workspaceID)
	params := []string{}
	if limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", limit))
	}
	if beforeID != "" {
		params = append(params, fmt.Sprintf("before_id=%s", beforeID))
	}
	if afterID != "" {
		params = append(params, fmt.Sprintf("after_id=%s", afterID))
	}
	if len(params) > 0 {
		path += "?"
		for i, p := range params {
			if i > 0 {
				path += "&"
			}
			path += p
		}
	}

	var result ListResponse[WorkspaceMember]
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	return &result, err
}

// GetWorkspaceMember retrieves a workspace member
func (c *Client) GetWorkspaceMember(ctx context.Context, workspaceID, userID string) (*WorkspaceMember, error) {
	var member WorkspaceMember
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/organizations/workspaces/%s/members/%s", workspaceID, userID), nil, &member)
	return &member, err
}

// AddWorkspaceMember adds a user to a workspace
func (c *Client) AddWorkspaceMember(ctx context.Context, workspaceID string, req *AddWorkspaceMemberRequest) (*WorkspaceMember, error) {
	var member WorkspaceMember
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/v1/organizations/workspaces/%s/members", workspaceID), req, &member)
	return &member, err
}

// UpdateWorkspaceMember updates a workspace member's role
func (c *Client) UpdateWorkspaceMember(ctx context.Context, workspaceID, userID string, req *UpdateWorkspaceMemberRequest) (*WorkspaceMember, error) {
	var member WorkspaceMember
	err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/v1/organizations/workspaces/%s/members/%s", workspaceID, userID), req, &member)
	return &member, err
}

// RemoveWorkspaceMember removes a user from a workspace
func (c *Client) RemoveWorkspaceMember(ctx context.Context, workspaceID, userID string) error {
	return c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/organizations/workspaces/%s/members/%s", workspaceID, userID), nil, nil)
}

// ============================================================================
// Organization Member Operations
// ============================================================================

// OrganizationMember represents a user's membership in the organization
type OrganizationMember struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"` // user, admin, developer
}

// UpdateOrganizationMemberRequest represents the request to update an org member
type UpdateOrganizationMemberRequest struct {
	Role string `json:"role"`
}

// ListOrganizationMembers retrieves all organization members
func (c *Client) ListOrganizationMembers(ctx context.Context, limit int, beforeID, afterID string) (*ListResponse[OrganizationMember], error) {
	path := "/v1/organizations/users"
	params := []string{}
	if limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", limit))
	}
	if beforeID != "" {
		params = append(params, fmt.Sprintf("before_id=%s", beforeID))
	}
	if afterID != "" {
		params = append(params, fmt.Sprintf("after_id=%s", afterID))
	}
	if len(params) > 0 {
		path += "?"
		for i, p := range params {
			if i > 0 {
				path += "&"
			}
			path += p
		}
	}

	var result ListResponse[OrganizationMember]
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	return &result, err
}

// GetOrganizationMember retrieves an organization member by ID
func (c *Client) GetOrganizationMember(ctx context.Context, userID string) (*OrganizationMember, error) {
	var member OrganizationMember
	err := c.doRequest(ctx, http.MethodGet, "/v1/organizations/users/"+userID, nil, &member)
	return &member, err
}

// UpdateOrganizationMember updates an organization member's role
func (c *Client) UpdateOrganizationMember(ctx context.Context, userID string, req *UpdateOrganizationMemberRequest) (*OrganizationMember, error) {
	var member OrganizationMember
	err := c.doRequest(ctx, http.MethodPost, "/v1/organizations/users/"+userID, req, &member)
	return &member, err
}

// RemoveOrganizationMember removes a user from the organization
func (c *Client) RemoveOrganizationMember(ctx context.Context, userID string) error {
	return c.doRequest(ctx, http.MethodDelete, "/v1/organizations/users/"+userID, nil, nil)
}

// ============================================================================
// Invite Operations
// ============================================================================

// Invite represents an invitation to join the organization
type Invite struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	Email          string `json:"email"`
	Role           string `json:"role"` // user, admin, developer
	Status         string `json:"status"` // pending, accepted, expired, deleted
	CreatedAt      string `json:"created_at"`
	ExpiresAt      string `json:"expires_at"`
	InviterID      string `json:"inviter_id,omitempty"`
	WorkspaceIDs   []string `json:"workspace_ids,omitempty"`
}

// CreateInviteRequest represents the request to create an invite
type CreateInviteRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// ListInvites retrieves all invites
func (c *Client) ListInvites(ctx context.Context, limit int, beforeID, afterID string) (*ListResponse[Invite], error) {
	path := "/v1/organizations/invites"
	params := []string{}
	if limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", limit))
	}
	if beforeID != "" {
		params = append(params, fmt.Sprintf("before_id=%s", beforeID))
	}
	if afterID != "" {
		params = append(params, fmt.Sprintf("after_id=%s", afterID))
	}
	if len(params) > 0 {
		path += "?"
		for i, p := range params {
			if i > 0 {
				path += "&"
			}
			path += p
		}
	}

	var result ListResponse[Invite]
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	return &result, err
}

// GetInvite retrieves an invite by ID
func (c *Client) GetInvite(ctx context.Context, inviteID string) (*Invite, error) {
	var invite Invite
	err := c.doRequest(ctx, http.MethodGet, "/v1/organizations/invites/"+inviteID, nil, &invite)
	return &invite, err
}

// CreateInvite creates a new invite
func (c *Client) CreateInvite(ctx context.Context, req *CreateInviteRequest) (*Invite, error) {
	var invite Invite
	err := c.doRequest(ctx, http.MethodPost, "/v1/organizations/invites", req, &invite)
	return &invite, err
}

// DeleteInvite deletes/cancels an invite
func (c *Client) DeleteInvite(ctx context.Context, inviteID string) error {
	return c.doRequest(ctx, http.MethodDelete, "/v1/organizations/invites/"+inviteID, nil, nil)
}
