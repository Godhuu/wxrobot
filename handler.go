package wxrobot

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (r *bot) handler(res http.ResponseWriter, req *http.Request) {
	if r.msgCrypt == nil {
		panic("msgCrypt is nil, please first call Serve(...) ")
	}

	switch req.Method {
	case "POST":
		r.processData(res, req)
	case "GET":
		r.processEcho(res, req)
	default:
		//不做处理
	}
}

//企业微信的echo消息
func (r *bot) processEcho(res http.ResponseWriter, req *http.Request) {
	wxRobot.logger.Info("processEcho rawQuery:", req.URL.RawQuery)
	querys, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		wxRobot.logger.Info("ParseQuery err:", err)
		return
	}

	if len(querys["msg_signature"]) <= 0 || len(querys["timestamp"]) <= 0 || len(
		querys["nonce"]) <= 0 || len(querys["echostr"]) <= 0 {
		wxRobot.logger.Error("参数检查失败 。。。")
		return
	}

	verifyMsgSign := querys["msg_signature"][0]
	verifyTimestamp := querys["timestamp"][0]
	verifyNonce := querys["nonce"][0]
	verifyEchoStr := querys["echostr"][0]
	echoStr, cryptErr := r.msgCrypt.VerifyURL(verifyMsgSign, verifyTimestamp, verifyNonce, verifyEchoStr)
	if nil != cryptErr {
		wxRobot.logger.Error("verifyUrl fail", cryptErr)
		return
	}

	wxRobot.logger.Info("verifyUrl success echoStr", string(echoStr))
	_, _ = res.Write(echoStr)
}

// processPostData 处理正常的@我的消息
func (r *bot) processData(res http.ResponseWriter, req *http.Request) {
	wxRobot.logger.Info("processData rawQuery:", req.URL.RawQuery)
	querys, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		wxRobot.logger.Error("ParseQuery err:", err)
		return
	}

	if len(querys["msg_signature"]) <= 0 || len(querys["timestamp"]) <= 0 || len(querys["nonce"]) <= 0 {
		wxRobot.logger.Error("参数检查失败 。。。")
		return
	}

	reqMsgSign := querys["msg_signature"][0]
	reqTimestamp := querys["timestamp"][0]
	reqNonce := querys["nonce"][0]

	defer func() { _ = req.Body.Close() }()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		wxRobot.logger.Error("ioutil.ReadAll err ", err)
		return
	}

	wxRobot.logger.Info("processData body: ", string(body))
	msg, cryptErr := r.msgCrypt.DecryptMsg(reqMsgSign, reqTimestamp, reqNonce, body)
	if nil != cryptErr {
		wxRobot.logger.Error("DecryptMsg fail", cryptErr)
		return
	}

	wxRobot.logger.Debug("msg:  ", string(msg))
	var msgContent FromCommonMsg
	err = xml.Unmarshal(msg, &msgContent)
	if nil != err {
		wxRobot.logger.Error("Unmarshal fail：", err)
		return
	} else {
		wxRobot.logger.Debug("struct", msgContent)
	}

	go r.replyHandler(&msgContent, msg)
}

// 异步响应
func (r *bot) replyHandler(msgContent *FromCommonMsg, msgBody []byte) {
	switch msgContent.GetMsgType() {
	case MsgTypeText:
		var msg FromTextMsg
		_ = xml.Unmarshal(msgBody, &msg)
		msg.bot = r
		if r.textHandler == nil {
			defaultTextHandler(&msg)
			return
		}
		r.textHandler(&msg)
	case MsgTypeEvent:
		var msg FromEventMsg
		_ = xml.Unmarshal(msgBody, &msg)
		msg.bot = r
		if r.eventHandler == nil {
			defaultEventHandler(&msg)
			return
		}
		r.eventHandler(&msg)
	case MsgTypeImage:
		var msg FromImageMsg
		_ = xml.Unmarshal(msgBody, &msg)
		msg.bot = r
		if r.imageHandler == nil {
			defaultImageHandler(&msg)
			return
		}
		r.imageHandler(&msg)
	case MsgTypeMixed:
		var msg FromMixedMsg
		_ = xml.Unmarshal(msgBody, &msg)
		msg.bot = r
		if r.mixedHandler == nil {
			defaultMixedHandler(&msg)
			return
		}
		r.mixedHandler(&msg)
	case MsgTypeAttachment:
		var msg FromAttachmentMsg
		_ = xml.Unmarshal(msgBody, &msg)
		msg.bot = r
		if r.attachmentHandler == nil {
			defaultAttachmentHandler(&msg)
			return
		}
		r.attachmentHandler(&msg)
	default: //其他都不支持
		wxRobot.logger.Error("不支持的MsgType :", msgContent.MsgType)
		return
	}
}

func defaultEventHandler(msg *FromEventMsg) {
	switch msg.EventType() {
	case EnterChatEvent:
		wxRobot.logger.Debug(fmt.Sprintf("用户[%s][%s],进入了机器人[%s]单聊.", msg.From.UserId, msg.From.Name, msg.bot.name))
	case AddToChatEvent:
		wxRobot.logger.Debug(fmt.Sprintf("用户[%s][%s],将机器人[%s]拉入群聊.", msg.From.UserId, msg.From.Name, msg.bot.name))
	case DeleteFromChatEvent:
		wxRobot.logger.Debug(fmt.Sprintf("用户[%s][%s],将机器人[%s]从群聊删除.", msg.From.UserId, msg.From.Name, msg.bot.name))
	}
}
func defaultTextHandler(msg *FromTextMsg) {
	_ = msg.ToTextMsg(msg.Text.Content).Send()
}
func defaultImageHandler(msg *FromImageMsg) {
	_ = msg.ToTextMsg(msg.Image.ImageUrl).Send()
}
func defaultAttachmentHandler(msg *FromAttachmentMsg) {
	var content = fmt.Sprintf("[%s]选择了", msg.From.Name)
	for _, v := range msg.Attachment.Actions {
		content += fmt.Sprintf("选项[%s]-(value:%s)", v.Name, v.Value)
	}
	_ = msg.ToTextMsg(content).ChatId(msg.ChatId).Send()
}
func defaultMixedHandler(msg *FromMixedMsg) {
	var resMsg string
	for _, v := range msg.MixedMessage {
		if v.MsgType == "text" {
			resMsg = fmt.Sprintf("%s\n收到:%s", resMsg, v.Text.Content)
		}

		if v.MsgType == "image" {
			resMsg = fmt.Sprintf("%s\n收到图片:%s", resMsg, v.Image.ImageUrl)
		}
	}
	_ = msg.ToTextMsg(resMsg).Send()
}

type eventHandler func(msg *FromEventMsg)
type textHandler func(msg *FromTextMsg)
type imageHandler func(msg *FromImageMsg)
type attachmentHandler func(msg *FromAttachmentMsg)
type mixedHandler func(msg *FromMixedMsg)

// RegisterHandlerForEvent 注册关于我的事件消息回调处理函数
func (r *bot) RegisterHandlerForEvent(handler eventHandler) *bot {
	r.eventHandler = handler
	return r
}

// RegisterHandlerForText 注册发送给我的文本消息回调处理函数
func (r *bot) RegisterHandlerForText(handler textHandler) *bot {
	r.textHandler = handler
	return r
}

// RegisterHandlerForImage 注册发送给我的图片消息回调处理函数 注意：图片消息目前仅支持私聊我的消息
func (r *bot) RegisterHandlerForImage(handler imageHandler) *bot {
	r.imageHandler = handler
	return r
}

// RegisterHandlerForAttachment 注册发送给我的markdown消息点击事件回调处理函数
func (r *bot) RegisterHandlerForAttachment(handler attachmentHandler) *bot {
	r.attachmentHandler = handler
	return r
}

// RegisterHandlerForMixed 注册发送给我的图文消息回调处理函数
func (r *bot) RegisterHandlerForMixed(handler mixedHandler) *bot {
	r.mixedHandler = handler
	return r
}
