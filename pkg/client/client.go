package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	debug      bool
}

type Option func(*Client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}

func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL: "https://api.funstat.info",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

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

func strPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (c *Client) GetGroup(ctx context.Context, id int64) ([]byte, error) {
	path := fmt.Sprintf("/api/v1/groups/%d", id)
	return c.doRequest(ctx, http.MethodGet, path, nil, nil)
}

func (c *Client) GetUserReputation(ctx context.Context, userID int64) ([]byte, error) {
	query := url.Values{}
	query.Set("id", strconv.FormatInt(userID, 10))
	return c.doRequest(ctx, http.MethodGet, "/api/v1/users/reputation", query, nil)
}

func (c *Client) ResolveUsernames(ctx context.Context, usernames []string) (*ResolvedUserArrayAPIAnswer, error) {
	query := url.Values{}
	for _, username := range usernames {
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

type GetUserMessagesOptions struct {
	GroupID      *int64
	TextContains *string
	MediaCode    *int32
	Page         int32
	PageSize     int32
}

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

func (c *Client) GetUserUsernames(ctx context.Context, userID int64) (*UsernameHistoryAPIAnswer, error) {
	path := fmt.Sprintf("/api/v1/users/%d/usernames", userID)

	respBody, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result UsernameHistoryAPIAnswer
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// TextSearchOptions contains options for text search
type TextSearchOptions struct {
	Page     int32
	PageSize int32
}

// TextSearch searches for who and where wrote specified text (COST: 0.1 per request)
func (c *Client) TextSearch(ctx context.Context, text string, opts *TextSearchOptions) (*TextSearchAPIAnswer, error) {
	query := url.Values{}
	query.Set("input", text)

	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(int(opts.Page)))
		}
		if opts.PageSize > 0 {
			query.Set("pageSize", strconv.Itoa(int(opts.PageSize)))
		}
	}

	respBody, err := c.doRequest(ctx, http.MethodGet, "/api/v1/text/search", query, nil)
	if err != nil {
		return nil, err
	}

	var result TextSearchAPIAnswer
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetCommonGroups returns common groups for specified users (COST: 0.5 per request)
// All requested users must be members of the returned groups
func (c *Client) GetCommonGroups(ctx context.Context, userIDs []int64) (*CommonGroupsAPIAnswer, error) {
	query := url.Values{}
	for _, id := range userIDs {
		query.Add("id", strconv.FormatInt(id, 10))
	}

	respBody, err := c.doRequest(ctx, http.MethodGet, "/api/v1/groups/common_groups", query, nil)
	if err != nil {
		return nil, err
	}

	var result CommonGroupsAPIAnswer
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetCommonGroupsStat returns users who have common groups with specified user (COST: 5)
func (c *Client) GetCommonGroupsStat(ctx context.Context, userID int64) (*CommonGroupsStatAPIAnswer, error) {
	path := fmt.Sprintf("/api/v1/users/%d/common_groups_stat", userID)

	respBody, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result CommonGroupsStatAPIAnswer
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetUsernameUsage searches username usage the same way as the bot (no cost specified)
// Returns: 1=actual users, 2=past usage by users, 3=group/channel actual usage, 4=mentions in group/channel descriptions
func (c *Client) GetUsernameUsage(ctx context.Context, username string) (*UsernameUsageAPIAnswer, error) {
	query := url.Values{}
	username = strings.TrimPrefix(username, "@")
	query.Set("username", username)

	respBody, err := c.doRequest(ctx, http.MethodGet, "/api/v1/users/username_usage", query, nil)
	if err != nil {
		return nil, err
	}

	var result UsernameUsageAPIAnswer
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
