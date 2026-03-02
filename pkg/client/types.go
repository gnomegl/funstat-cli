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

// ChatInfo represents basic group/channel info (chat_inf)
type ChatInfo struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	IsPrivate bool    `json:"isPrivate"`
	Username  *string `json:"username,omitempty"`
}

// ChatInfoExt represents extended group/channel info (chat_inf_ext)
type ChatInfoExt struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	IsPrivate bool    `json:"isPrivate"`
	IsChannel bool    `json:"isChannel"`
	Username  *string `json:"username,omitempty"`
	Link      *string `json:"link,omitempty"`
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
	StarsVal          *int64     `json:"stars_val,omitempty"`
	PersonalChannelID *int64     `json:"personal_channel_id,omitempty"`
	GiftCount         *int32     `json:"gift_count,omitempty"`
	StarsLevel        *int32     `json:"stars_level,omitempty"`
	BirthDay          *int32     `json:"birth_day,omitempty"`
	BirthMonth        *int32     `json:"birth_month,omitempty"`
	BirthYear         *int32     `json:"birth_year,omitempty"`
	About             *string    `json:"about,omitempty"`
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
	Success bool      `json:"success"`
	Tech    TechInfo  `json:"tech"`
	Paging  Paging    `json:"paging"`
	Data    []UserMsg `json:"data,omitempty"`
}

type UserChatInfo struct {
	Chat          *ChatInfo  `json:"chat"`
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

// UsernameHistoryItem represents a name/username history entry (user_name_inf)
type UsernameHistoryItem struct {
	Name     string    `json:"name"`
	DateTime time.Time `json:"date_time"`
}

type UsernameHistoryAPIAnswer struct {
	Success bool                  `json:"success"`
	Tech    TechInfo              `json:"tech"`
	Data    []UsernameHistoryItem `json:"data,omitempty"`
}

// TextSearchResult represents a text search result (who_wrote_text)
type TextSearchResult struct {
	MessageID int32       `json:"message_id"`
	UserID    int64       `json:"user_id"`
	Date      time.Time   `json:"date"`
	Name      *string     `json:"name,omitempty"`
	Username  *string     `json:"username,omitempty"`
	IsActive  bool        `json:"is_active"`
	Group     ChatInfoExt `json:"group"`
	Text      string      `json:"text"`
}

// TextSearchData represents the paginated text search data (who_wrote_textPagedNoCount)
type TextSearchData struct {
	Total       int32              `json:"total"`
	Data        []TextSearchResult `json:"data"`
	IsLastPage  bool               `json:"isLastPage"`
	PageSize    int32              `json:"pageSize"`
	CurrentPage int32              `json:"currentPage"`
	TotalPages  int32              `json:"totalPages"`
	IsSliding   bool               `json:"isSliding"`
}

type TextSearchAPIAnswer struct {
	Success bool           `json:"success"`
	Tech    TechInfo       `json:"tech"`
	Data    TextSearchData `json:"data"`
}

// CommonGroupsAPIAnswer is the response for common groups endpoint
type CommonGroupsAPIAnswer struct {
	Success bool          `json:"success"`
	Tech    TechInfo      `json:"tech"`
	Data    []ChatInfoExt `json:"data,omitempty"`
}

// CommonGroupsStatUser represents a user with common groups (ucommon_group_inf)
type CommonGroupsStatUser struct {
	UserID       int64   `json:"user_id"`
	CommonGroups int32   `json:"common_groups"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Username     *string `json:"username,omitempty"`
	IsUserActive bool    `json:"is_user_active"`
}

type CommonGroupsStatAPIAnswer struct {
	Success bool                   `json:"success"`
	Tech    TechInfo               `json:"tech"`
	Data    []CommonGroupsStatUser `json:"data,omitempty"`
}

// UsernameUsageModel represents categorized username usage results (username_usage_model)
type UsernameUsageModel struct {
	ActualUsers                []ResolvedUser `json:"actualUsers,omitempty"`
	UsageByUsersInThePast      []ResolvedUser `json:"usageByUsersInThePast,omitempty"`
	ActualGroupsOrChannels     []ChatInfoExt  `json:"actualGroupsOrChannels,omitempty"`
	MentionByChannelOrGroupDesc []ChatInfoExt  `json:"mentionByChannelOrGroupDesc,omitempty"`
}

type UsernameUsageAPIAnswer struct {
	Success bool               `json:"success"`
	Tech    TechInfo           `json:"tech"`
	Data    UsernameUsageModel `json:"data"`
}

// GroupMember represents a group member (group_member)
type GroupMember struct {
	ID        int64   `json:"id"`
	Username  *string `json:"username,omitempty"`
	Name      *string `json:"name,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	IsAdmin   *bool   `json:"is_admin,omitempty"`
	IsActive  bool    `json:"is_active"`
	TodayMsg  int32   `json:"today_msg"`
	HasPrem   *bool   `json:"has_prem,omitempty"`
	HasPhoto  bool    `json:"has_photo"`
	DcID      *int32  `json:"dc_id,omitempty"`
}

type GroupMemberArrayAPIAnswer struct {
	Success bool          `json:"success"`
	Tech    TechInfo      `json:"tech"`
	Data    []GroupMember `json:"data,omitempty"`
}

// GiftRelation represents a gift relationship between users (gift_relations_inf)
type GiftRelation struct {
	LastGiftDate     *time.Time `json:"last_gift_date,omitempty"`
	FromUserID       int64      `json:"from_user_id"`
	FromFirstName    *string    `json:"from_first_name,omitempty"`
	FromLastName     *string    `json:"from_last_name,omitempty"`
	FromMainUsername *string    `json:"from_mainUsername,omitempty"`
	FromIsActive     bool       `json:"from_is_active"`
	ToUserID         int64      `json:"to_user_id"`
	ToFirstName      *string    `json:"to_first_name,omitempty"`
	ToLastName       *string    `json:"to_last_name,omitempty"`
	ToMainUsername   *string    `json:"to_mainUsername,omitempty"`
	ToIsActive       bool       `json:"to_is_active"`
}

type GiftRelationArrayAPIAnswer struct {
	Success bool           `json:"success"`
	Tech    TechInfo       `json:"tech"`
	Data    []GiftRelation `json:"data,omitempty"`
}

// StickerInfo represents a sticker pack created by a user (sticker_inf)
type StickerInfo struct {
	StickerSetID  int64      `json:"sticker_set_id"`
	LastSeen      string     `json:"last_seen"`
	MinSeen       string     `json:"min_seen"`
	Resolved      *time.Time `json:"resolved,omitempty"`
	Title         *string    `json:"title,omitempty"`
	ShortName     *string    `json:"short_name,omitempty"`
	StickersCount *int32     `json:"stickers_count,omitempty"`
}

type StickerArrayAPIAnswer struct {
	Success bool          `json:"success"`
	Tech    TechInfo      `json:"tech"`
	Data    []StickerInfo `json:"data,omitempty"`
}
