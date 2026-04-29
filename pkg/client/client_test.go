package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		apiKey string
		opts   []Option
		want   func(*Client) bool
	}{
		{
			name:   "basic client",
			apiKey: "test-key",
			want: func(c *Client) bool {
				return c.apiKey == "test-key" &&
					c.baseURL.String() == "http://telelog.org" &&
					c.httpClient.Timeout == 30*time.Second
			},
		},
		{
			name:   "with custom base URL",
			apiKey: "test-key",
			opts:   []Option{WithBaseURL("https://custom.api")},
			want: func(c *Client) bool {
				return c.baseURL.String() == "https://custom.api"
			},
		},
		{
			name:   "with debug enabled",
			apiKey: "test-key",
			opts:   []Option{WithDebug(true)},
			want: func(c *Client) bool {
				return c.debug == true
			},
		},
		{
			name:   "with custom http client",
			apiKey: "test-key",
			opts: []Option{WithHTTPClient(&http.Client{
				Timeout: 10 * time.Second,
			})},
			want: func(c *Client) bool {
				return c.httpClient.Timeout == 10*time.Second
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.apiKey, tt.opts...)
			assert.True(t, tt.want(got))
		})
	}
}

func TestResolveUsernames(t *testing.T) {
	tests := []struct {
		name       string
		usernames  []string
		response   string
		statusCode int
		wantErr    bool
		validate   func(*testing.T, *ResolvedUserArrayAPIAnswer)
	}{
		{
			name:      "single username success",
			usernames: []string{"testuser"},
			response: `{
				"success": true,
				"tech": {"request_cost": 0.1, "current_ballance": 100, "request_duration": "10ms"},
				"data": [{"id": 123456, "username": "testuser", "is_active": true, "is_bot": false}]
			}`,
			statusCode: 200,
			wantErr:    false,
			validate: func(t *testing.T, result *ResolvedUserArrayAPIAnswer) {
				require.NotNil(t, result)
				assert.True(t, result.Success)
				assert.Len(t, result.Data, 1)
				assert.Equal(t, int64(123456), result.Data[0].ID)
				assert.Equal(t, "testuser", *result.Data[0].Username)
			},
		},
		{
			name:      "multiple usernames",
			usernames: []string{"user1", "@user2"},
			response: `{
				"success": true,
				"tech": {"request_cost": 0.2, "current_ballance": 99.8, "request_duration": "20ms"},
				"data": [
					{"id": 111, "username": "user1", "is_active": true, "is_bot": false},
					{"id": 222, "username": "user2", "is_active": true, "is_bot": false}
				]
			}`,
			statusCode: 200,
			wantErr:    false,
			validate: func(t *testing.T, result *ResolvedUserArrayAPIAnswer) {
				require.NotNil(t, result)
				assert.Len(t, result.Data, 2)
			},
		},
		{
			name:       "API error",
			usernames:  []string{"invalid"},
			response:   `{"title": "Not Found", "detail": "User not found"}`,
			statusCode: 404,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/v1/users/resolve_username", r.URL.Path)
				assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := New("test-key", WithBaseURL(server.URL))
			result, err := client.ResolveUsernames(context.Background(), tt.usernames)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestGetUserStats(t *testing.T) {
	tests := []struct {
		name       string
		userID     int64
		response   string
		statusCode int
		wantErr    bool
		validate   func(*testing.T, *UserStatsAPIAnswer)
	}{
		{
			name:   "full stats success",
			userID: 123456,
			response: `{
				"success": true,
				"tech": {"request_cost": 1, "current_ballance": 99, "request_duration": "50ms"},
				"data": {
					"id": 123456,
					"first_name": "John",
					"last_name": "Doe",
					"is_bot": false,
					"is_active": true,
					"total_msg_count": 1000,
					"msg_in_groups_count": 800,
					"adm_in_groups": 5,
					"total_groups": 20,
					"reply_percent": 15.5,
					"media_percent": 30.2,
					"link_percent": 10.1
				}
			}`,
			statusCode: 200,
			wantErr:    false,
			validate: func(t *testing.T, result *UserStatsAPIAnswer) {
				require.NotNil(t, result)
				assert.True(t, result.Success)
				assert.NotNil(t, result.Data)
				assert.Equal(t, int64(123456), result.Data.ID)
				assert.Equal(t, "John", *result.Data.FirstName)
				assert.Equal(t, int64(1000), result.Data.TotalMsgCount)
				assert.Equal(t, float32(1), result.Tech.RequestCost)
			},
		},
		{
			name:       "user not found",
			userID:     999999,
			response:   `{"title": "Not Found", "detail": "User not found"}`,
			statusCode: 404,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, "/api/v1/users/")
				assert.Contains(t, r.URL.Path, "/stats")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := New("test-key", WithBaseURL(server.URL))
			result, err := client.GetUserStats(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestGetUserStatsMin(t *testing.T) {
	tests := []struct {
		name       string
		userID     int64
		response   string
		statusCode int
		wantErr    bool
		validate   func(*testing.T, *UserStatsMinAPIAnswer)
	}{
		{
			name:   "minimal stats success",
			userID: 123456,
			response: `{
				"success": true,
				"tech": {"request_cost": 0, "current_ballance": 100, "request_duration": "10ms"},
				"data": {
					"id": 123456,
					"is_bot": false,
					"is_active": true,
					"total_msg_count": 500,
					"total_groups": 10
				}
			}`,
			statusCode: 200,
			wantErr:    false,
			validate: func(t *testing.T, result *UserStatsMinAPIAnswer) {
				require.NotNil(t, result)
				assert.True(t, result.Success)
				assert.Equal(t, int64(123456), result.Data.ID)
				assert.Equal(t, int64(500), result.Data.TotalMsgCount)
				assert.Equal(t, float32(0), result.Tech.RequestCost)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := New("test-key", WithBaseURL(server.URL))
			result, err := client.GetUserStatsMin(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestGetUsersByID(t *testing.T) {
	tests := []struct {
		name       string
		userIDs    []int64
		response   string
		statusCode int
		wantErr    bool
		validate   func(*testing.T, *ResolvedUserArrayAPIAnswer)
	}{
		{
			name:    "single ID success",
			userIDs: []int64{123456},
			response: `{
				"success": true,
				"tech": {"request_cost": 0.1, "current_ballance": 99.9, "request_duration": "15ms"},
				"data": [{"id": 123456, "first_name": "Alice", "is_active": true, "is_bot": false}]
			}`,
			statusCode: 200,
			wantErr:    false,
			validate: func(t *testing.T, result *ResolvedUserArrayAPIAnswer) {
				require.NotNil(t, result)
				assert.Len(t, result.Data, 1)
				assert.Equal(t, int64(123456), result.Data[0].ID)
			},
		},
		{
			name:    "multiple IDs",
			userIDs: []int64{111, 222, 333},
			response: `{
				"success": true,
				"tech": {"request_cost": 0.3, "current_ballance": 99.7, "request_duration": "25ms"},
				"data": [
					{"id": 111, "is_active": true, "is_bot": false},
					{"id": 222, "is_active": true, "is_bot": false},
					{"id": 333, "is_active": false, "is_bot": true}
				]
			}`,
			statusCode: 200,
			wantErr:    false,
			validate: func(t *testing.T, result *ResolvedUserArrayAPIAnswer) {
				assert.Len(t, result.Data, 3)
				assert.True(t, result.Data[2].IsBot)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/v1/users/basic_info_by_id", r.URL.Path)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := New("test-key", WithBaseURL(server.URL))
			result, err := client.GetUsersByID(context.Background(), tt.userIDs)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestGetUserGroups(t *testing.T) {
	tests := []struct {
		name       string
		userID     int64
		response   string
		statusCode int
		wantErr    bool
		validate   func(*testing.T, *UserChatInfoArrayAPIAnswer)
	}{
		{
			name:   "groups list success",
			userID: 123456,
			response: `{
				"success": true,
				"tech": {"request_cost": 5, "current_ballance": 95, "request_duration": "100ms"},
				"data": [
					{
						"chat": {"id": 111, "title": "Test Group", "isPrivate": false},
						"messagesCount": 50,
						"isAdmin": true,
						"isLeft": false
					},
					{
						"chat": {"id": 222, "title": "Another Group", "isPrivate": true},
						"messagesCount": 30,
						"isAdmin": false,
						"isLeft": false
					}
				]
			}`,
			statusCode: 200,
			wantErr:    false,
			validate: func(t *testing.T, result *UserChatInfoArrayAPIAnswer) {
				require.NotNil(t, result)
				assert.Len(t, result.Data, 2)
				assert.Equal(t, "Test Group", result.Data[0].Chat.Title)
				assert.True(t, result.Data[0].IsAdmin)
				assert.Equal(t, float32(5), result.Tech.RequestCost)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := New("test-key", WithBaseURL(server.URL))
			result, err := client.GetUserGroups(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestGetUserGroupsCount(t *testing.T) {
	tests := []struct {
		name             string
		userID           int64
		onlyWithMessages bool
		response         string
		statusCode       int
		wantErr          bool
		expectedCount    int32
	}{
		{
			name:             "groups count with messages",
			userID:           123456,
			onlyWithMessages: true,
			response:         `15`,
			statusCode:       200,
			wantErr:          false,
			expectedCount:    15,
		},
		{
			name:             "all groups count",
			userID:           123456,
			onlyWithMessages: false,
			response:         `25`,
			statusCode:       200,
			wantErr:          false,
			expectedCount:    25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, r.URL.Query().Get("onlyMsg"), strconv.FormatBool(tt.onlyWithMessages))
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := New("test-key", WithBaseURL(server.URL))
			result, err := client.GetUserGroupsCount(context.Background(), tt.userID, tt.onlyWithMessages)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedCount, result)
		})
	}
}

func TestGetUserMessages(t *testing.T) {
	tests := []struct {
		name       string
		userID     int64
		opts       GetUserMessagesOptions
		response   string
		statusCode int
		wantErr    bool
		validate   func(*testing.T, *UserMsgArrayAPIAnswerPaged)
	}{
		{
			name:   "messages with pagination",
			userID: 123456,
			opts: GetUserMessagesOptions{
				Page:     1,
				PageSize: 10,
			},
			response: `{
				"success": true,
				"tech": {"request_cost": 10, "current_ballance": 90, "request_duration": "200ms"},
				"paging": {"total": 100, "currentPage": 1, "pageSize": 10, "totalPages": 10},
				"data": [
					{
						"date": "2024-01-01T12:00:00Z",
						"messageId": 1,
						"text": "Hello",
						"group": {"id": 111, "title": "Test", "isPrivate": false}
					}
				]
			}`,
			statusCode: 200,
			wantErr:    false,
			validate: func(t *testing.T, result *UserMsgArrayAPIAnswerPaged) {
				require.NotNil(t, result)
				assert.True(t, result.Success)
				assert.Equal(t, int32(100), result.Paging.Total)
				assert.Len(t, result.Data, 1)
			},
		},
		{
			name:   "messages with filters",
			userID: 123456,
			opts: GetUserMessagesOptions{
				GroupID:      ptr(int64(111)),
				TextContains: ptr("test"),
				MediaCode:    ptr(int32(1)),
				Page:         1,
				PageSize:     5,
			},
			response: `{
				"success": true,
				"tech": {"request_cost": 10, "current_ballance": 90, "request_duration": "150ms"},
				"paging": {"total": 5, "currentPage": 1, "pageSize": 5, "totalPages": 1},
				"data": []
			}`,
			statusCode: 200,
			wantErr:    false,
			validate: func(t *testing.T, result *UserMsgArrayAPIAnswerPaged) {
				assert.Equal(t, int32(5), result.Paging.Total)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				query := r.URL.Query()
				assert.Equal(t, strconv.Itoa(int(tt.opts.Page)), query.Get("page"))
				assert.Equal(t, strconv.Itoa(int(tt.opts.PageSize)), query.Get("pageSize"))
				if tt.opts.GroupID != nil {
					assert.Equal(t, strconv.FormatInt(*tt.opts.GroupID, 10), query.Get("group_id"))
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := New("test-key", WithBaseURL(server.URL))
			result, err := client.GetUserMessages(context.Background(), tt.userID, tt.opts)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestGetUserMessagesCount(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		response      string
		statusCode    int
		wantErr       bool
		expectedCount int32
	}{
		{
			name:          "messages count success",
			userID:        123456,
			response:      `500`,
			statusCode:    200,
			wantErr:       false,
			expectedCount: 500,
		},
		{
			name:          "zero messages",
			userID:        999999,
			response:      `0`,
			statusCode:    200,
			wantErr:       false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := New("test-key", WithBaseURL(server.URL))
			result, err := client.GetUserMessagesCount(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedCount, result)
		})
	}
}

func TestGetUserNames(t *testing.T) {
	tests := []struct {
		name       string
		userID     int64
		response   string
		statusCode int
		wantErr    bool
		validate   func(*testing.T, *UsernameHistoryAPIAnswer)
	}{
		{
			name:   "names history success",
			userID: 123456,
			response: `{
				"success": true,
				"tech": {"request_cost": 3, "current_ballance": 97, "request_duration": "50ms"},
				"data": [
					{"name": "John Doe", "date_time": "2023-01-01T00:00:00Z"},
					{"name": "John Smith", "date_time": "2023-06-02T00:00:00Z"}
				]
			}`,
			statusCode: 200,
			wantErr:    false,
			validate: func(t *testing.T, result *UsernameHistoryAPIAnswer) {
				require.NotNil(t, result)
				assert.True(t, result.Success)
				assert.Len(t, result.Data, 2)
				assert.Equal(t, "John Doe", result.Data[0].Name)
				assert.Equal(t, "John Smith", result.Data[1].Name)
				assert.Equal(t, float32(3), result.Tech.RequestCost)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, "/names")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := New("test-key", WithBaseURL(server.URL))
			result, err := client.GetUserNames(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestGetUserUsernames(t *testing.T) {
	tests := []struct {
		name       string
		userID     int64
		response   string
		statusCode int
		wantErr    bool
		validate   func(*testing.T, *UsernameHistoryAPIAnswer)
	}{
		{
			name:   "usernames history success",
			userID: 123456,
			response: `{
				"success": true,
				"tech": {"request_cost": 3, "current_ballance": 97, "request_duration": "50ms"},
				"data": [
					{"name": "oldusername1", "date_time": "2023-01-01T00:00:00Z"},
					{"name": "oldusername2", "date_time": "2023-06-02T00:00:00Z"}
				]
			}`,
			statusCode: 200,
			wantErr:    false,
			validate: func(t *testing.T, result *UsernameHistoryAPIAnswer) {
				require.NotNil(t, result)
				assert.True(t, result.Success)
				assert.Len(t, result.Data, 2)
				assert.Equal(t, "oldusername1", result.Data[0].Name)
				assert.Equal(t, "oldusername2", result.Data[1].Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, "/usernames")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := New("test-key", WithBaseURL(server.URL))
			result, err := client.GetUserUsernames(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestGetGroup(t *testing.T) {
	tests := []struct {
		name       string
		groupID    int64
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:    "group info success",
			groupID: 111222333,
			response: `{
				"id": 111222333,
				"title": "Test Group",
				"username": "testgroup",
				"members_count": 1000,
				"is_private": false
			}`,
			statusCode: 200,
			wantErr:    false,
		},
		{
			name:       "group not found",
			groupID:    999999,
			response:   `{"title": "Not Found", "detail": "Group not found"}`,
			statusCode: 404,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Contains(t, r.URL.Path, "/api/v1/groups/")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := New("test-key", WithBaseURL(server.URL))
			result, err := client.GetGroup(context.Background(), tt.groupID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)

			var data map[string]interface{}
			err = json.Unmarshal(result, &data)
			require.NoError(t, err)
		})
	}
}

func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(200)
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.GetUserStatsMin(ctx, 123456)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestErrorWithNonAppProblemJSON(t *testing.T) {
	// When the API returns a 400+ status with valid JSON that isn't an AppProblem,
	// we should fall back to the generic status code error instead of "API error:  - "
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(502)
		w.Write([]byte(`{"error": "bad gateway"}`))
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	_, err := client.GetUserStatsMin(context.Background(), 123456)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status 502")
	assert.NotContains(t, err.Error(), "API error:  - ")
}

func TestAuthenticationHeader(t *testing.T) {
	apiKey := "test-secret-key"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer "+apiKey, authHeader)
		w.WriteHeader(200)
		w.Write([]byte(`{"success": true, "tech": {"request_cost": 0}, "data": {"id": 123, "is_active": true, "is_bot": false}}`))
	}))
	defer server.Close()

	client := New(apiKey, WithBaseURL(server.URL))
	_, err := client.GetUserStatsMin(context.Background(), 123456)
	require.NoError(t, err)
}

func TestTextSearch(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		opts       TextSearchOptions
		response   string
		statusCode int
		wantErr    bool
		validate   func(*testing.T, *TextSearchAPIAnswer)
	}{
		{
			name: "text search success",
			text: "hello world",
			opts: TextSearchOptions{Page: 1, PageSize: 10},
			response: `{
				"success": true,
				"tech": {"request_cost": 0.1, "current_ballance": 99.9, "request_duration": "30ms"},
				"data": {
					"total": 5,
					"data": [
						{"message_id": 100, "user_id": 111, "date": "2024-01-01T12:00:00Z", "text": "hello world!", "is_active": true, "group": {"id": 222, "title": "Test Group", "isPrivate": false, "isChannel": false}}
					],
					"isLastPage": false,
					"pageSize": 10,
					"currentPage": 1,
					"totalPages": 1,
					"isSliding": false
				}
			}`,
			statusCode: 200,
			wantErr:    false,
			validate: func(t *testing.T, result *TextSearchAPIAnswer) {
				require.NotNil(t, result)
				assert.True(t, result.Success)
				assert.Equal(t, int32(5), result.Data.Total)
				assert.Len(t, result.Data.Data, 1)
				assert.Equal(t, "hello world!", result.Data.Data[0].Text)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/v1/text/search", r.URL.Path)
				assert.Equal(t, tt.text, r.URL.Query().Get("text"))
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := New("test-key", WithBaseURL(server.URL))
			result, err := client.TextSearch(context.Background(), tt.text, tt.opts)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestGetCommonGroups(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 0.5, "current_ballance": 99.5, "request_duration": "40ms"},
		"data": [
			{"id": 111, "title": "Shared Group", "isPrivate": false, "isChannel": false}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/groups/common_groups", r.URL.Path)
		ids := r.URL.Query()["id"]
		assert.Len(t, ids, 2)
		w.WriteHeader(200)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	result, err := client.GetCommonGroups(context.Background(), []int64{111, 222})
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "Shared Group", result.Data[0].Title)
}

func TestGetCommonGroupsStat(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 5, "current_ballance": 95, "request_duration": "80ms"},
		"data": [
			{"user_id": 222, "common_groups": 3, "first_name": "Alice", "is_user_active": true},
			{"user_id": 333, "common_groups": 1, "first_name": "Bob", "is_user_active": true}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/common_groups_stat")
		w.WriteHeader(200)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	result, err := client.GetCommonGroupsStat(context.Background(), 111)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, int64(222), result.Data[0].UserID)
	assert.Equal(t, int32(3), result.Data[0].CommonGroups)
}

func TestGetUsernameUsage(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 0.1, "current_ballance": 99.9, "request_duration": "20ms"},
		"data": {
			"actualUsers": [{"id": 111, "username": "testuser", "is_active": true, "is_bot": false}],
			"usageByUsersInThePast": [{"id": 222, "username": "testuser", "is_active": false, "is_bot": false}],
			"actualGroupsOrChannels": [],
			"mentionByChannelOrGroupDesc": []
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/users/username_usage", r.URL.Path)
		assert.Equal(t, "testuser", r.URL.Query().Get("username"))
		w.WriteHeader(200)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	result, err := client.GetUsernameUsage(context.Background(), "@testuser")
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data.ActualUsers, 1)
	assert.Len(t, result.Data.UsageByUsersInThePast, 1)
}

func TestGetGroupMembers(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 15, "current_ballance": 85, "request_duration": "200ms"},
		"data": [
			{"id": 111, "first_name": "Alice", "is_active": true, "today_msg": 5, "has_photo": true},
			{"id": 222, "first_name": "Bob", "username": "bob", "is_active": true, "today_msg": 0, "has_photo": false}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/members")
		w.WriteHeader(200)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	result, err := client.GetGroupMembers(context.Background(), 999)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, int64(111), result.Data[0].ID)
	assert.Equal(t, "bob", *result.Data[1].Username)
}

func TestGetGiftsRelation(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 5, "current_ballance": 95, "request_duration": "100ms"},
		"data": [
			{
				"from_user_id": 111,
				"from_first_name": "Alice",
				"from_is_active": true,
				"to_user_id": 222,
				"to_first_name": "Bob",
				"to_is_active": true,
				"last_gift_date": "2024-06-15T10:00:00Z"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/gifts_relation")
		assert.Equal(t, "1", r.URL.Query().Get("page"))
		assert.Equal(t, "20", r.URL.Query().Get("pageSize"))
		w.WriteHeader(200)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	opts := GiftsRelationOptions{Page: 1, PageSize: 20}
	result, err := client.GetGiftsRelation(context.Background(), 111, opts)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, int64(111), result.Data[0].FromUserID)
	assert.Equal(t, int64(222), result.Data[0].ToUserID)
}

func TestGetUserStickers(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 1, "current_ballance": 99, "request_duration": "50ms"},
		"data": [
			{"sticker_set_id": 12345, "last_seen": "2024-01-01", "min_seen": "2023-01-01", "title": "My Pack", "short_name": "mypack", "stickers_count": 30}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/stickers")
		w.WriteHeader(200)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	result, err := client.GetUserStickers(context.Background(), 111)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "My Pack", *result.Data[0].Title)
	assert.Equal(t, int64(12345), result.Data[0].StickerSetID)
}

func TestGetUserReputation(t *testing.T) {
	response := `{"spam_score": 0, "is_scammer": false}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/users/reputation", r.URL.Path)
		assert.Equal(t, "111", r.URL.Query().Get("id"))
		w.WriteHeader(200)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	result, err := client.GetUserReputation(context.Background(), 111)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetBotRandom(t *testing.T) {
	response := `{"user_id": 12345, "username": "randombot"}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/bot/random", r.URL.Path)
		w.WriteHeader(200)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := New("test-key", WithBaseURL(server.URL))
	result, err := client.GetBotRandom(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)

	var data map[string]interface{}
	err = json.Unmarshal(result, &data)
	require.NoError(t, err)
	assert.Equal(t, "randombot", data["username"])
}

func ptr[T any](v T) *T {
	return &v
}
