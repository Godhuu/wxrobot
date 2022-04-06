package wxrobot

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"sync"
	"time"
)

// wxRobot 全局机器人管理器
var wxRobot = &wxBot{logger: std, httpClient: defaultHttpClient}

var defaultHttpClient = &http.Client{
	Timeout: 30 * time.Second,
}

type wxBot struct {
	robots     sync.Map
	logger     Logger
	httpClient *http.Client
}

// SetLogger 设置日志记录器
func SetLogger(log Logger) {
	wxRobot.logger = log
}

// bot 企业微信机器人
type bot struct {
	name       string
	token      string // 接入验证的token
	aesKey     string // 接入验证的encodingAesKey
	receiverId string
	webhookURL string // 主动推送消息的地址
	uploadUrl  string // 上传文件地址
	debug      bool
	router     *mux.Router
	msgCrypt   *WXBizMsgCrypt

	eventHandler      eventHandler
	textHandler       textHandler
	imageHandler      imageHandler
	attachmentHandler attachmentHandler
	mixedHandler      mixedHandler
}

// Bot 新建或获取一个机器人
// name 机器人的名字/别名 name已存在时会返回已有的机器人，否则会新建一个
func Bot(name ...string) *bot {
	var _name string
	if len(name) > 0 {
		_name = name[0]
	}
	_bot, _ := wxRobot.robots.LoadOrStore(_name, &bot{name: _name})
	return _bot.(*bot)
}

// HttpClient 设置http.Client
func HttpClient(client *http.Client) *wxBot {
	wxRobot.httpClient = client
	return wxRobot
}

// WebhookURL 设置机器人的WebhookURL 用于推送消息
func (r *bot) WebhookURL(url string) *bot {
	r.webhookURL = url
	r.uploadUrl = strings.ReplaceAll(url, "send", "upload_media") + "&type=file"
	return r
}

// SwitchDebugMode 设置机器人的debug模式
func (r *bot) SwitchDebugMode(open bool) *bot {
	r.debug = open
	return r
}

// Serve 设置机器人的接收消息配置参数 用于与机器人的交互回调
func (r *bot) Serve(token, aesKey string, receiverId ...string) *bot {
	r.token = token
	r.aesKey = aesKey

	if len(receiverId) > 0 {
		r.receiverId = receiverId[0]
	}

	r.msgCrypt = NewWXBizMsgCrypt(token, aesKey, r.receiverId, XmlType)
	return r
}

// Router 返回接收消息的路由
func (r *bot) Router(router *mux.Router, path ...string) *mux.Router {
	if router == nil {
		r.router = mux.NewRouter()
	}
	_path := "/"
	if len(path) > 0 {
		_path = path[0]
	}

	r.router = router
	r.router.HandleFunc(_path, r.handler).Methods("POST", "GET")
	return r.router
}

// NewRouter 返回接收消息的路由
func (r *bot) NewRouter(path ...string) *mux.Router {
	r.router = mux.NewRouter()

	_path := "/"
	if len(path) > 0 {
		_path = path[0]
	}

	r.router.HandleFunc(_path, r.handler).Methods("POST", "GET")
	return r.router
}

// ToTextMsg 新建要发送的文本消息
func (r *bot) ToTextMsg(msg string) *toMsgText {
	text := new(toMsgText)
	text.bot = r
	text.msgType = "text"
	text.Content = msg
	return text
}

// ToMarkdownMsg 新建要发送的markdown消息
func (r *bot) ToMarkdownMsg(markdown string) *toMsgMarkdown {
	t := new(toMsgMarkdown)
	t.bot = r
	t.msgType = "markdown"
	t.Content = markdown
	return t
}

// ToImageMsg 新建要发送的图片消息
func (r *bot) ToImageMsg() *toMsgImage {
	t := new(toMsgImage)
	t.bot = r
	t.msgType = "image"
	return t
}

// ToNewsMsg 新建要发送的图文消息
func (r *bot) ToNewsMsg() *toMsgNews {
	t := new(toMsgNews)
	t.bot = r
	t.msgType = "news"
	return t
}

// ToFileMsg 新建要发送的文件消息
func (r *bot) ToFileMsg() *toMsgFile {
	t := new(toMsgFile)
	t.bot = r
	t.msgType = "file"
	return t
}
