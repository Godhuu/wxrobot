package wxrobot

type EventType string
type ChatType string
type MsgType string

const (
	AddToChatEvent      EventType = "add_to_chat"      // 机器人被添加进会话
	DeleteFromChatEvent EventType = "delete_from_chat" // 机器人被移出会话
	EnterChatEvent      EventType = "enter_chat"       // 用户进入机器人单聊
)

const (
	ChatTypeSingle          ChatType = "single"           // 单聊
	ChatTypeGroup           ChatType = "group"            // 群聊
	ChatTypeBlackboard      ChatType = "blackboard"       // 小黑板帖子
	ChatTypeBlackboardReply ChatType = "blackboard_reply" // 小黑板帖子回复
)

const (
	MsgTypeEvent MsgType = "event" // 事件消息
	/*
		attachment事件消息 机器人可以通过接口发送带attachment的markdown消息，目前attachment
		支持按钮类型 当用户点击按钮时，企业微信往机器人回调相应的事件
	*/
	MsgTypeAttachment MsgType = "attachment" //

	MsgTypeMixed MsgType = "mixed" // 图文混排消息
	MsgTypeText  MsgType = "text"  // 文本消息
	MsgTypeImage MsgType = "image" // 图片消息 目前仅支持 chatTypeSingle 会话 ｜注意：目前仅单聊支持回调图片消息
)
