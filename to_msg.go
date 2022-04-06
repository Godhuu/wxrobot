package wxrobot

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
)

var md5Pool sync.Pool

func init() {
	md5Pool.New = func() interface{} {
		return md5.New()
	}
}

// toCommonMsg 回复消息
type toCommonMsg struct {
	*bot
	MsgType       string         `json:"msgtype"`
	ChatID        string         `json:"chatid,omitempty"`
	PostId        string         `json:"post_id,omitempty"`
	VisibleToUser string         `json:"visible_to_user,omitempty"`
	Text          *toMsgText     `json:"text,omitempty"`
	Markdown      *toMsgMarkdown `json:"markdown,omitempty"`
	News          *toMsgNews     `json:"news,omitempty"`
	Image         *toMsgImage    `json:"image,omitempty"`
	File          *toMsgFile     `json:"file,omitempty"`
}

// sendResponse 推送消息响应
type sendResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (tc *toCommonMsg) send() error {
	bt, err := json.Marshal(tc)
	if err != nil {
		return err
	}

	webhookURL := tc.webhookURL
	if tc.debug {
		webhookURL += "&debug=1"
	}
	wxRobot.logger.Debug("send json", string(bt))
	res, err := wxRobot.httpClient.Post(webhookURL, "application/json", bytes.NewReader(bt))
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var sendRes sendResponse
	_ = json.Unmarshal(body, &sendRes)

	wxRobot.logger.Debug("send return ", sendRes)

	if sendRes.ErrCode != 0 {
		return fmt.Errorf("%v", sendRes)
	}

	return nil
}

type toBaseMsg struct {
	bot            *bot
	msgType        string
	visibleToUsers []string
	chatids        []string
	postId         string
}

func (t *toBaseMsg) chatId(chatId ...string) {
	t.chatids = append(t.chatids, chatId...)
}

func (t *toBaseMsg) visible(user ...string) {
	t.visibleToUsers = append(t.visibleToUsers, user...)
}

func (t *toBaseMsg) buildCommonMsg() *toCommonMsg {
	cmsg := new(toCommonMsg)
	cmsg.bot = t.bot
	cmsg.MsgType = t.msgType
	cmsg.PostId = t.postId
	cmsg.ChatID = strings.Join(t.chatids, "|")
	cmsg.VisibleToUser = strings.Join(t.visibleToUsers, "|")

	return cmsg
}

// toMsgText 文本消息类型
type toMsgText struct {
	toBaseMsg
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list"`
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}

// MentionUser @用户(输入用户id)
func (t *toMsgText) MentionUser(user ...string) *toMsgText {
	t.MentionedList = append(t.MentionedList, user...)
	return t
}

// MentionMobile @用户(输入用户手机号)
func (t *toMsgText) MentionMobile(mobile ...string) *toMsgText {
	t.MentionedMobileList = append(t.MentionedMobileList, mobile...)
	return t
}

// PostId 小黑板帖子id，当前消息为小黑板回帖消息时带上，有且只有chatid指定了一个小黑板的时候生效
func (t *toMsgText) PostId(id string) *toMsgText {
	t.postId = id
	return t
}

// ChatId
/*
会话id，支持最多传100个。
可能是群聊会话，也可能是单聊会话或者小黑板会话，通过消息回调获得，也可以是userid。
特殊的，当chatid为“@all_group”时，表示对所有群广播，为“@all_subscriber”时表示对订阅范围内员工广播单聊消息，
为“@all_blackboard”时，表示对所有小黑板广播，为“@all”时，表示对所有群、所有订阅范围
*/
func (t *toMsgText) ChatId(chatId ...string) *toMsgText {
	t.chatId(chatId...)
	return t
}

// Visible
// 该消息只有指定的群成员或小黑板成员可见（其他成员不可见），有且只有chatid指定了一个群或一个小黑板的时候生效
func (t *toMsgText) Visible(user ...string) *toMsgText {
	t.visible(user...)
	return t
}

// Send 发送消息
func (t *toMsgText) Send() error {
	cmsg := t.buildCommonMsg()
	cmsg.Text = t
	return cmsg.send()
}

// toMsgMarkdown markdown消息类型
type toMsgMarkdown struct {
	toBaseMsg
	Content     string           `xml:"Content" json:"content"`
	ShortName   bool             `json:"at_short_name"`
	Attachments []*MsgAttachment `json:"attachments,omitempty"`
}

// PostId 小黑板帖子id，当前消息为小黑板回帖消息时带上，有且只有chatid指定了一个小黑板的时候生效
func (tm *toMsgMarkdown) PostId(id string) *toMsgMarkdown {
	tm.postId = id
	return tm
}

// AtShortName 设置为true
//markdown内容中@人指定为的短名字的方式，类型为bool值，设置为true则markdown中@xxx的表现为短名
func (tm *toMsgMarkdown) AtShortName() *toMsgMarkdown {
	tm.ShortName = true
	return tm
}

// Attachment attachments内容，目前仅支持button类型
func (tm *toMsgMarkdown) Attachment(attach ...*MsgAttachment) *toMsgMarkdown {
	tm.Attachments = append(tm.Attachments, attach...)
	return tm
}

// ChatId
/*
会话id，支持最多传100个。
可能是群聊会话，也可能是单聊会话或者小黑板会话，通过消息回调获得，也可以是userid。
特殊的，当chatid为“@all_group”时，表示对所有群广播，为“@all_subscriber”时表示对订阅范围内员工广播单聊消息，
为“@all_blackboard”时，表示对所有小黑板广播，为“@all”时，表示对所有群、所有订阅范围
*/
func (tm *toMsgMarkdown) ChatId(chatId ...string) *toMsgMarkdown {
	tm.chatId(chatId...)
	return tm
}

// Visible
// 该消息只有指定的群成员或小黑板成员可见（其他成员不可见），有且只有chatid指定了一个群或一个小黑板的时候生效
func (tm *toMsgMarkdown) Visible(user ...string) *toMsgMarkdown {
	tm.visible(user...)
	return tm
}

// Send 发送消息
func (tm *toMsgMarkdown) Send() error {
	cmsg := tm.buildCommonMsg()
	cmsg.Markdown = tm
	return cmsg.send()
}

// toMsgNews 图文消息
type toMsgNews struct {
	toBaseMsg
	Articles_ []*NewsArticle `json:"articles"`
}

// NewsArticle 图文项
type NewsArticle struct {
	Title       string `xml:"-" json:"title"`
	Description string `xml:"-" json:"description"`
	URL         string `xml:"-" json:"url"`
	PicURL      string `xml:"-" json:"picurl"`
}

// Articles 图文消息，一个图文消息支持1到8条图文 注意：NewsArticle.URL 回到URL不能为空，否则会报错
func (tm *toMsgNews) Articles(articles ...*NewsArticle) *toMsgNews {
	tm.Articles_ = append(tm.Articles_, articles...)
	return tm
}

// ChatId
/*
会话id，支持最多传100个。
可能是群聊会话，也可能是单聊会话或者小黑板会话，通过消息回调获得，也可以是userid。
特殊的，当chatid为“@all_group”时，表示对所有群广播，为“@all_subscriber”时表示对订阅范围内员工广播单聊消息，
为“@all_blackboard”时，表示对所有小黑板广播，为“@all”时，表示对所有群、所有订阅范围
*/
func (tm *toMsgNews) ChatId(chatId ...string) *toMsgNews {
	tm.chatId(chatId...)
	return tm
}

// Visible
// 该消息只有指定的群成员或小黑板成员可见（其他成员不可见），有且只有chatid指定了一个群或一个小黑板的时候生效
func (tm *toMsgNews) Visible(user ...string) *toMsgNews {
	tm.visible(user...)
	return tm
}

// Send 发送消息
func (tm *toMsgNews) Send() error {
	cmsg := tm.buildCommonMsg()
	cmsg.News = tm
	return cmsg.send()
}

// toMsgImage 图片消息
type toMsgImage struct {
	toBaseMsg
	Base64 string `json:"base64"`
	MD5    string `json:"md5"`
}

// Image 图片最大不能超过2M，支持JPG,PNG格式
func (tm *toMsgImage) Image(bts []byte) *toMsgImage {
	// 计算文件base64
	base64Str := base64.StdEncoding.EncodeToString(bts)
	// 计算文件md5
	md := md5Pool.Get().(hash.Hash)
	_, _ = md.Write(bts)
	md5Str := fmt.Sprintf("%x", md.Sum(nil))
	md.Reset()
	md5Pool.Put(md)
	tm.Base64 = base64Str
	tm.MD5 = md5Str
	return tm
}

// ChatId
/*
会话id，支持最多传100个。
可能是群聊会话，也可能是单聊会话或者小黑板会话，通过消息回调获得，也可以是userid。
特殊的，当chatid为“@all_group”时，表示对所有群广播，为“@all_subscriber”时表示对订阅范围内员工广播单聊消息，
为“@all_blackboard”时，表示对所有小黑板广播，为“@all”时，表示对所有群、所有订阅范围
*/
func (tm *toMsgImage) ChatId(chatId ...string) *toMsgImage {
	tm.chatId(chatId...)
	return tm
}

// Visible
// 该消息只有指定的群成员或小黑板成员可见（其他成员不可见），有且只有chatid指定了一个群或一个小黑板的时候生效
func (tm *toMsgImage) Visible(user ...string) *toMsgImage {
	tm.visible(user...)
	return tm
}

// Send 发送消息
func (tm *toMsgImage) Send() error {
	cmsg := tm.buildCommonMsg()
	cmsg.Image = tm
	return cmsg.send()
}

// toMsgFile 文件消息
type toMsgFile struct {
	toBaseMsg
	MediaID string `xml:"-" json:"media_id"`
}

// File 要求文件大小在5B~20M之间
func (tm *toMsgFile) File(fileName string, bts []byte) *toMsgFile {
	upLoadRes, err := tm.uploadFile(fileName, bts)
	if err != nil {
		return tm
	}

	tm.MediaID = upLoadRes.MediaId
	return tm
}

// ChatId
/*
会话id，支持最多传100个。
可能是群聊会话，也可能是单聊会话或者小黑板会话，通过消息回调获得，也可以是userid。
特殊的，当chatid为“@all_group”时，表示对所有群广播，为“@all_subscriber”时表示对订阅范围内员工广播单聊消息，
为“@all_blackboard”时，表示对所有小黑板广播，为“@all”时，表示对所有群、所有订阅范围
*/
func (tm *toMsgFile) ChatId(chatId ...string) *toMsgFile {
	tm.chatId(chatId...)
	return tm
}

// Visible
// 该消息只有指定的群成员或小黑板成员可见（其他成员不可见），有且只有chatid指定了一个群或一个小黑板的时候生效
func (tm *toMsgFile) Visible(user ...string) *toMsgFile {
	tm.visible(user...)
	return tm
}

// Send 发送消息
func (tm *toMsgFile) Send() error {
	if tm.MediaID == "" {
		return fmt.Errorf("请先调用File(...)方法上传文件， 或检查日志是否上传失败")
	}
	cmsg := tm.buildCommonMsg()
	cmsg.File = tm
	return cmsg.send()
}

// MsgAttachment markdown附加数据
type MsgAttachment struct {
	CallbackID string      `xml:"CallbackId" json:"callback_id"`
	Actions    []MsgAction `xml:"Actions" json:"actions"`
}

// MsgAction 操作
type MsgAction struct {
	Name        string `xml:"Name" json:"name"`
	Value       string `xml:"Value" json:"value"`
	Text        string `xml:"Text" json:"text"`
	Type        string `xml:"Type" json:"type"`
	BorderColor string `xml:"BorderColor" json:"border_color"`
	TextColor   string `xml:"TextColor" json:"text_color"`
	ReplaceText string `xml:"ReplaceText" json:"replace_text"`
}

// upLoadRes 上传文件的响应信息
type upLoadRes struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	MediaId   string `json:"media_id"` // 微信返回的文件id，机器人可以拿这个发送文件给指定用户/群
	Type      string `json:"type"`
	CreatedAt int64  `json:"created_at"`
}

func (tm *toMsgFile) uploadFile(filename string, pdfContent []byte) (*upLoadRes, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	part, err := bodyWriter.CreateFormFile("file1", filename)
	if err != nil {
		wxRobot.logger.Error(fmt.Sprintf("Cannot CreateFormFile for: %s , err: %v", filename, err))
		return nil, err
	}

	_, err = part.Write(pdfContent)
	if err != nil {
		wxRobot.logger.Error(fmt.Sprintf("Cannot Write file: %s , err: %v", filename, err))
		return nil, err
	}
	_ = bodyWriter.Close()
	req, err := http.NewRequest("POST", tm.bot.uploadUrl, bodyBuf)
	if err != nil {
		wxRobot.logger.Error("NewRequest err:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	resp, err := wxRobot.httpClient.Do(req)
	if err != nil {
		wxRobot.logger.Error("uploadFile send http request err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wxRobot.logger.Error("ioutil.ReadAll err ", err)
		return nil, err
	}

	uploadRes := &upLoadRes{}
	_ = json.Unmarshal(body, uploadRes)
	wxRobot.logger.Debug("uploadRes :", uploadRes)

	return uploadRes, nil
}
