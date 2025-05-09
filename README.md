# 使用Go + WebSocket 实现带语音和图片发送的即时通讯系统

在开发实时聊天系统时，除了基础的文本消息，语音与图片的发送也是关键的功能之一。本文将介绍如何基于 Go、Gin、WebSocket 和 GORM 实现一个支持语音和图片消息的即时通讯系统。

## 一、项目结构概览

```bash
├── main.go
├── router/
│   └── router.go
├── ws/
│   └── handler.go         // WebSocket 处理器
├── service/
│   └── message_service.go // 消息处理服务
├── models/
│   ├── user.go
│   ├── message.go
│   ├── friendship.go
│   ├── group.go
│   └── group_member.go
├── uploads/
│   ├── audio/
│   └── image/
```

## 二、消息模型设计

```go
// models/message.go
type Message struct {
	gorm.Model
	Type      string `json:"type"`
	From      string `json:"from"`
	To        string `json:"to"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Extra     string `json:"extra"`
}
```

支持的消息类型包括：
- text（文本）
- image（图片）
- audio（语音）
- emoji（表情）
- voice_invite / video_invite / call_offer / ... （WebRTC 信令）

## 三、文件上传接口实现

```go
func UploadMedia(c *gin.Context) (string, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", appError.NewAppError(500, "文件上传失败", nil)
	}
	fileType := c.PostForm("type")
	var saveDir string
	if fileType == "audio" {
		saveDir = "./uploads/audio/"
	} else if fileType == "image" {
		saveDir = "./uploads/image/"
	} else {
		return "", appError.NewAppError(500, "未知的文件类型", nil)
	}
	_ = os.MkdirAll(saveDir, os.ModePerm)
	filename := time.Now().Format("20060102150405") + "_" + file.Filename
	savePath := filepath.Join(saveDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return "", appError.NewAppError(500, "保存文件失败", nil)
	}
	return "/uploads/" + fileType + "/" + filename, nil
}
```

### 文件上传接口路由

```go
// POST /api/v1/upload
func UploadMediaHandler() gin.HandlerFunc {
	return utils.HandleRequest(
		func(c *gin.Context, req *string) error {
			path, err := UploadMedia(c)
			if err != nil {
				return err
			}
			*req = path
			return nil
		},
		func(c *gin.Context, req *string) (any, *appError.AppError) {
			return gin.H{"url": *req}, nil
		},
	)
}
```

## 四、WebSocket 消息类型处理

```go
switch msg.Type {
case "audio", "image":
	if err := service.HandleMessage(db, msg); err != nil {
		c.SendError("发送" + msg.Type + "消息失败")
	}
case "message":
	// 普通文本消息
	...
}
```

## 五、客户端接收语音与图片

前端通过 WebSocket 接收到消息后，可根据 `type` 字段判断是 `audio` 还是 `image`：

```js
if (msg.type === 'audio') {
	let audio = new Audio(msg.content);
	document.body.appendChild(audio);
	audio.controls = true;
} else if (msg.type === 'image') {
	let img = new Image();
	img.src = msg.content;
	document.body.appendChild(img);
}
```

## 六、数据库建表 SQL

```sql
CREATE TABLE `message` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `created_at` DATETIME,
  `updated_at` DATETIME,
  `deleted_at` DATETIME,
  `type` VARCHAR(50),
  `from` VARCHAR(100),
  `to` VARCHAR(100),
  `content` TEXT,
  `timestamp` BIGINT,
  `extra` TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

其他表如 `user_basic`、`friendship`、`group`、`group_member` 等也一并建好。

---

## 七、总结

本文介绍了如何结合 Go 实现音频、图片上传及消息发送功能，构建一个类似微信语音发送的聊天系统。整体逻辑清晰，可自由拓展文字、语音、图片以外的富媒体消息。

 
