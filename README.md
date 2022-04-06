## 企业微信机器人使用工具

**该工具旨在使用go语言发送机器人消息变得简单易用**

更多的使用示例可以参考[目录test下的测试用例](https://git.woa.com/mingkunhu/wxrobot/blob/master/test/wxrobot_test.go)
### 使用示例

> **1.简单的往群聊中发送文本消息**
> ```
> wxrobot.Bot("机器人的名字").WebhookURL("机器人的webhook地址").ToTextMsg("hello, I am miniBot.").Send()
> ```
> wxrobot.Bot("机器人的名字") 多次调用使用不同名字即可创建多个机器人，并分别配置自己的WebhookURL
> 
> 发送其他类型消息可以参考test目录下的测试用例
 

> **2.当需要接收并处理@机器人消息时**
> ```
> //1.只需要调用Serve方法，并传入机器人接收消息配置的Token， EncodingAESKey
> wxrobot.Bot("机器人的名字").WebhookURL("机器人的webhook地址").Serve("Token", "EncodingAESKey")
> 
> //2. 在你的服务器代码中注册一下机器人的处理器的路由
> r := mux.NewRouter()
> r = bot.Router(r)
> // r = bot.Router(r, "/{设置自定义path}")
> http.ListenAndServe(":80", r)
> ```
> 
> 将代码部署之后，在机器人接收消息配置页面的回调URL写上自己服务的URL， 如：http://www.demo.com/{自定义path} 
> 
> 现在，你的机器人应该可以正常接收用户的消息了！
 

> **3.接收并处理@机器人消息时，工具提供了几个handler注册函数以便于你进行处理不同的用户消息类型**
> ```
> //---- 注册各种事件消息回调处理函数 -------
>	// 注册发送给我的文本消息回调处理函数
>	bot.RegisterHandlerForText(processTextMsgCallback)
>	// 注册发送给我的图文消息回调处理函数
>	bot.RegisterHandlerForMixed(processMixedMsgCallback)
>	// 注册发送给我的图片消息回调处理函数 注意：图片消息目前仅支持私聊我的消息
>	bot.RegisterHandlerForImage(processImageMsgCallback)
>	// 注册发送给我的markdown消息点击事件回调处理函数
>	bot.RegisterHandlerForAttachment(processAttachmentMsgCallback)
>	// 注册关于我的事件消息回调处理函数
>	bot.RegisterHandlerForEvent(processEventMsgCallback)
> ```
> 
> 可以参考test目录下的测试用例