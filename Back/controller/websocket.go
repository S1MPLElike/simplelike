package controller

import (
	"GoWorkingProject/Back/model"
	"GoWorkingProject/Back/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket连接管理
var (
	clients    = make(map[uint64]*websocket.Conn) // 用户ID到WebSocket连接的映射
	clientsMux sync.RWMutex                       // 保护clients的互斥锁
)

// WebSocket消息类型
const (
	MessageTypeChat    = "chat"    // 聊天消息
	MessageTypeStatus  = "status"  // 状态消息
	MessageTypeError   = "error"   // 错误消息
	MessageTypeHistory = "history" // 历史消息
)

// WebSocket消息结构
type WSMessage struct {
	Type    string      `json:"type"`    // 消息类型
	Payload interface{} `json:"payload"` // 消息内容
}

// 聊天消息结构
type ChatMessage struct {
	ID         uint64 `json:"id"`          // 消息ID
	SenderID   uint64 `json:"sender_id"`   // 发送者ID
	ReceiverID uint64 `json:"receiver_id"` // 接收者ID
	Content    string `json:"content"`     // 消息内容
	Type       string `json:"type"`        // 消息类型：text, image, emoji
	CreateTime string `json:"create_time"` // 创建时间
	IsSelf     bool   `json:"is_self"`     // 是否是自己发送的
}

// 状态消息结构
type StatusMessage struct {
	Type    string `json:"type"`    // 状态类型：online, offline
	UserID  uint64 `json:"user_id"` // 用户ID
	Message string `json:"message"` // 状态消息
}

// 错误消息结构
type ErrorMessage struct {
	Code    int    `json:"code"`    // 错误代码
	Message string `json:"message"` // 错误消息
}

// WebSocket升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有跨域请求
		return true
	},
}

// WebSocketHandler 处理WebSocket连接
func WebSocketHandler(c *gin.Context) {
	// 获取当前登录用户ID
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录"))
		return
	}

	log.Printf("WebSocket连接请求，用户ID: %d", userID)

	// 升级HTTP连接为WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	// 注册客户端连接
	clientsMux.Lock()
	clients[userID] = conn
	clientsMux.Unlock()

	log.Printf("WebSocket连接成功，用户ID: %d", userID)

	// 发送在线状态消息
	sendStatusMessage(userID, "online", "用户上线")

	// 处理连接关闭
	defer func() {
		log.Printf("WebSocket连接关闭，用户ID: %d", userID)

		// 移除客户端连接
		clientsMux.Lock()
		delete(clients, userID)
		clientsMux.Unlock()

		// 发送离线状态消息
		sendStatusMessage(userID, "offline", "用户下线")

		// 关闭连接
		conn.Close()
	}()

	// 处理消息
	for {
		// 读取消息
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket读取失败，用户ID: %d, 错误: %v", userID, err)
			break
		}

		// 处理文本消息
		if messageType == websocket.TextMessage {
			log.Printf("收到WebSocket消息，用户ID: %d, 消息: %s", userID, string(message))
			handleWSMessage(userID, message)
		}
	}
}

// 处理WebSocket消息
func handleWSMessage(userID uint64, message []byte) {
	// 解析消息
	var wsMsg WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		log.Printf("解析WebSocket消息失败: %v", err)
		sendErrorMessage(userID, utils.CodeParamError, "消息格式错误")
		return
	}

	// 根据消息类型处理
	switch wsMsg.Type {
	case MessageTypeChat:
		handleChatMessage(userID, wsMsg.Payload)
	default:
		log.Printf("未知消息类型: %s", wsMsg.Type)
		sendErrorMessage(userID, utils.CodeParamError, "未知消息类型")
	}
}

// 处理聊天消息
func handleChatMessage(userID uint64, payload interface{}) {
	log.Printf("处理聊天消息，发送者ID: %d", userID)

	// 解析聊天消息
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Printf("序列化消息失败: %v", err)
		sendErrorMessage(userID, utils.CodeParamError, "消息格式错误")
		return
	}

	var chatMsg struct {
		ReceiverID uint64 `json:"receiver_id"`
		Content    string `json:"content"`
		Type       string `json:"type"`
	}

	if err := json.Unmarshal(payloadJSON, &chatMsg); err != nil {
		log.Printf("解析聊天消息失败: %v", err)
		sendErrorMessage(userID, utils.CodeParamError, "消息格式错误")
		return
	}

	log.Printf("聊天消息解析成功，接收者ID: %d, 内容: %s", chatMsg.ReceiverID, chatMsg.Content)

	// 验证接收者是否存在
	var receiver model.User
	if err := model.DB.Where("id = ?", chatMsg.ReceiverID).First(&receiver).Error; err != nil {
		sendErrorMessage(userID, utils.CodeUserNotFound, "接收者不存在")
		return
	}

	// 验证是否是好友关系
	var count1, count2 int64
	model.DB.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", userID, chatMsg.ReceiverID).Count(&count1)
	model.DB.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", chatMsg.ReceiverID, userID).Count(&count2)

	if count1 == 0 || count2 == 0 {
		sendErrorMessage(userID, utils.CodeOperationFailed, "只能给好友发送消息")
		return
	}

	// 设置消息类型默认值
	messageType := chatMsg.Type
	if messageType == "" {
		messageType = "text"
	}

	// 创建消息
	message := model.Message{
		SenderID:   userID,
		ReceiverID: chatMsg.ReceiverID,
		Content:    chatMsg.Content,
		Type:       messageType,
		Status:     "unread",
	}

	if err := model.DB.Create(&message).Error; err != nil {
		sendErrorMessage(userID, utils.CodeDBError, "发送消息失败")
		return
	}

	// 构建发送给发送者的消息
	senderMsg := ChatMessage{
		ID:         message.ID,
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		Content:    message.Content,
		Type:       message.Type,
		CreateTime: message.CreateTime.Format("2006-01-02 15:04:05"),
		IsSelf:     true,
	}

	// 构建发送给接收者的消息
	receiverMsg := ChatMessage{
		ID:         message.ID,
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		Content:    message.Content,
		Type:       message.Type,
		CreateTime: message.CreateTime.Format("2006-01-02 15:04:05"),
		IsSelf:     false,
	}

	// 发送消息给发送者
	log.Printf("发送消息给发送者，用户ID: %d", userID)
	sendWSMessage(userID, MessageTypeChat, senderMsg)

	// 发送消息给接收者
	log.Printf("发送消息给接收者，用户ID: %d", chatMsg.ReceiverID)
	sendWSMessage(chatMsg.ReceiverID, MessageTypeChat, receiverMsg)
}

// 发送WebSocket消息
func sendWSMessage(userID uint64, messageType string, payload interface{}) {
	// 构建消息
	wsMsg := WSMessage{
		Type:    messageType,
		Payload: payload,
	}

	// 序列化消息
	message, err := json.Marshal(wsMsg)
	if err != nil {
		log.Printf("序列化WebSocket消息失败: %v", err)
		return
	}

	// 发送消息
	clientsMux.RLock()
	conn, ok := clients[userID]
	clientsMux.RUnlock()

	if ok {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("发送WebSocket消息失败: %v", err)
			// 移除无效连接
			clientsMux.Lock()
			delete(clients, userID)
			clientsMux.Unlock()
			conn.Close()
		} else {
			log.Printf("WebSocket消息发送成功，用户ID: %d", userID)
		}
	} else {
		log.Printf("用户 %d 不在线，无法发送WebSocket消息", userID)
	}
}

// 发送状态消息
func sendStatusMessage(userID uint64, statusType string, message string) {
	statusMsg := StatusMessage{
		Type:    statusType,
		UserID:  userID,
		Message: message,
	}

	sendWSMessage(userID, MessageTypeStatus, statusMsg)
}

// 发送错误消息
func sendErrorMessage(userID uint64, code int, message string) {
	errorMsg := ErrorMessage{
		Code:    code,
		Message: message,
	}

	sendWSMessage(userID, MessageTypeError, errorMsg)
}

// GetOnlineStatus 获取用户在线状态
func GetOnlineStatus(c *gin.Context) {
	// 获取目标用户ID
	targetIDStr := c.Param("id")
	var targetID uint64
	_, err := fmt.Sscanf(targetIDStr, "%d", &targetID)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "无效的用户ID"))
		return
	}

	// 检查用户是否在线
	clientsMux.RLock()
	_, online := clients[targetID]
	clientsMux.RUnlock()

	c.JSON(200, utils.Success(gin.H{
		"online": online,
	}))
}
