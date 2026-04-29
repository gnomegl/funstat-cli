package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gnomegl/funstat-cli/pkg/client"
)

func main() {
	apiKey := os.Getenv("FUNSTAT_API_KEY")
	if apiKey == "" {
		log.Fatal("FUNSTAT_API_KEY environment variable is required")
	}

	c := client.New(apiKey, client.WithDebug(false))
	ctx := context.Background()

	resolvedUsers, err := c.ResolveUsernames(ctx, []string{"durov", "telegram"})
	if err != nil {
		log.Printf("Error resolving usernames: %v", err)
	} else {
		fmt.Printf("Resolved %d users\n", len(resolvedUsers.Data))
		fmt.Printf("Request cost: %.2f\n", resolvedUsers.Tech.RequestCost)
		fmt.Printf("Current balance: %.2f\n", resolvedUsers.Tech.CurrentBalance)

		for _, user := range resolvedUsers.Data {
			fmt.Printf("  User: %s %s (@%s) - ID: %d, Active: %v\n",
				strPtr(user.FirstName), strPtr(user.LastName),
				strPtr(user.Username), user.ID, user.IsActive)
		}
	}

	if len(resolvedUsers.Data) > 0 {
		userID := resolvedUsers.Data[0].ID
		statsMin, err := c.GetUserStatsMin(ctx, userID)
		if err != nil {
			log.Printf("Error getting minimal stats: %v", err)
		} else {
			if statsMin.Data != nil {
				fmt.Printf("User: %s %s\n",
					strPtr(statsMin.Data.FirstName),
					strPtr(statsMin.Data.LastName))
				fmt.Printf("  Total messages: %d\n", statsMin.Data.TotalMsgCount)
				fmt.Printf("  Groups with messages: %d\n", statsMin.Data.MsgInGroupsCount)
				fmt.Printf("  Total groups: %d\n", statsMin.Data.TotalGroups)
				fmt.Printf("  Admin in groups: %d\n", statsMin.Data.AdmInGroups)
			}
		}
	}

	if len(resolvedUsers.Data) > 0 {
		userID := resolvedUsers.Data[0].ID
		count, err := c.GetUserGroupsCount(ctx, userID, true)
		if err != nil {
			log.Printf("Error getting groups count: %v", err)
		} else {
			fmt.Printf("User is in %d groups (with messages)\n", count)
		}
	}

	if len(resolvedUsers.Data) > 0 {
		userID := resolvedUsers.Data[0].ID
		count, err := c.GetUserMessagesCount(ctx, userID)
		if err != nil {
			log.Printf("Error getting messages count: %v", err)
		} else {
			fmt.Printf("User has %d total messages\n", count)
		}
	}

	usersByID, err := c.GetUsersByID(ctx, []int64{1, 777000})
	if err != nil {
		log.Printf("Error getting users by ID: %v", err)
	} else {
		fmt.Printf("Found %d users by ID\n", len(usersByID.Data))
		for _, user := range usersByID.Data {
			fmt.Printf("  User ID %d: %s %s (@%s)\n",
				user.ID, strPtr(user.FirstName),
				strPtr(user.LastName), strPtr(user.Username))
		}
	}
}

func strPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}


