package controller

import (
	"GoWorkingProject/Back/model"
	"GoWorkingProject/Back/utils"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// SendMessage 发送消息接口
func SendMessage(c *gin.Context) {
	// 获取当前登录用户ID
	senderID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录"))
		return
	}

	// 解析请求参数
	var req struct {
		ReceiverID uint64 `json:"receiver_id" binding:"required"`
		Content    string `json:"content" binding:"required"`
		Type       string `json:"type" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "参数错误: "+err.Error()))
		return
	}

	// 验证接收者是否存在
	var receiver model.User
	if err := model.DB.Where("id = ?", req.ReceiverID).First(&receiver).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeUserNotFound, "接收者不存在"))
		return
	}

	// 验证是否是好友关系
	var count1, count2 int64
	model.DB.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", senderID, req.ReceiverID).Count(&count1)
	model.DB.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", req.ReceiverID, senderID).Count(&count2)

	if count1 == 0 || count2 == 0 {
		c.JSON(200, utils.Fail(utils.CodeOperationFailed, "只能给好友发送消息"))
		return
	}

	// 设置消息类型默认值
	messageType := req.Type
	if messageType == "" {
		messageType = "text"
	}

	// 创建消息
	message := model.Message{
		SenderID:   senderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		Type:       messageType,
		Status:     "unread",
	}

	if err := model.DB.Create(&message).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "发送消息失败: "+err.Error()))
		return
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

	// 发送消息给接收者
	sendWSMessage(req.ReceiverID, MessageTypeChat, receiverMsg)

	c.JSON(200, utils.Success(gin.H{
		"message_id":  message.ID,
		"content":     message.Content,
		"type":        message.Type,
		"create_time": message.CreateTime,
	}))
}

// GetMessageHistory 获取消息历史接口
func GetMessageHistory(c *gin.Context) {
	// 获取当前登录用户ID
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录"))
		return
	}

	// 获取聊天对象ID
	friendIDStr := c.Param("id")
	var friendID uint64
	_, err = fmt.Sscanf(friendIDStr, "%d", &friendID)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "无效的好友ID"))
		return
	}

	// 验证是否是好友关系
	var count1, count2 int64
	model.DB.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", userID, friendID).Count(&count1)
	model.DB.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", friendID, userID).Count(&count2)

	if count1 == 0 || count2 == 0 {
		c.JSON(200, utils.Fail(utils.CodeOperationFailed, "只能查看好友的消息历史"))
		return
	}

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	var page int
	_, err = fmt.Sscanf(pageStr, "%d", &page)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize := 20
	offset := (page - 1) * pageSize

	// 查询消息历史
	var messages []model.Message
	if err := model.DB.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, friendID, friendID, userID,
	).Order("create_time desc").Limit(pageSize).Offset(offset).Find(&messages).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "获取消息历史失败: "+err.Error()))
		return
	}

	// 标记对方发送的消息为已读
	model.DB.Model(&model.Message{}).Where(
		"sender_id = ? AND receiver_id = ? AND status = ?",
		friendID, userID, "unread",
	).Update("status", "read")

	c.JSON(200, utils.Success(messages))
}

// GetUnreadMessageCount 获取未读消息数量接口
func GetUnreadMessageCount(c *gin.Context) {
	// 获取当前登录用户ID
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录"))
		return
	}

	// 查询未读消息数量
	var count int64
	model.DB.Model(&model.Message{}).Where("receiver_id = ? AND status = ?", userID, "unread").Count(&count)

	c.JSON(200, utils.Success(gin.H{
		"unread_count": count,
	}))
}

// GetUnreadMessageCountByFriend 获取按好友分组的未读消息数量接口
func GetUnreadMessageCountByFriend(c *gin.Context) {
	// 获取当前登录用户ID
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录"))
		return
	}

	// 查询按好友分组的未读消息数量
	type UnreadCount struct {
		SenderID uint64 `json:"sender_id"`
		Count    int64  `json:"count"`
	}

	var unreadCounts []UnreadCount
	model.DB.Model(&model.Message{}).Select("sender_id, count(*) as count").Where("receiver_id = ? AND status = ?", userID, "unread").Group("sender_id").Scan(&unreadCounts)

	// 转换为map格式
	unreadMap := make(map[uint64]int64)
	for _, uc := range unreadCounts {
		unreadMap[uc.SenderID] = uc.Count
	}

	c.JSON(200, utils.Success(gin.H{
		"unread_counts": unreadMap,
	}))
}

// UploadImage 上传图片接口
func UploadImage(c *gin.Context) {
	// 获取当前登录用户ID
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录"))
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "获取文件失败: "+err.Error()))
		return
	}

	// 验证文件类型
	ext := getFileExtension(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		c.JSON(200, utils.Fail(utils.CodeParamError, "只支持jpg、jpeg、png、gif格式的图片"))
		return
	}

	// 创建保存目录
	uploadDir := "./uploads/images"
	if err := ensureDir(uploadDir); err != nil {
		c.JSON(200, utils.Fail(utils.CodeServerError, "创建目录失败: "+err.Error()))
		return
	}

	// 生成文件名
	fileName := generateFileName(userID, ext)
	filePath := uploadDir + "/" + fileName

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(200, utils.Fail(utils.CodeServerError, "保存文件失败: "+err.Error()))
		return
	}

	// 生成图片URL
	imageURL := "/uploads/images/" + fileName

	c.JSON(200, utils.Success(gin.H{
		"image_url": imageURL,
	}))
}

// 辅助函数：获取文件扩展名
func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}

// 辅助函数：确保目录存在
func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// 辅助函数：生成文件名
func generateFileName(userID uint64, ext string) string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%d_%d%s", userID, timestamp, ext)
}
