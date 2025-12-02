package client

import (
	"time"
)

type AppProblem struct {
	Detail     *string    `json:"detail,omitempty"`
	Instance   *string    `json:"instance,omitempty"`
	Status     *int32     `json:"status,omitempty"`
	Title      *string    `json:"title,omitempty"`
	Type       *string    `json:"type,omitempty"`
	Method     *string    `json:"method,omitempty"`
	AppVersion *string    `json:"appVersion,omitempty"`
	DateTime   *time.Time `json:"dateTime,omitempty"`
}

type TechInfo struct {
	RequestCost     float32 `json:"request_cost"`
	CurrentBalance  float32 `json:"current_ballance"`
	RequestDuration string  `json:"request_duration"`
}

type ChatInfo struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	IsPrivate bool    `json:"isPrivate"`
	Username  *string `json:"username,omitempty"`
}

type ResolvedUser struct {
	ID         int64   `json:"id"`
	Username   *string `json:"username,omitempty"`
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	IsActive   bool    `json:"is_active"`
	IsBot      bool    `json:"is_bot"`
	HasPremium *bool   `json:"has_premium,omitempty"`
}

type ResolvedUserArrayAPIAnswer struct {
	Success bool           `json:"success"`
	Tech    TechInfo       `json:"tech"`
	Data    []ResolvedUser `json:"data,omitempty"`
}

type UserStatsMin struct {
	ID               int64      `json:"id"`
	FirstName        *string    `json:"first_name,omitempty"`
	LastName         *string    `json:"last_name,omitempty"`
	IsBot            bool       `json:"is_bot"`
	IsActive         bool       `json:"is_active"`
	FirstMsgDate     *time.Time `json:"first_msg_date,omitempty"`
	LastMsgDate      *time.Time `json:"last_msg_date,omitempty"`
	TotalMsgCount    int64      `json:"total_msg_count"`
	MsgInGroupsCount int64      `json:"msg_in_groups_count"`
	AdmInGroups      int32      `json:"adm_in_groups"`
	UsernamesCount   int32      `json:"usernames_count"`
	NamesCount       int32      `json:"names_count"`
	TotalGroups      int32      `json:"total_groups"`
}

type UserStatsMinAPIAnswer struct {
	Success bool          `json:"success"`
	Tech    TechInfo      `json:"tech"`
	Data    *UserStatsMin `json:"data,omitempty"`
}

type UserStats struct {
	ID                int64      `json:"id"`
	FirstName         *string    `json:"first_name,omitempty"`
	LastName          *string    `json:"last_name,omitempty"`
	IsBot             bool       `json:"is_bot"`
	IsActive          bool       `json:"is_active"`
	FirstMsgDate      *time.Time `json:"first_msg_date,omitempty"`
	LastMsgDate       *time.Time `json:"last_msg_date,omitempty"`
	TotalMsgCount     int64      `json:"total_msg_count"`
	MsgInGroupsCount  int64      `json:"msg_in_groups_count"`
	AdmInGroups       int32      `json:"adm_in_groups"`
	UsernamesCount    int32      `json:"usernames_count"`
	NamesCount        int32      `json:"names_count"`
	TotalGroups       int32      `json:"total_groups"`
	IsCyrillicPrimary *bool      `json:"is_cyrillic_primary,omitempty"`
	LangCode          *string    `json:"lang_code,omitempty"`
	UniquePercent     *float32   `json:"unique_percent,omitempty"`
	CircleCount       int32      `json:"circle_count"`
	VoiceCount        int32      `json:"voice_count"`
	ReplyPercent      float32    `json:"reply_percent"`
	MediaPercent      float32    `json:"media_percent"`
	LinkPercent       float32    `json:"link_percent"`
	FavoriteChat      *ChatInfo  `json:"favorite_chat,omitempty"`
	MediaUsage        *string    `json:"media_usage,omitempty"`
}

type UserStatsAPIAnswer struct {
	Success bool       `json:"success"`
	Tech    TechInfo   `json:"tech"`
	Data    *UserStats `json:"data,omitempty"`
}

type UserMsg struct {
	Date             time.Time `json:"date"`
	MessageID        int32     `json:"messageId"`
	ReplyToMessageID *int32    `json:"replyToMessageId,omitempty"`
	MediaCode        *int32    `json:"mediaCode,omitempty"`
	MediaName        *string   `json:"mediaName,omitempty"`
	Text             string    `json:"text"`
	Group            ChatInfo  `json:"group"`
}

type Paging struct {
	Total       int32 `json:"total"`
	CurrentPage int32 `json:"currentPage"`
	PageSize    int32 `json:"pageSize"`
	TotalPages  int32 `json:"totalPages"`
}

type UserMsgArrayAPIAnswerPaged struct {
	Success bool        `json:"success"`
	Tech    TechInfo    `json:"tech"`
	Paging  Paging      `json:"paging"`
	Data    [][]UserMsg `json:"data,omitempty"`
}

type UserChatInfo struct {
	Chat          *ChatInfo  `json:"chat,omitempty"`
	LastMessageID *int32     `json:"lastMessageId,omitempty"`
	MessagesCount *int32     `json:"messagesCount,omitempty"`
	LastMessage   *time.Time `json:"lastMessage,omitempty"`
	FirstMessage  *time.Time `json:"firstMessage,omitempty"`
	IsAdmin       bool       `json:"isAdmin"`
	IsLeft        bool       `json:"isLeft"`
}

type UserChatInfoArrayAPIAnswer struct {
	Success bool           `json:"success"`
	Tech    TechInfo       `json:"tech"`
	Data    []UserChatInfo `json:"data,omitempty"`
}

// TextSearchGroup represents group info in text search results
type TextSearchGroup struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	IsPrivate bool    `json:"isPrivate"`
	IsChannel bool    `json:"isChannel"`
	Username  *string `json:"username,omitempty"`
	Link      *string `json:"link,omitempty"`
}

// TextSearchResult represents a search result for text search
type TextSearchResult struct {
	MessageID int32           `json:"message_id"`
	UserID    int64           `json:"user_id"`
	Date      time.Time       `json:"date"`
	Name      *string         `json:"name,omitempty"`
	Username  *string         `json:"username,omitempty"`
	IsActive  bool            `json:"is_active"`
	Group     TextSearchGroup `json:"group"`
	Text      string          `json:"text"`
}

// TextSearchData represents the data object in text search response
type TextSearchData struct {
	IsLastPage  bool               `json:"isLastPage"`
	PageSize    int32              `json:"pageSize"`
	CurrentPage int32              `json:"currentPage"`
	TotalPages  int32              `json:"totalPages"`
	IsSliding   bool               `json:"isSliding"`
	Total       int32              `json:"total"`
	Data        []TextSearchResult `json:"data"`
}

// TextSearchAPIAnswer is the response for text search endpoint
type TextSearchAPIAnswer struct {
	Success bool           `json:"success"`
	Tech    TechInfo       `json:"tech"`
	Data    TextSearchData `json:"data"`
}

// CommonGroupsAPIAnswer is the response for common groups endpoint
type CommonGroupsAPIAnswer struct {
	Success bool       `json:"success"`
	Tech    TechInfo   `json:"tech"`
	Data    []ChatInfo `json:"data,omitempty"`
}

// CommonGroupsStatUser represents a user with common groups
type CommonGroupsStatUser struct {
	UserID       int64   `json:"user_id"`
	Username     *string `json:"username,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	CommonGroups int32   `json:"common_groups"`
}

// CommonGroupsStatAPIAnswer is the response for common groups stat endpoint
type CommonGroupsStatAPIAnswer struct {
	Success bool                   `json:"success"`
	Tech    TechInfo               `json:"tech"`
	Data    []CommonGroupsStatUser `json:"data,omitempty"`
}

// UsernameUsageResult represents a username usage result
type UsernameUsageResult struct {
	// Type: 1=actual user, 2=past user usage, 3=group/channel, 4=mentioned in description
	Type        int32      `json:"type"`
	UserID      *int64     `json:"user_id,omitempty"`
	GroupID     *int64     `json:"group_id,omitempty"`
	Username    string     `json:"username"`
	Title       *string    `json:"title,omitempty"`
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
	DateFrom    *time.Time `json:"date_from,omitempty"`
	DateTo      *time.Time `json:"date_to,omitempty"`
	Description *string    `json:"description,omitempty"`
}

// UsernameUsageAPIAnswer is the response for username usage endpoint
type UsernameUsageAPIAnswer struct {
	Success bool                  `json:"success"`
	Tech    TechInfo              `json:"tech"`
	Data    []UsernameUsageResult `json:"data,omitempty"`
}
