//go:build integration
// +build integration

package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gnomegl/funstat-api/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getIntegrationClient(t *testing.T) *client.Client {
	apiKey := os.Getenv("FUNSTAT_API_KEY")
	if apiKey == "" {
		t.Skip("FUNSTAT_API_KEY not set, skipping integration tests")
	}
	return client.New(apiKey)
}

func TestIntegrationResolveUsernames(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	c := getIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := c.ResolveUsernames(ctx, []string{"telegram"})
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotEmpty(t, result.Data)
}

func TestIntegrationGetUserStatsMin(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	c := getIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testUserID := int64(777000)

	result, err := c.GetUserStatsMin(ctx, testUserID)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Data)
}

func TestIntegrationGetUserStats(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	c := getIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testUserID := int64(777000)

	result, err := c.GetUserStats(ctx, testUserID)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Data)
	assert.Equal(t, testUserID, result.Data.ID)
}

func TestIntegrationGetUsersByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	c := getIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := c.GetUsersByID(ctx, []int64{777000})
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotEmpty(t, result.Data)
}

func TestIntegrationGetUserGroupsCount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	c := getIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testUserID := int64(777000)

	count, err := c.GetUserGroupsCount(ctx, testUserID, true)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int32(0))
}

func TestIntegrationGetUserMessagesCount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	c := getIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testUserID := int64(777000)

	count, err := c.GetUserMessagesCount(ctx, testUserID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int32(0))
}

func TestIntegrationRateLimiting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	c := getIntegrationClient(t)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_, err := c.GetUserStatsMin(ctx, int64(777000+i))
		require.NoError(t, err)
		time.Sleep(100 * time.Millisecond)
	}
}

func TestIntegrationErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	c := getIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	invalidUserID := int64(999999999999)
	_, err := c.GetUserStats(ctx, invalidUserID)
	assert.Error(t, err)
}

func TestIntegrationContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	c := getIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond)

	_, err := c.GetUserStats(ctx, 777000)
	assert.Error(t, err)
}
