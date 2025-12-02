package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gnomegl/funstat-api/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	apiKey  string
	debug   bool
)

var rootCmd = &cobra.Command{
	Use:   "funstat",
	Short: "Funstat API CLI client",
	Long:  `A command-line interface for interacting with the Funstat API.`,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.funstat.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key (overrides FUNSTAT_API_KEY env var)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")

	userCmd := &cobra.Command{
		Use:   "user",
		Short: "User-related operations",
		Long:  `Perform operations related to Telegram users.`,
	}

	resolveCmd := &cobra.Command{
		Use:   "resolve [usernames...]",
		Short: "Resolve Telegram usernames to user info (Cost: 0.10 per success)",
		Args:  cobra.MinimumNArgs(1),
		RunE:  resolveUsernames,
	}

	statsCmd := &cobra.Command{
		Use:   "stats [user-id]",
		Short: "Get full user statistics (Cost: 1)",
		Args:  cobra.ExactArgs(1),
		RunE:  getUserStats,
	}

	statsMinCmd := &cobra.Command{
		Use:   "stats-min [user-id]",
		Short: "Get minimal user statistics (FREE)",
		Args:  cobra.ExactArgs(1),
		RunE:  getUserStatsMin,
	}

	getByIDCmd := &cobra.Command{
		Use:   "get-by-id [user-ids...]",
		Short: "Get users by Telegram ID (Cost: 0.10 per success)",
		Args:  cobra.MinimumNArgs(1),
		RunE:  getUsersByID,
	}

	groupsCmd := &cobra.Command{
		Use:   "groups [user-id]",
		Short: "Get user's groups (Cost: 5)",
		Args:  cobra.ExactArgs(1),
		RunE:  getUserGroups,
	}

	groupsCountCmd := &cobra.Command{
		Use:   "groups-count [user-id]",
		Short: "Get count of user's groups (FREE)",
		Args:  cobra.ExactArgs(1),
		RunE:  getUserGroupsCount,
	}
	groupsCountCmd.Flags().Bool("only-with-messages", true, "Only count groups where user has messages")

	messagesCmd := &cobra.Command{
		Use:   "messages [user-id]",
		Short: "Get user messages (Cost: 10 per user if success)",
		Args:  cobra.ExactArgs(1),
		RunE:  getUserMessages,
	}
	messagesCmd.Flags().Int64("group-id", 0, "Filter by group ID")
	messagesCmd.Flags().String("text-contains", "", "Filter by message text")
	messagesCmd.Flags().Int32("media-code", 0, "Filter by media code")
	messagesCmd.Flags().Int32("page", 1, "Page number")
	messagesCmd.Flags().Int32("page-size", 10, "Page size")

	messagesCountCmd := &cobra.Command{
		Use:   "messages-count [user-id]",
		Short: "Get count of user's messages (FREE)",
		Args:  cobra.ExactArgs(1),
		RunE:  getUserMessagesCount,
	}

	namesCmd := &cobra.Command{
		Use:   "names [user-id]",
		Short: "Get user's name history (Cost: 3)",
		Args:  cobra.ExactArgs(1),
		RunE:  getUserNames,
	}

	usernamesCmd := &cobra.Command{
		Use:   "usernames [user-id]",
		Short: "Get user's @username history (Cost: 3)",
		Args:  cobra.ExactArgs(1),
		RunE:  getUserUsernames,
	}

	commonGroupsStatCmd := &cobra.Command{
		Use:   "common-groups-stat [user-id]",
		Short: "Get users who have common groups with specified user (Cost: 5)",
		Args:  cobra.ExactArgs(1),
		RunE:  getUserCommonGroupsStat,
	}

	usernameUsageCmd := &cobra.Command{
		Use:   "username-usage [username]",
		Short: "Search username usage (actual users, past usage, groups, mentions)",
		Args:  cobra.ExactArgs(1),
		RunE:  getUsernameUsage,
	}

	userCmd.AddCommand(resolveCmd, statsCmd, statsMinCmd, getByIDCmd,
		groupsCmd, groupsCountCmd, messagesCmd, messagesCountCmd,
		namesCmd, usernamesCmd, commonGroupsStatCmd, usernameUsageCmd)

	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Group-related operations",
		Long:  `Perform operations related to Telegram groups.`,
	}

	groupInfoCmd := &cobra.Command{
		Use:   "info [group-id]",
		Short: "Get group basic info (Cost: 0.01)",
		Args:  cobra.ExactArgs(1),
		RunE:  getGroupInfo,
	}

	commonGroupsCmd := &cobra.Command{
		Use:   "common [user-ids...]",
		Short: "Get common groups for specified users (Cost: 0.5)",
		Args:  cobra.MinimumNArgs(2),
		RunE:  getCommonGroups,
	}

	groupCmd.AddCommand(groupInfoCmd, commonGroupsCmd)

	textCmd := &cobra.Command{
		Use:   "text",
		Short: "Text-related operations",
		Long:  `Perform operations related to text search.`,
	}

	textSearchCmd := &cobra.Command{
		Use:   "search [text]",
		Short: "Search who and where wrote specified text (Cost: 0.1)",
		Args:  cobra.ExactArgs(1),
		RunE:  textSearch,
	}
	textSearchCmd.Flags().Int32("page", 1, "Page number")
	textSearchCmd.Flags().Int32("page-size", 10, "Page size")

	textCmd.AddCommand(textSearchCmd)

	rootCmd.AddCommand(userCmd, groupCmd, textCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigName(".funstat")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("FUNSTAT")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if apiKey == "" {
		apiKey = viper.GetString("API_KEY")
	}

	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: API key is required. Set FUNSTAT_API_KEY environment variable or use --api-key flag")
		os.Exit(1)
	}
}

func getClient() *client.Client {
	opts := []client.Option{}
	if debug {
		opts = append(opts, client.WithDebug(true))
	}
	return client.New(apiKey, opts...)
}

func resolveUsernames(cmd *cobra.Command, args []string) error {
	c := getClient()
	ctx := context.Background()

	result, err := c.ResolveUsernames(ctx, args)
	if err != nil {
		return err
	}

	return printJSON(result)
}

func getUserStats(cmd *cobra.Command, args []string) error {
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	c := getClient()
	ctx := context.Background()

	result, err := c.GetUserStats(ctx, userID)
	if err != nil {
		return err
	}

	return printJSON(result)
}

func getUserStatsMin(cmd *cobra.Command, args []string) error {
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	c := getClient()
	ctx := context.Background()

	result, err := c.GetUserStatsMin(ctx, userID)
	if err != nil {
		return err
	}

	return printJSON(result)
}

func getUsersByID(cmd *cobra.Command, args []string) error {
	var userIDs []int64
	for _, arg := range args {
		id, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid user ID %s: %w", arg, err)
		}
		userIDs = append(userIDs, id)
	}

	c := getClient()
	ctx := context.Background()

	result, err := c.GetUsersByID(ctx, userIDs)
	if err != nil {
		return err
	}

	return printJSON(result)
}

func getUserGroups(cmd *cobra.Command, args []string) error {
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	c := getClient()
	ctx := context.Background()

	result, err := c.GetUserGroups(ctx, userID)
	if err != nil {
		return err
	}

	return printJSON(result)
}

func getUserGroupsCount(cmd *cobra.Command, args []string) error {
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	onlyWithMessages, _ := cmd.Flags().GetBool("only-with-messages")

	c := getClient()
	ctx := context.Background()

	count, err := c.GetUserGroupsCount(ctx, userID, onlyWithMessages)
	if err != nil {
		return err
	}

	fmt.Printf("Groups count: %d\n", count)
	return nil
}

func getUserMessages(cmd *cobra.Command, args []string) error {
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	groupID, _ := cmd.Flags().GetInt64("group-id")
	textContains, _ := cmd.Flags().GetString("text-contains")
	mediaCode, _ := cmd.Flags().GetInt32("media-code")
	page, _ := cmd.Flags().GetInt32("page")
	pageSize, _ := cmd.Flags().GetInt32("page-size")

	opts := client.GetUserMessagesOptions{
		Page:     page,
		PageSize: pageSize,
	}

	if groupID != 0 {
		opts.GroupID = &groupID
	}
	if textContains != "" {
		opts.TextContains = &textContains
	}
	if mediaCode != 0 {
		opts.MediaCode = &mediaCode
	}

	c := getClient()
	ctx := context.Background()

	result, err := c.GetUserMessages(ctx, userID, opts)
	if err != nil {
		return err
	}

	return printJSON(result)
}

func getUserMessagesCount(cmd *cobra.Command, args []string) error {
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	c := getClient()
	ctx := context.Background()

	count, err := c.GetUserMessagesCount(ctx, userID)
	if err != nil {
		return err
	}

	fmt.Printf("Messages count: %d\n", count)
	return nil
}

func getUserNames(cmd *cobra.Command, args []string) error {
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	c := getClient()
	ctx := context.Background()

	result, err := c.GetUserNames(ctx, userID)
	if err != nil {
		return err
	}

	return printJSON(result)
}

func getUserUsernames(cmd *cobra.Command, args []string) error {
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	c := getClient()
	ctx := context.Background()

	result, err := c.GetUserUsernames(ctx, userID)
	if err != nil {
		return err
	}

	return printJSON(result)
}

func getGroupInfo(cmd *cobra.Command, args []string) error {
	groupID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	c := getClient()
	ctx := context.Background()

	result, err := c.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	var data interface{}
	if err := json.Unmarshal(result, &data); err != nil {
		return err
	}
	return printJSON(data)
}

func getUserCommonGroupsStat(cmd *cobra.Command, args []string) error {
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	c := getClient()
	ctx := context.Background()

	result, err := c.GetCommonGroupsStat(ctx, userID)
	if err != nil {
		return err
	}

	return printJSON(result)
}

func getUsernameUsage(cmd *cobra.Command, args []string) error {
	c := getClient()
	ctx := context.Background()

	result, err := c.GetUsernameUsage(ctx, args[0])
	if err != nil {
		return err
	}

	return printJSON(result)
}

func textSearch(cmd *cobra.Command, args []string) error {
	page, _ := cmd.Flags().GetInt32("page")
	pageSize, _ := cmd.Flags().GetInt32("page-size")

	c := getClient()
	ctx := context.Background()

	opts := &client.TextSearchOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := c.TextSearch(ctx, args[0], opts)
	if err != nil {
		return err
	}

	return printJSON(result)
}

func getCommonGroups(cmd *cobra.Command, args []string) error {
	var userIDs []int64
	for _, arg := range args {
		id, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid user ID %s: %w", arg, err)
		}
		userIDs = append(userIDs, id)
	}

	c := getClient()
	ctx := context.Background()

	result, err := c.GetCommonGroups(ctx, userIDs)
	if err != nil {
		return err
	}

	return printJSON(result)
}

func printJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
