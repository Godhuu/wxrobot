package wxrobot

import "encoding/xml"

// from 来自
type from struct {
	UserId string `xml:"UserId"` // 发送者的userid
	Name   string `xml:"Name"`   // 发送者姓名 公司内为中文名
	Alias  string `xml:"Alias"`  // 发送者别名 公司内为rtx英文名
}

//
type fromText struct {
	Content string `xml:"Content"` // 消息内容
}

//
type fromImage struct {
	ImageUrl string `xml:"ImageUrl"` // 图片的url，注意不可在网页引用该图片
}

// FromCommonMsg 基础消息结构体
type FromCommonMsg struct {
	bot   *bot
	From  from   `xml:"From"`  // 发送者信息
	MsgId string `xml:"MsgId"` // 消息Id，可用于去重
	/*
		会话类型，single，group，blackboard和blackboard_reply，分别表示：单聊，群聊，小黑板帖子和小黑板帖子回复，目前仅单聊支持回调图片
	*/
	ChatType       string `xml:"ChatType"`
	MsgType        string `xml:"MsgType"`        // 消息类型
	WebhookUrl     string `xml:"WebhookUrl"`     // 机器人主动推送消息的url
	ChatId         string `xml:"ChatId"`         // 会话id，可能是群聊，也可能是单聊，也可能是小黑板
	GetChatInfoUrl string `xml:"GetChatInfoUrl"` // 获取群信息的URL，有效时间5分钟，且仅能调用一次，当ChatType是single时不提供该字段。
}

func (r *FromCommonMsg) ToTextMsg(msg string) *toMsgText {
	text := new(toMsgText)
	text.bot = r.bot
	text.ChatId(r.ChatId)
	text.msgType = "text"
	text.Content = msg
	return text
}

func (r *FromCommonMsg) ToMarkdownMsg(markdown string) *toMsgMarkdown {
	t := new(toMsgMarkdown)
	t.bot = r.bot
	t.ChatId(r.ChatId)
	t.msgType = "markdown"
	t.Content = markdown
	return t
}

func (r *FromCommonMsg) ToImageMsg() *toMsgImage {
	t := new(toMsgImage)
	t.bot = r.bot
	t.ChatId(r.ChatId)
	t.msgType = "image"
	return t
}

func (r *FromCommonMsg) ToNewsMsg() *toMsgNews {
	t := new(toMsgNews)
	t.bot = r.bot
	t.ChatId(r.ChatId)
	t.msgType = "news"
	return t
}

func (r *FromCommonMsg) ToFileMsg() *toMsgFile {
	t := new(toMsgFile)
	t.bot = r.bot
	t.ChatId(r.ChatId)
	t.msgType = "file"
	return t
}

func (r *FromCommonMsg) GetChatType() ChatType {
	return ChatType(r.ChatType)
}

func (r *FromCommonMsg) GetMsgType() MsgType {
	return MsgType(r.MsgType)
}

// FromTextMsg 文本消息
type FromTextMsg struct {
	FromCommonMsg
	Text   fromText `xml:"Text"`
	PostId string   `xml:"PostId"` // 小黑板帖子id，当前消息为小黑板回帖消息时带上
}

// FromImageMsg 图片消息
type FromImageMsg struct {
	FromCommonMsg
	Image fromImage `xml:"Image"`
}

// FromEventMsg 事件消息
type FromEventMsg struct {
	FromCommonMsg
	Event      Event  `xml:"Event"`
	AppVersion string `xml:"AppVersion"` // 客户端版本号，当ChatType为blackboard时不提供该字段
}

func (fm *FromEventMsg) EventType() EventType {
	return EventType(fm.Event.EventType)
}

// FromAttachmentMsg 事件消息
//
// 机器人可以通过接口发送带attachment的markdown消息，目前attachment支持按钮类型，当用户点击按钮时，企业微信往机器人回调相应的事件
type FromAttachmentMsg struct {
	FromCommonMsg
	PostId     string        `xml:"PostId"`     // 小黑板帖子id，当前消息为小黑板回帖消息时带上
	Attachment MsgAttachment `xml:"Attachment"` // 用户点击的attachment，目前只支持button
}

// FromMixedMsg 图文混排消息
type FromMixedMsg struct {
	FromCommonMsg
	MixedMessage []MsgItem `xml:"MixedMessage>MsgItem"` // 图文混排消息，可由多个MsgItem组成
}

// Event
type Event struct {
	// 事件类型:目前可能是add_to_chat表示被添加进会话,或者delete_from_chat表示被移出会话,enter_chat 表示用户进入机器人单聊
	EventType string `xml:"EventType"`
}

// MsgItem
type MsgItem struct {
	XMLName xml.Name  `xml:"MsgItem"`
	MsgType string    `xml:"MsgType"` // 消息类型 text 文本 image 图片
	Text    fromText  `xml:"Text"`    // 消息类型 text 时存在
	Image   fromImage `xml:"Image"`   // 消息类型 image 时存在
}
