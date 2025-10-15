package client

import (
	"time"
)

// AppProblem represents an API error response
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

// TechInfo contains technical information about the request
type TechInfo struct {
	RequestCost     float32 `json:"request_cost"`
	CurrentBalance  float32 `json:"current_ballance"`
	RequestDuration string  `json:"request_duration"`
}

// ChatInfo represents chat/group information
type ChatInfo struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	IsPrivate bool    `json:"isPrivate"`
	Username  *string `json:"username,omitempty"`
}

// ResolvedUser represents resolved user information
type ResolvedUser struct {
	ID         int64   `json:"id"`
	Username   *string `json:"username,omitempty"`
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	IsActive   bool    `json:"is_active"`
	IsBot      bool    `json:"is_bot"`
	HasPremium *bool   `json:"has_premium,omitempty"`
}

// ResolvedUserArrayAPIAnswer represents the API response for resolved users
type ResolvedUserArrayAPIAnswer struct {
	Success bool           `json:"success"`
	Tech    TechInfo       `json:"tech"`
	Data    []ResolvedUser `json:"data,omitempty"`
}

// UserStatsMin represents minimal user statistics
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

// UserStatsMinAPIAnswer represents the API response for minimal user stats
type UserStatsMinAPIAnswer struct {
	Success bool          `json:"success"`
	Tech    TechInfo      `json:"tech"`
	Data    *UserStatsMin `json:"data,omitempty"`
}

// UserStats represents full user statistics
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

// UserStatsAPIAnswer represents the API response for full user stats
type UserStatsAPIAnswer struct {
	Success bool       `json:"success"`
	Tech    TechInfo   `json:"tech"`
	Data    *UserStats `json:"data,omitempty"`
}

// UserMsg represents a user message
type UserMsg struct {
	Date             time.Time `json:"date"`
	MessageID        int32     `json:"messageId"`
	ReplyToMessageID *int32    `json:"replyToMessageId,omitempty"`
	MediaCode        *int32    `json:"mediaCode,omitempty"`
	MediaName        *string   `json:"mediaName,omitempty"`
	Text             string    `json:"text"`
	Group            ChatInfo  `json:"group"`
}

// Paging represents pagination information
type Paging struct {
	Total       int32 `json:"total"`
	CurrentPage int32 `json:"currentPage"`
	PageSize    int32 `json:"pageSize"`
	TotalPages  int32 `json:"totalPages"`
}

// UserMsgArrayAPIAnswerPaged represents the paginated API response for user messages
type UserMsgArrayAPIAnswerPaged struct {
	Success bool        `json:"success"`
	Tech    TechInfo    `json:"tech"`
	Paging  Paging      `json:"paging"`
	Data    [][]UserMsg `json:"data,omitempty"`
}

// UserChatInfo represents user chat information
type UserChatInfo struct {
	Chat          *ChatInfo  `json:"chat,omitempty"`
	LastMessageID *int32     `json:"lastMessageId,omitempty"`
	MessagesCount *int32     `json:"messagesCount,omitempty"`
	LastMessage   *time.Time `json:"lastMessage,omitempty"`
	FirstMessage  *time.Time `json:"firstMessage,omitempty"`
	IsAdmin       bool       `json:"isAdmin"`
	IsLeft        bool       `json:"isLeft"`
}

// UserChatInfoArrayAPIAnswer represents the API response for user chat info
type UserChatInfoArrayAPIAnswer struct {
	Success bool           `json:"success"`
	Tech    TechInfo       `json:"tech"`
	Data    []UserChatInfo `json:"data,omitempty"`
}
