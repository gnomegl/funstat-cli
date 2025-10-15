package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client represents the Funstat API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	debug      bool
}

// Option is a functional option for configuring the client
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithDebug enables debug mode
func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}

// New creates a new Funstat API client
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL: "http://api.funstat.info",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(ctx context.Context, method, path string, query url.Values, body interface{}) ([]byte, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	u.Path = path
	if query != nil {
		u.RawQuery = query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.debug {
		fmt.Printf("Request: %s %s\n", method, u.String())
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if c.debug {
		fmt.Printf("Response Status: %d\n", resp.StatusCode)
		fmt.Printf("Response Body: %s\n", string(respBody))
	}

	if resp.StatusCode >= 400 {
		var problem AppProblem
		if err := json.Unmarshal(respBody, &problem); err == nil {
			return nil, fmt.Errorf("API error: %s - %s",
				strPtr(problem.Title), strPtr(problem.Detail))
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	return respBody, nil
}

// strPtr safely dereferences a string pointer
func strPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Groups API

// GetGroup gets basic info, links and today stats for a group
func (c *Client) GetGroup(ctx context.Context, id int64) ([]byte, error) {
	path := fmt.Sprintf("/api/v1/groups/%d", id)
	return c.doRequest(ctx, http.MethodGet, path, nil, nil)
}

// Users API

// GetUserReputation returns user reputation information (FREE)
func (c *Client) GetUserReputation(ctx context.Context, userID int64) ([]byte, error) {
	query := url.Values{}
	query.Set("id", strconv.FormatInt(userID, 10))
	return c.doRequest(ctx, http.MethodGet, "/api/v1/users/reputation", query, nil)
}

// ResolveUsernames resolves telegram usernames to user info (Cost: 0.10 per success)
func (c *Client) ResolveUsernames(ctx context.Context, usernames []string) (*ResolvedUserArrayAPIAnswer, error) {
	query := url.Values{}
	for _, username := range usernames {
		// Remove @ if present
		username = strings.TrimPrefix(username, "@")
		query.Add("name", username)
	}

	respBody, err := c.doRequest(ctx, http.MethodGet, "/api/v1/users/resolve_username", query, nil)
	if err != nil {
		return nil, err
	}

	var result ResolvedUserArrayAPIAnswer
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetUserStatsMin returns basic user stats (FREE)
func (c *Client) GetUserStatsMin(ctx context.Context, userID int64) (*UserStatsMinAPIAnswer, error) {
	path := fmt.Sprintf("/api/v1/users/%d/stats_min", userID)
	respBody, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result UserStatsMinAPIAnswer
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetUserStats returns full user stats (Cost: 1)
func (c *Client) GetUserStats(ctx context.Context, userID int64) (*UserStatsAPIAnswer, error) {
	path := fmt.Sprintf("/api/v1/users/%d/stats", userID)
	respBody, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result UserStatsAPIAnswer
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetUsersByID gets user info by telegram ID (Cost: 0.10 per success)
func (c *Client) GetUsersByID(ctx context.Context, userIDs []int64) (*ResolvedUserArrayAPIAnswer, error) {
	query := url.Values{}
	for _, id := range userIDs {
		query.Add("id", strconv.FormatInt(id, 10))
	}

	respBody, err := c.doRequest(ctx, http.MethodGet, "/api/v1/users/basic_info_by_id", query, nil)
	if err != nil {
		return nil, err
	}

	var result ResolvedUserArrayAPIAnswer
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetUserGroupsCount returns total count of user groups (FREE)
func (c *Client) GetUserGroupsCount(ctx context.Context, userID int64, onlyWithMessages bool) (int32, error) {
	path := fmt.Sprintf("/api/v1/users/%d/groups_count", userID)
	query := url.Values{}
	query.Set("onlyMsg", strconv.FormatBool(onlyWithMessages))

	respBody, err := c.doRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return 0, err
	}

	var count int32
	if err := json.Unmarshal(respBody, &count); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return count, nil
}

// GetUserMessagesOptions represents options for GetUserMessages
type GetUserMessagesOptions struct {
	GroupID      *int64
	TextContains *string
	MediaCode    *int32
	Page         int32
	PageSize     int32
}

// GetUserMessages gets user messages (Cost: 10 per user if success)
func (c *Client) GetUserMessages(ctx context.Context, userID int64, opts GetUserMessagesOptions) (*UserMsgArrayAPIAnswerPaged, error) {
	path := fmt.Sprintf("/api/v1/users/%d/messages", userID)
	query := url.Values{}

	if opts.GroupID != nil {
		query.Set("group_id", strconv.FormatInt(*opts.GroupID, 10))
	}
	if opts.TextContains != nil {
		query.Set("text_contains", *opts.TextContains)
	}
	if opts.MediaCode != nil {
		query.Set("media_code", strconv.Itoa(int(*opts.MediaCode)))
	}
	query.Set("page", strconv.Itoa(int(opts.Page)))
	query.Set("pageSize", strconv.Itoa(int(opts.PageSize)))

	respBody, err := c.doRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}

	var result UserMsgArrayAPIAnswerPaged
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetUserMessagesCount returns total count of user messages (FREE)
func (c *Client) GetUserMessagesCount(ctx context.Context, userID int64) (int32, error) {
	path := fmt.Sprintf("/api/v1/users/%d/messages_count", userID)

	respBody, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return 0, err
	}

	var count int32
	if err := json.Unmarshal(respBody, &count); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return count, nil
}

// GetUserGroups returns known user groups (Cost: 5)
func (c *Client) GetUserGroups(ctx context.Context, userID int64) (*UserChatInfoArrayAPIAnswer, error) {
	path := fmt.Sprintf("/api/v1/users/%d/groups", userID)

	respBody, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result UserChatInfoArrayAPIAnswer
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetUserNames returns user (firstname + lastname) history (Cost: 3)
func (c *Client) GetUserNames(ctx context.Context, userID int64) (*UserChatInfoArrayAPIAnswer, error) {
	path := fmt.Sprintf("/api/v1/users/%d/names", userID)

	respBody, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result UserChatInfoArrayAPIAnswer
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetUserUsernames returns @usernames history (Cost: 3)
func (c *Client) GetUserUsernames(ctx context.Context, userID int64) (*UserChatInfoArrayAPIAnswer, error) {
	path := fmt.Sprintf("/api/v1/users/%d/usernames", userID)

	respBody, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result UserChatInfoArrayAPIAnswer
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
