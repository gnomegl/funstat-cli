package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gnomegl/funstat-cli/pkg/client"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T, statusCode int, response string) (*httptest.Server, *client.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	}))

	c := client.New("test-key", client.WithBaseURL(server.URL))
	return server, c
}

func TestResolveUsernamesHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 0.1, "current_ballance": 100, "request_duration": "10ms"},
		"data": [{"id": 123456, "username": "testuser", "is_active": true, "is_bot": false}]
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.ResolveUsernames(ctx, []string{"testuser"})
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data, 1)
}

func TestGetUserStatsHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 1, "current_ballance": 99, "request_duration": "50ms"},
		"data": {
			"id": 123456,
			"first_name": "John",
			"is_bot": false,
			"is_active": true,
			"total_msg_count": 1000
		}
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetUserStats(ctx, 123456)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, int64(123456), result.Data.ID)
}

func TestGetUserStatsMinHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 0, "current_ballance": 100, "request_duration": "10ms"},
		"data": {
			"id": 123456,
			"is_bot": false,
			"is_active": true,
			"total_msg_count": 500
		}
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetUserStatsMin(ctx, 123456)
	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestGetUsersByIDHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 0.2, "current_ballance": 99.8, "request_duration": "20ms"},
		"data": [
			{"id": 111, "is_active": true, "is_bot": false},
			{"id": 222, "is_active": true, "is_bot": false}
		]
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetUsersByID(ctx, []int64{111, 222})
	require.NoError(t, err)
	assert.Len(t, result.Data, 2)
}

func TestGetUserGroupsHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 5, "current_ballance": 95, "request_duration": "100ms"},
		"data": [
			{
				"chat": {"id": 111, "title": "Test Group", "isPrivate": false},
				"messagesCount": 50,
				"isAdmin": true,
				"isLeft": false
			}
		]
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetUserGroups(ctx, 123456)
	require.NoError(t, err)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "Test Group", result.Data[0].Chat.Title)
}

func TestGetUserGroupsCountHandler(t *testing.T) {
	server, c := setupTestServer(t, 200, `25`)
	defer server.Close()

	ctx := context.Background()
	count, err := c.GetUserGroupsCount(ctx, 123456, false)
	require.NoError(t, err)
	assert.Equal(t, int32(25), count)
}

func TestGetUserMessagesHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 10, "current_ballance": 90, "request_duration": "200ms"},
		"paging": {"total": 100, "currentPage": 1, "pageSize": 10, "totalPages": 10},
		"data": [[
			{
				"date": "2024-01-01T12:00:00Z",
				"messageId": 1,
				"text": "Hello",
				"group": {"id": 111, "title": "Test", "isPrivate": false}
			}
		]]
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	opts := client.GetUserMessagesOptions{
		Page:     1,
		PageSize: 10,
	}
	result, err := c.GetUserMessages(ctx, 123456, opts)
	require.NoError(t, err)
	assert.Equal(t, int32(100), result.Paging.Total)
}

func TestGetUserMessagesCountHandler(t *testing.T) {
	server, c := setupTestServer(t, 200, `500`)
	defer server.Close()

	ctx := context.Background()
	count, err := c.GetUserMessagesCount(ctx, 123456)
	require.NoError(t, err)
	assert.Equal(t, int32(500), count)
}

func TestGetUserNamesHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 3, "current_ballance": 97, "request_duration": "50ms"},
		"data": [
			{"firstMessage": "2023-01-01T00:00:00Z"},
			{"firstMessage": "2023-06-02T00:00:00Z"}
		]
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetUserNames(ctx, 123456)
	require.NoError(t, err)
	assert.Len(t, result.Data, 2)
}

func TestGetUserUsernamesHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 3, "current_ballance": 97, "request_duration": "50ms"},
		"data": [
			{"name": "oldusername1", "date_time": "2023-01-01T00:00:00Z"},
			{"name": "oldusername2", "date_time": "2023-06-02T00:00:00Z"}
		]
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetUserUsernames(ctx, 123456)
	require.NoError(t, err)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, "oldusername1", result.Data[0].Name)
}

func TestGetGroupHandler(t *testing.T) {
	response := `{
		"id": 111222333,
		"title": "Test Group",
		"username": "testgroup",
		"members_count": 1000
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetGroup(ctx, 111222333)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCobraCommandStructure(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *cobra.Command
		wantUse string
		minArgs int
	}{
		{"root command", rootCmd, "funstat", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantUse, tt.cmd.Use)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	server, c := setupTestServer(t, 404, `{"title": "Not Found", "detail": "User not found"}`)
	defer server.Close()

	ctx := context.Background()
	_, err := c.GetUserStats(ctx, 999999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Not Found")
}
