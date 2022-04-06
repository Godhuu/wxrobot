package test

import (
	_ "embed" //embed
	"fmt"
	"git.woa.com/mingkunhu/wxrobot"

	"net/http"
	"testing"
)

// 初始化自己的机器人配置 要配置多个机器人只需要多次调用以下方法即可
var bot = wxrobot.Bot("demo").WebhookURL(
	"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=****").
	Serve("***", "***")

func init() {
	//wxrobot.SetLogger(myLogger) //设置自己的日志记录器，默认日志打印在控制台
}

//------------------- 以下测试及提供示例 机器人提供交互及发送消息的方法 --------------

func TestBot(t *testing.T) {
	bot.SwitchDebugMode(true)

	//---- 注册各种事件消息回调处理函数 -------
	// 注册发送给我的文本消息回调处理函数
	bot.RegisterHandlerForText(processTextMsgCallback)
	// 注册发送给我的图文消息回调处理函数
	bot.RegisterHandlerForMixed(processMixedMsgCallback)
	// 注册发送给我的图片消息回调处理函数 注意：图片消息目前仅支持私聊我的消息
	bot.RegisterHandlerForImage(processImageMsgCallback)
	// 注册发送给我的markdown消息点击事件回调处理函数
	bot.RegisterHandlerForAttachment(processAttachmentMsgCallback)
	// 注册关于我的事件消息回调处理函数
	bot.RegisterHandlerForEvent(processEventMsgCallback)

	if err := http.ListenAndServe(":8000", bot.NewRouter()); err != nil {
		panic(err)
	}
}

func processTextMsgCallback(msg *wxrobot.FromTextMsg) {
	fmt.Println("开始进行业务处理...")
	switch msg.GetChatType() {
	case wxrobot.ChatTypeGroup:
		fmt.Println("这是群聊中有人@我的消息。")
	case wxrobot.ChatTypeSingle:
		fmt.Println("这是私聊中有人发给我的消息。")
	case wxrobot.ChatTypeBlackboard:
		fmt.Println("这是小黑板帖子消息。")
	case wxrobot.ChatTypeBlackboardReply:
		fmt.Println("这是小黑板帖子有人回复我的消息。")
	}
	fmt.Println("业务处理完成。")

	//如果需要回复必要信息 默认直接回复消息发送人，无需指定chatid
	_ = msg.ToTextMsg(fmt.Sprintf("消息[%s]已收到！", msg.Text.Content)).Send()

	// --------------------------
	// 同样可以回复其他任何消息类型
	//_ = msg.ToNewsMsg().Send()
	//_ = msg.ToImageMsg().Send()
	// 。。。。
}

func processMixedMsgCallback(msg *wxrobot.FromMixedMsg) {
	fmt.Println("开始进行业务处理...")
	switch msg.GetChatType() {
	case wxrobot.ChatTypeGroup:
		fmt.Println("这是群聊中有人@我的消息。")
	case wxrobot.ChatTypeSingle:
		fmt.Println("这是私聊中有人发给我的消息。")
	case wxrobot.ChatTypeBlackboard:
		fmt.Println("这是小黑板帖子消息。")
	case wxrobot.ChatTypeBlackboardReply:
		fmt.Println("这是小黑板帖子有人回复我的消息。")
	}
	fmt.Println("业务处理完成。")

	//如果需要回复必要信息 默认直接回复消息发送人，无需指定chatid
	_ = msg.ToTextMsg("消息已收到！").Send()

	// --------------------------
	// 同样可以回复其他任何消息类型
	//_ = msg.ToNewsMsg().Send()
	//_ = msg.ToImageMsg().Send()
	// 。。。。
}

func processImageMsgCallback(msg *wxrobot.FromImageMsg) {
	fmt.Println("开始进行业务处理...")
	fmt.Println("图片消息回调当前仅支持私聊我的消息。")
	//switch msg.GetChatType() {
	//case wxrobot.ChatTypeGroup:
	//	fmt.Println("这是群聊中有人@我的消息。")
	//case wxrobot.ChatTypeSingle:
	//	fmt.Println("这是私聊中有人发给我的消息。")
	//case wxrobot.ChatTypeBlackboard:
	//	fmt.Println("这是小黑板帖子消息。")
	//case wxrobot.ChatTypeBlackboardReply:
	//	fmt.Println("这是小黑板帖子有人回复我的消息。")
	//}
	fmt.Println("业务处理完成。")

	//如果需要回复必要信息 默认直接回复消息发送人，无需指定chatid
	_ = msg.ToTextMsg("图片消息已收到！").Send()

	// --------------------------
	// 同样可以回复其他任何消息类型
	//_ = msg.ToNewsMsg().Send()
	//_ = msg.ToImageMsg().Send()
	// 。。。。
}

func processAttachmentMsgCallback(msg *wxrobot.FromAttachmentMsg) {
	fmt.Println("开始进行业务处理...")
	switch msg.GetChatType() {
	case wxrobot.ChatTypeGroup:
		fmt.Println("这是群聊中有人@我的消息。")
	case wxrobot.ChatTypeSingle:
		fmt.Println("这是私聊中有人发给我的消息。")
	case wxrobot.ChatTypeBlackboard:
		fmt.Println("这是小黑板帖子消息。")
	case wxrobot.ChatTypeBlackboardReply:
		fmt.Println("这是小黑板帖子有人回复我的消息。")
	}

	//如果markdown中有按钮选项等供用户点击，此处可获取用户点击的内容
	for _, action := range msg.Attachment.Actions {
		fmt.Println("用户", msg.From.Name, "选中的是", action.Name, "value:", action.Value)
	}
	fmt.Println("业务处理完成。")

	//如果需要回复必要信息 默认直接回复消息发送人，无需指定chatid
	_ = msg.ToTextMsg("markdown消息已收到！").Send()

	// --------------------------
	// 同样可以回复其他任何消息类型
	//_ = msg.ToNewsMsg().Send()
	//_ = msg.ToImageMsg().Send()
	// 。。。。
}

func processEventMsgCallback(msg *wxrobot.FromEventMsg) {
	fmt.Println("开始进行业务处理...")
	switch msg.EventType() {
	case wxrobot.AddToChatEvent:
		fmt.Println("我被添加进某个群了。")
	case wxrobot.DeleteFromChatEvent:
		fmt.Println("我从群聊中被管理员删除了。")
	case wxrobot.EnterChatEvent:
		fmt.Println("用户", msg.From.Name, "开始跟我私聊了。")
	}
	fmt.Println("业务处理完成。")

	//如果需要回复必要信息 默认直接回复消息发送人，无需指定chatid
	_ = msg.ToTextMsg("事件消息已收到！").Send()

	// --------------------------
	// 同样可以回复其他任何消息类型
	//_ = msg.ToNewsMsg().Send()
	//_ = msg.ToImageMsg().Send()
	// 。。。。
}

//------------------- 以下测试及提供示例 机器人单独发送消息的方法 --------------

func TestSendText(t *testing.T) {
	bot.SwitchDebugMode(true)
	_ = bot.ToTextMsg("hello, I am miniBot.").Send()
}

func TestSendMarkdown(t *testing.T) {
	bot.SwitchDebugMode(true)
	at := &wxrobot.MsgAttachment{
		CallbackID: "123456",
		Actions: []wxrobot.MsgAction{{Name: "name1", Type: "button", Text: "展示文本1", Value: "111", ReplaceText: "已点击1"},
			{Name: "name2", Type: "button", Text: "展示文本2", Value: "222", ReplaceText: "已点击2"}},
	}
	_ = bot.ToMarkdownMsg(
		"**2019公司文化衫尺码收集**\n\n主题：2019文化衫尺码收集\n范围：所有<font color=\"warning\">正式员工+实习生</font>\n服装：统一为蓝色logo+白色T\n\n请选择你需要的尺码\n").
		Attachment(at).Send()

}

//go:embed test.png
var img []byte

func TestSendImage(t *testing.T) {
	bot.SwitchDebugMode(true)
	_ = bot.ToImageMsg().Image(img).Send()
}

func TestSendNews(t *testing.T) {
	bot.SwitchDebugMode(true)

	var news = &wxrobot.NewsArticle{
		Title:       "中秋节礼品领取",
		Description: "今年中秋节公司有豪礼相送",
		URL:         "https://www.qq.com/",
		PicURL:      "http://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png",
	}

	var news2 = &wxrobot.NewsArticle{
		Title:       "中秋节礼品领取",
		Description: "今年中秋节公司有豪礼相送",
		URL:         "https://www.qq.com/",
		PicURL:      "http://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png",
	}
	_ = bot.ToNewsMsg().Articles(news, news2).Send()
}

//go:embed test.docx
var file []byte

func TestSendFile(t *testing.T) {
	bot.SwitchDebugMode(true)
	_ = bot.ToFileMsg().File("test.docx", file).Send()
}
