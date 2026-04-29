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
		"data": [
			{
				"date": "2024-01-01T12:00:00Z",
				"messageId": 1,
				"text": "Hello",
				"group": {"id": 111, "title": "Test", "isPrivate": false}
			}
		]
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
			{"name": "John Doe", "date_time": "2023-01-01T00:00:00Z"},
			{"name": "John Smith", "date_time": "2023-06-02T00:00:00Z"}
		]
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetUserNames(ctx, 123456)
	require.NoError(t, err)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, "John Doe", result.Data[0].Name)
	assert.Equal(t, "John Smith", result.Data[1].Name)
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

func TestGetCommonGroupsStatHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 5, "current_ballance": 95, "request_duration": "80ms"},
		"data": [
			{"user_id": 222, "common_groups": 3, "first_name": "Alice", "is_user_active": true}
		]
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetCommonGroupsStat(ctx, 111)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, int32(3), result.Data[0].CommonGroups)
}

func TestGetUsernameUsageHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 0.1, "current_ballance": 99.9, "request_duration": "20ms"},
		"data": {
			"actualUsers": [{"id": 111, "username": "testuser", "is_active": true, "is_bot": false}],
			"usageByUsersInThePast": [],
			"actualGroupsOrChannels": [],
			"mentionByChannelOrGroupDesc": []
		}
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetUsernameUsage(ctx, "testuser")
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data.ActualUsers, 1)
}

func TestTextSearchHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 0.1, "current_ballance": 99.9, "request_duration": "30ms"},
		"data": {
			"total": 1,
			"data": [{"message_id": 100, "user_id": 111, "date": "2024-01-01T12:00:00Z", "text": "hello", "is_active": true, "group": {"id": 222, "title": "Test", "isPrivate": false, "isChannel": false}}],
			"isLastPage": true,
			"pageSize": 10,
			"currentPage": 1,
			"totalPages": 1,
			"isSliding": false
		}
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	opts := client.TextSearchOptions{Page: 1, PageSize: 10}
	result, err := c.TextSearch(ctx, "hello", opts)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, int32(1), result.Data.Total)
}

func TestGetCommonGroupsHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 0.5, "current_ballance": 99.5, "request_duration": "40ms"},
		"data": [
			{"id": 111, "title": "Shared Group", "isPrivate": false, "isChannel": false}
		]
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetCommonGroups(ctx, []int64{111, 222})
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data, 1)
}

func TestGetGroupMembersHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 15, "current_ballance": 85, "request_duration": "200ms"},
		"data": [
			{"id": 111, "first_name": "Alice", "is_active": true, "today_msg": 5, "has_photo": true},
			{"id": 222, "first_name": "Bob", "is_active": true, "today_msg": 0, "has_photo": false}
		]
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetGroupMembers(ctx, 999)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data, 2)
}

func TestGetUserStickersHandler(t *testing.T) {
	response := `{
		"success": true,
		"tech": {"request_cost": 1, "current_ballance": 99, "request_duration": "50ms"},
		"data": [
			{"sticker_set_id": 12345, "last_seen": "2024-01-01", "min_seen": "2023-01-01", "title": "My Pack", "short_name": "mypack", "stickers_count": 30}
		]
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetUserStickers(ctx, 111)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "My Pack", *result.Data[0].Title)
}

func TestGetGiftsRelationHandler(t *testing.T) {
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
				"to_is_active": true
			}
		]
	}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	opts := client.GiftsRelationOptions{Page: 1, PageSize: 20}
	result, err := c.GetGiftsRelation(ctx, 111, opts)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, int64(111), result.Data[0].FromUserID)
}

func TestGetUserReputationHandler(t *testing.T) {
	response := `{"spam_score": 0, "is_scammer": false}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetUserReputation(ctx, 111)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetBotRandomHandler(t *testing.T) {
	response := `{"user_id": 12345, "username": "randombot"}`

	server, c := setupTestServer(t, 200, response)
	defer server.Close()

	ctx := context.Background()
	result, err := c.GetBotRandom(ctx)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestErrorHandling(t *testing.T) {
	server, c := setupTestServer(t, 404, `{"title": "Not Found", "detail": "User not found"}`)
	defer server.Close()

	ctx := context.Background()
	_, err := c.GetUserStats(ctx, 999999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Not Found")
}
